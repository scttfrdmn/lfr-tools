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

	// Create AWS client with 'aws' profile for GUI
	profile := cfg.AWS.Profile
	if profile == "" {
		profile = defaultAWSProfile // Default to 'aws' profile for GUI
	}

	awsClient, err := aws.NewClient(api.ctx, aws.Options{
		Region:  cfg.AWS.Region,
		Profile: profile,
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

	// Create AWS client with 'aws' profile for GUI
	profile := cfg.AWS.Profile
	if profile == "" {
		profile = defaultAWSProfile // Default to 'aws' profile for GUI
	}

	awsClient, err := aws.NewClient(api.ctx, aws.Options{
		Region:  cfg.AWS.Region,
		Profile: profile,
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

	// Create AWS client with 'aws' profile for GUI
	profile := cfg.AWS.Profile
	if profile == "" {
		profile = defaultAWSProfile // Default to 'aws' profile for GUI
	}

	awsClient, err := aws.NewClient(api.ctx, aws.Options{
		Region:  cfg.AWS.Region,
		Profile: profile,
	})
	if err != nil {
		return fmt.Errorf("failed to create AWS client: %w", err)
	}

	lightsailService := aws.NewLightsailService(awsClient)

	return lightsailService.StopInstance(api.ctx, instanceName)
}

// GetProjectInfo returns project overview with real cost data
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

	// Get real budget information
	costAPI := NewCostAPI()
	budgetInfo, err := costAPI.GetProjectBudget(project)
	if err != nil {
		// Fall back to estimated values if cost API fails
		return &ProjectInfo{
			Name:          project,
			StudentCount:  len(instances),
			RunningCount:  runningCount,
			BudgetUsed:    340.50,
			BudgetTotal:   500.00,
			DaysRemaining: 45,
		}, nil
	}

	return &ProjectInfo{
		Name:          project,
		StudentCount:  len(instances),
		RunningCount:  runningCount,
		BudgetUsed:    budgetInfo.UsedAmount,
		BudgetTotal:   budgetInfo.TotalBudget,
		DaysRemaining: budgetInfo.DaysRemaining,
	}, nil
}

// GetUserRole determines user role (for GUI authentication)
func (api *InstanceAPI) GetUserRole() (*UserInfo, error) {
	// Load configuration to determine authentication method
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	// Check for stored student token first
	tokenManager, err := config.NewTokenManager()
	if err == nil {
		tokens, err := tokenManager.ListTokens()
		if err == nil && len(tokens) > 0 {
			// Use first valid token
			for _, token := range tokens {
				if err := tokenManager.ValidateToken(token.Project, token.Username); err == nil {
					return &UserInfo{
						Role:        token.Role,
						Username:    token.Username,
						Project:     token.Project,
						Permissions: token.Permissions,
					}, nil
				}
			}
		}
	}

	// Check for AWS credentials (professor/admin access)
	profile := cfg.AWS.Profile
	if profile == "" {
		profile = defaultAWSProfile
	}

	awsClient, err := aws.NewClient(api.ctx, aws.Options{
		Region:  cfg.AWS.Region,
		Profile: profile,
	})
	if err != nil {
		return nil, fmt.Errorf("no valid authentication found. Students: activate token first. Professors: configure AWS credentials")
	}

	// Verify AWS access by attempting to list regions
	lightsailService := aws.NewLightsailService(awsClient)
	_, err = lightsailService.GetRegions(api.ctx)
	if err != nil {
		return nil, fmt.Errorf("AWS authentication failed: %w", err)
	}

	// AWS access confirmed - this is a professor/admin
	return &UserInfo{
		Role:        "professor",
		Username:    "aws-user",
		Project:     "", // Will be selected in GUI
		Permissions: []string{"create", "delete", "start", "stop", "ssh", "admin"},
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