package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/scttfrdmn/lfr-tools/internal/aws"
	"github.com/scttfrdmn/lfr-tools/internal/config"
	"github.com/scttfrdmn/lfr-tools/internal/types"
)

// LFRService provides the main service interface for the GUI
type LFRService struct {
	ctx context.Context
}

// NewLFRService creates a new LFR service instance
func NewLFRService() *LFRService {
	return &LFRService{
		ctx: context.Background(),
	}
}

// UserInfo represents user information for the frontend
type UserInfo struct {
	Role        string `json:"role"`        // student, ta, professor, admin
	Username    string `json:"username"`
	Project     string `json:"project"`
	Permissions []string `json:"permissions"`
}

// InstanceInfo represents instance information for the frontend
type InstanceInfo struct {
	Name      string            `json:"name"`
	State     string            `json:"state"`
	PublicIP  string            `json:"public_ip"`
	Blueprint string            `json:"blueprint"`
	Bundle    string            `json:"bundle"`
	Region    string            `json:"region"`
	Tags      map[string]string `json:"tags"`
	Username  string            `json:"username"` // Extracted from instance name
}

// ProjectInfo represents project/class information
type ProjectInfo struct {
	Name          string `json:"name"`
	StudentCount  int    `json:"student_count"`
	RunningCount  int    `json:"running_count"`
	BudgetUsed    float64 `json:"budget_used"`
	BudgetTotal   float64 `json:"budget_total"`
	DaysRemaining int    `json:"days_remaining"`
}

// GetUserRole determines the current user's role and permissions
func (s *LFRService) GetUserRole() (*UserInfo, error) {
	// Load configuration to determine role
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	// For now, detect role based on AWS profile access
	// In production, this would check stored tokens or authentication
	userInfo := &UserInfo{
		Role:        "professor", // Default for GUI development
		Username:    "demo-user",
		Project:     "demo-class",
		Permissions: []string{"create", "delete", "start", "stop", "ssh"},
	}

	return userInfo, nil
}

// ListInstances returns all instances for the current user/project
func (s *LFRService) ListInstances(project string) ([]*InstanceInfo, error) {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	// Create AWS client
	awsClient, err := aws.NewClient(s.ctx, aws.Options{
		Region:  cfg.AWS.Region,
		Profile: cfg.AWS.Profile,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS client: %w", err)
	}

	lightsailService := aws.NewLightsailService(awsClient)

	// Get instances
	instances, err := lightsailService.ListInstances(s.ctx, project)
	if err != nil {
		return nil, fmt.Errorf("failed to list instances: %w", err)
	}

	// Convert to frontend format
	var instancesInfo []*InstanceInfo
	for _, instance := range instances {
		// Extract username from instance name
		username := extractUsername(instance.Name)

		instanceInfo := &InstanceInfo{
			Name:      instance.Name,
			State:     instance.State,
			PublicIP:  instance.PublicIP,
			Blueprint: instance.Blueprint,
			Bundle:    instance.Bundle,
			Region:    instance.Region,
			Tags:      instance.Tags,
			Username:  username,
		}

		instancesInfo = append(instancesInfo, instanceInfo)
	}

	return instancesInfo, nil
}

// StartInstance starts a specific instance
func (s *LFRService) StartInstance(instanceName string) error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Create AWS client
	awsClient, err := aws.NewClient(s.ctx, aws.Options{
		Region:  cfg.AWS.Region,
		Profile: cfg.AWS.Profile,
	})
	if err != nil {
		return fmt.Errorf("failed to create AWS client: %w", err)
	}

	lightsailService := aws.NewLightsailService(awsClient)

	// Start the instance
	err = lightsailService.StartInstance(s.ctx, instanceName)
	if err != nil {
		return fmt.Errorf("failed to start instance: %w", err)
	}

	return nil
}

// StopInstance stops a specific instance
func (s *LFRService) StopInstance(instanceName string) error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Create AWS client
	awsClient, err := aws.NewClient(s.ctx, aws.Options{
		Region:  cfg.AWS.Region,
		Profile: cfg.AWS.Profile,
	})
	if err != nil {
		return fmt.Errorf("failed to create AWS client: %w", err)
	}

	lightsailService := aws.NewLightsailService(awsClient)

	// Stop the instance
	err = lightsailService.StopInstance(s.ctx, instanceName)
	if err != nil {
		return fmt.Errorf("failed to stop instance: %w", err)
	}

	return nil
}

// GetProjectInfo returns project/class overview information
func (s *LFRService) GetProjectInfo(project string) (*ProjectInfo, error) {
	// Get instances for project
	instances, err := s.ListInstances(project)
	if err != nil {
		return nil, fmt.Errorf("failed to get instances: %w", err)
	}

	// Calculate statistics
	runningCount := 0
	for _, instance := range instances {
		if instance.State == "running" {
			runningCount++
		}
	}

	projectInfo := &ProjectInfo{
		Name:          project,
		StudentCount:  len(instances),
		RunningCount:  runningCount,
		BudgetUsed:    340.50, // TODO: Calculate from actual costs
		BudgetTotal:   500.00, // TODO: Get from project configuration
		DaysRemaining: 45,     // TODO: Calculate from project end date
	}

	return projectInfo, nil
}

// ConnectToInstance initiates SSH connection for a student
func (s *LFRService) ConnectToInstance(username, project string) (string, error) {
	// This would integrate with our existing SSH connection logic
	// For now, return connection information

	instances, err := s.ListInstances(project)
	if err != nil {
		return "", fmt.Errorf("failed to list instances: %w", err)
	}

	// Find user's instance
	for _, instance := range instances {
		if instance.Username == username {
			if instance.State != "running" {
				return "", fmt.Errorf("instance is %s, not running", instance.State)
			}

			if instance.PublicIP == "" {
				return "", fmt.Errorf("instance has no public IP")
			}

			// Return SSH connection command
			sshCommand := fmt.Sprintf("ssh -i ~/.ssh/lfr-tools/LightsailDefaultKey.pem -o StrictHostKeyChecking=no ubuntu@%s", instance.PublicIP)
			return sshCommand, nil
		}
	}

	return "", fmt.Errorf("no instance found for user: %s", username)
}

// Helper function to extract username from instance name
func extractUsername(instanceName string) string {
	// Instance names follow pattern: username-blueprint
	parts := strings.Split(instanceName, "-")
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}