// Package api provides public APIs for GUI integration
package api

import (
	"context"
	"fmt"

	"strings"

	"github.com/scttfrdmn/lfr-tools/internal/aws"
	"github.com/scttfrdmn/lfr-tools/internal/config"
)

// InstanceAPI provides instance management operations for GUI
type InstanceAPI struct {
	ctx context.Context
}

// NewInstanceAPI creates a new instance API
func NewInstanceAPI() *InstanceAPI {
	return &InstanceAPI{
		ctx: context.Background(),
	}
}

// InstanceInfo represents instance information for GUI consumption
type InstanceInfo struct {
	Name      string            `json:"name"`
	State     string            `json:"state"`
	PublicIP  string            `json:"public_ip"`
	Blueprint string            `json:"blueprint"`
	Bundle    string            `json:"bundle"`
	Region    string            `json:"region"`
	Tags      map[string]string `json:"tags"`
	Username  string            `json:"username"`
}

// ProjectInfo represents project overview information
type ProjectInfo struct {
	Name          string  `json:"name"`
	StudentCount  int     `json:"student_count"`
	RunningCount  int     `json:"running_count"`
	BudgetUsed    float64 `json:"budget_used"`
	BudgetTotal   float64 `json:"budget_total"`
	DaysRemaining int     `json:"days_remaining"`
}

// UserInfo represents user role and permissions
type UserInfo struct {
	Role        string   `json:"role"`
	Username    string   `json:"username"`
	Project     string   `json:"project"`
	Permissions []string `json:"permissions"`
}

// ListInstances returns instances for a project
func (api *InstanceAPI) ListInstances(project string) ([]*InstanceInfo, error) {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	// Create AWS client
	awsClient, err := aws.NewClient(api.ctx, aws.Options{
		Region:  cfg.AWS.Region,
		Profile: cfg.AWS.Profile,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS client: %w", err)
	}

	lightsailService := aws.NewLightsailService(awsClient)

	// Get instances
	instances, err := lightsailService.ListInstances(api.ctx, project)
	if err != nil {
		return nil, fmt.Errorf("failed to list instances: %w", err)
	}

	// Convert to GUI format
	var guiInstances []*InstanceInfo
	for _, instance := range instances {
		guiInstance := &InstanceInfo{
			Name:      instance.Name,
			State:     instance.State,
			PublicIP:  instance.PublicIP,
			Blueprint: instance.Blueprint,
			Bundle:    instance.Bundle,
			Region:    instance.Region,
			Tags:      instance.Tags,
			Username:  extractUsernameFromInstanceName(instance.Name),
		}
		guiInstances = append(guiInstances, guiInstance)
	}

	return guiInstances, nil
}

// StartInstance starts an instance
func (api *InstanceAPI) StartInstance(instanceName string) error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Create AWS client
	awsClient, err := aws.NewClient(api.ctx, aws.Options{
		Region:  cfg.AWS.Region,
		Profile: cfg.AWS.Profile,
	})
	if err != nil {
		return fmt.Errorf("failed to create AWS client: %w", err)
	}

	lightsailService := aws.NewLightsailService(awsClient)

	return lightsailService.StartInstance(api.ctx, instanceName)
}

// StopInstance stops an instance
func (api *InstanceAPI) StopInstance(instanceName string) error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Create AWS client
	awsClient, err := aws.NewClient(api.ctx, aws.Options{
		Region:  cfg.AWS.Region,
		Profile: cfg.AWS.Profile,
	})
	if err != nil {
		return fmt.Errorf("failed to create AWS client: %w", err)
	}

	lightsailService := aws.NewLightsailService(awsClient)

	return lightsailService.StopInstance(api.ctx, instanceName)
}

// GetProjectInfo returns project overview
func (api *InstanceAPI) GetProjectInfo(project string) (*ProjectInfo, error) {
	instances, err := api.ListInstances(project)
	if err != nil {
		return nil, fmt.Errorf("failed to get instances: %w", err)
	}

	runningCount := 0
	for _, instance := range instances {
		if instance.State == "running" {
			runningCount++
		}
	}

	return &ProjectInfo{
		Name:          project,
		StudentCount:  len(instances),
		RunningCount:  runningCount,
		BudgetUsed:    340.50, // TODO: Implement real cost calculation
		BudgetTotal:   500.00, // TODO: Get from project configuration
		DaysRemaining: 45,     // TODO: Calculate from time boundaries
	}, nil
}

// GetUserRole determines user role (for GUI authentication)
func (api *InstanceAPI) GetUserRole() (*UserInfo, error) {
	// For now, return demo user info
	// TODO: Implement real authentication
	return &UserInfo{
		Role:        "professor",
		Username:    "demo-user",
		Project:     "demo-class",
		Permissions: []string{"create", "delete", "start", "stop", "ssh"},
	}, nil
}

// Helper function to extract username from instance name
func extractUsernameFromInstanceName(instanceName string) string {
	// Instance names follow pattern: username-blueprint
	parts := strings.Split(instanceName, "-")
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}