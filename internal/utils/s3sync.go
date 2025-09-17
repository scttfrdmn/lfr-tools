// Package utils provides S3 synchronization utilities for educational access.
package utils

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"

	"github.com/scttfrdmn/lfr-tools/internal/aws"
	"github.com/scttfrdmn/lfr-tools/internal/types"
)

// S3SyncConfig holds configuration for S3 status synchronization.
type S3SyncConfig struct {
	Enabled   bool   `json:"enabled"`
	Bucket    string `json:"bucket"`
	Project   string `json:"project"`
	AutoSync  bool   `json:"auto_sync"`
}

// UpdateInstanceStatusInS3 automatically updates instance status in S3 if configured.
func UpdateInstanceStatusInS3(ctx context.Context, instance *types.Instance) error {
	// Check if S3 sync is enabled
	syncConfig := getS3SyncConfig()
	if !syncConfig.Enabled || syncConfig.Bucket == "" {
		return nil // S3 sync not configured, skip silently
	}

	// Extract project from instance tags or use configured project
	project := instance.Tags["Project"]
	if project == "" {
		project = syncConfig.Project
	}
	if project == "" {
		return nil // No project context, skip
	}

	// Extract username from instance name
	username := ExtractUsernameFromInstance(instance.Name)
	if username == "" {
		return nil // Can't determine username, skip
	}

	// Create AWS client
	awsClient, err := aws.NewClient(ctx, aws.Options{
		Region:  viper.GetString("aws.region"),
		Profile: viper.GetString("aws.profile"),
	})
	if err != nil {
		// Don't fail the main operation if S3 sync fails
		fmt.Fprintf(os.Stderr, "Warning: Failed to create AWS client for S3 sync: %v\n", err)
		return nil
	}

	s3Service := aws.NewS3Service(awsClient)

	// Create status object
	status := &aws.StudentStatus{
		State:       instance.State,
		PublicIP:    instance.PublicIP,
		LastUpdated: time.Now(),
	}

	// Update status in S3
	err = s3Service.UpdateStudentStatus(ctx, syncConfig.Bucket, project, username, status)
	if err != nil {
		// Don't fail the main operation if S3 sync fails
		fmt.Fprintf(os.Stderr, "Warning: Failed to update S3 status for %s: %v\n", username, err)
		return nil
	}

	return nil
}

// UpdateMultipleInstancesInS3 updates status for multiple instances.
func UpdateMultipleInstancesInS3(ctx context.Context, instances []*types.Instance) {
	for _, instance := range instances {
		// Update each instance status (errors are logged, not returned)
		_ = UpdateInstanceStatusInS3(ctx, instance)
	}
}

// getS3SyncConfig gets S3 sync configuration from environment or config.
func getS3SyncConfig() S3SyncConfig {
	// Check for S3 sync configuration
	bucket := viper.GetString("students.s3_bucket")
	project := viper.GetString("students.project")
	enabled := viper.GetBool("students.s3_sync_enabled")

	// Also check environment variables
	if envBucket := os.Getenv("LFR_STUDENTS_S3_BUCKET"); envBucket != "" {
		bucket = envBucket
		enabled = true
	}
	if envProject := os.Getenv("LFR_STUDENTS_PROJECT"); envProject != "" {
		project = envProject
	}

	return S3SyncConfig{
		Enabled:  enabled,
		Bucket:   bucket,
		Project:  project,
		AutoSync: enabled,
	}
}

// ExtractUsernameFromInstance extracts username from instance name.
func ExtractUsernameFromInstance(instanceName string) string {
	// Instance names follow pattern: username-blueprint
	parts := strings.Split(instanceName, "-")
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}

// GetS3SyncConfig exposes S3 sync configuration for external use.
func GetS3SyncConfig() S3SyncConfig {
	return getS3SyncConfig()
}

// EnableS3Sync enables S3 synchronization for a project.
func EnableS3Sync(project, bucket string) error {
	viper.Set("students.s3_sync_enabled", true)
	viper.Set("students.s3_bucket", bucket)
	viper.Set("students.project", project)

	// Save configuration
	return viper.WriteConfig()
}

// DisableS3Sync disables S3 synchronization.
func DisableS3Sync() error {
	viper.Set("students.s3_sync_enabled", false)
	return viper.WriteConfig()
}