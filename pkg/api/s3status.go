// Package api provides S3 status communication for educational access
package api

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	awsInternal "github.com/scttfrdmn/lfr-tools/internal/aws"
	"github.com/scttfrdmn/lfr-tools/internal/config"
	"github.com/scttfrdmn/lfr-tools/internal/lfrutils"
)

// S3StatusAPI provides S3 status communication functionality
type S3StatusAPI struct {
	ctx context.Context
}

// NewS3StatusAPI creates a new S3 status API
func NewS3StatusAPI() *S3StatusAPI {
	return &S3StatusAPI{
		ctx: context.Background(),
	}
}

// StudentStatusInfo represents comprehensive student status
type StudentStatusInfo struct {
	Username        string    `json:"username"`
	Project         string    `json:"project"`
	InstanceName    string    `json:"instance_name"`
	State           string    `json:"state"`
	PublicIP        string    `json:"public_ip"`
	LastUpdated     time.Time `json:"last_updated"`
	BudgetUsed      float64   `json:"budget_used"`
	BudgetTotal     float64   `json:"budget_total"`
	AccessExpires   time.Time `json:"access_expires"`
	StartRequested  bool      `json:"start_requested"`
	RequestedAt     time.Time `json:"requested_at,omitempty"`
}

// UpdateStudentStatus updates a student's status in S3 automatically
func (api *S3StatusAPI) UpdateStudentStatus(username, project string) error {
	// Check if S3 sync is enabled
	syncConfig := lfrutils.GetS3SyncConfig()
	if !syncConfig.Enabled {
		// S3 sync not configured, skip silently
		return nil
	}

	// Get current instance information
	instanceAPI := NewInstanceAPI()
	instances, err := instanceAPI.ListInstances(project)
	if err != nil {
		return fmt.Errorf("failed to get instance information: %w", err)
	}

	// Find user's instance
	var userInstance *InstanceInfo
	for _, instance := range instances {
		if instance.Username == username {
			userInstance = instance
			break
		}
	}

	if userInstance == nil {
		return fmt.Errorf("no instance found for user %s in project %s", username, project)
	}

	// Get budget information
	costAPI := NewCostAPI()
	budgetInfo, err := costAPI.GetProjectBudget(project)
	if err != nil {
		// Continue with default budget if cost API fails
		budgetInfo = &BudgetInfo{
			UsedAmount:  0,
			TotalBudget: 25.0, // Default student budget
		}
	}

	// Create comprehensive status
	status := &StudentStatusInfo{
		Username:     username,
		Project:      project,
		InstanceName: userInstance.Name,
		State:        userInstance.State,
		PublicIP:     userInstance.PublicIP,
		LastUpdated:  time.Now(),
		BudgetUsed:   budgetInfo.UsedAmount / float64(len(instances)), // Rough per-student estimate
		BudgetTotal:  25.0, // Default student budget
		AccessExpires: time.Now().AddDate(0, 6, 0), // 6 months default
	}

	// Update S3 with status
	return api.uploadStatusToS3(syncConfig.Bucket, status)
}

// uploadStatusToS3 uploads student status to S3
func (api *S3StatusAPI) uploadStatusToS3(bucket string, status *StudentStatusInfo) error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Create AWS client
	profile := cfg.AWS.Profile
	if profile == "" {
		profile = defaultAWSProfile
	}

	awsClient, err := awsInternal.NewClient(api.ctx, awsInternal.Options{
		Region:  cfg.AWS.Region,
		Profile: profile,
	})
	if err != nil {
		return fmt.Errorf("failed to create AWS client: %w", err)
	}

	s3Service := awsInternal.NewS3Service(awsClient)

	// Convert to S3 status format
	s3Status := &awsInternal.StudentStatus{
		State:           status.State,
		PublicIP:        status.PublicIP,
		LastUpdated:     status.LastUpdated,
		BudgetRemaining: status.BudgetTotal - status.BudgetUsed,
		AccessExpires:   status.AccessExpires,
	}

	// Update in S3
	return s3Service.UpdateStudentStatus(api.ctx, bucket, status.Project, status.Username, s3Status)
}

// CheckStartRequests checks for pending student start requests in S3
func (api *S3StatusAPI) CheckStartRequests(project string) ([]*StudentStartRequest, error) {
	// Check if S3 sync is enabled
	syncConfig := lfrutils.GetS3SyncConfig()
	if !syncConfig.Enabled {
		return nil, fmt.Errorf("S3 sync not enabled for project %s", project)
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	// Create AWS client
	profile := cfg.AWS.Profile
	if profile == "" {
		profile = defaultAWSProfile
	}

	awsClient, err := awsInternal.NewClient(api.ctx, awsInternal.Options{
		Region:  cfg.AWS.Region,
		Profile: profile,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS client: %w", err)
	}

	s3Service := awsInternal.NewS3Service(awsClient)

	// Check for start requests
	requests, err := s3Service.CheckStartRequests(api.ctx, syncConfig.Bucket, project)
	if err != nil {
		return nil, fmt.Errorf("failed to check start requests: %w", err)
	}

	// Convert to API format
	var apiRequests []*StudentStartRequest
	for username, request := range requests {
		apiRequest := &StudentStartRequest{
			Username:    username,
			StudentID:   request.StudentID,
			Token:       request.Token,
			RequestedAt: request.RequestedAt,
			MachineHash: request.MachineHash,
		}
		apiRequests = append(apiRequests, apiRequest)
	}

	return apiRequests, nil
}

// StudentStartRequest represents a student's request to start their instance
type StudentStartRequest struct {
	Username    string    `json:"username"`
	StudentID   string    `json:"student_id"`
	Token       string    `json:"token"`
	RequestedAt time.Time `json:"requested_at"`
	MachineHash string    `json:"machine_hash"`
}

// SubmitStartRequest submits a start request for a student
func (api *S3StatusAPI) SubmitStartRequest(username, project, studentID, token, machineHash string) error {
	// Check if S3 sync is enabled
	syncConfig := lfrutils.GetS3SyncConfig()
	if !syncConfig.Enabled {
		return fmt.Errorf("S3 sync not enabled for project %s", project)
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Create AWS client
	profile := cfg.AWS.Profile
	if profile == "" {
		profile = defaultAWSProfile
	}

	awsClient, err := awsInternal.NewClient(api.ctx, awsInternal.Options{
		Region:  cfg.AWS.Region,
		Profile: profile,
	})
	if err != nil {
		return fmt.Errorf("failed to create AWS client: %w", err)
	}

	s3Service := awsInternal.NewS3Service(awsClient)

	// Create start request
	request := &awsInternal.StudentStartRequest{
		Username:    username,
		StudentID:   studentID,
		Token:       token,
		RequestedAt: time.Now(),
		MachineHash: machineHash,
	}

	// Submit to S3
	return s3Service.SubmitStartRequest(api.ctx, syncConfig.Bucket, project, request)
}

// ProcessStartRequests processes pending start requests and starts instances
func (api *S3StatusAPI) ProcessStartRequests(project string, autoApprove bool) (int, error) {
	// Get pending requests
	requests, err := api.CheckStartRequests(project)
	if err != nil {
		return 0, fmt.Errorf("failed to check start requests: %w", err)
	}

	if len(requests) == 0 {
		return 0, nil
	}

	// Process requests
	instanceAPI := NewInstanceAPI()
	processedCount := 0

	for _, request := range requests {
		if autoApprove {
			// Find user's instance and start it
			instances, err := instanceAPI.ListInstances(project)
			if err != nil {
				log.Printf("Failed to list instances for %s: %v", request.Username, err)
				continue
			}

			for _, instance := range instances {
				if instance.Username == request.Username {
					err = instanceAPI.StartInstance(instance.Name)
					if err != nil {
						log.Printf("Failed to start instance for %s: %v", request.Username, err)
					} else {
						processedCount++
						// Update status
						_ = api.UpdateStudentStatus(request.Username, project)
					}
					break
				}
			}
		}
	}

	return processedCount, nil
}

// EnableS3SyncForProject enables S3 status synchronization for a project
func (api *S3StatusAPI) EnableS3SyncForProject(project, bucket string) error {
	return lfrutils.EnableS3Sync(project, bucket)
}

// UpdateAllStudentStatuses updates S3 status for all students in a project
func (api *S3StatusAPI) UpdateAllStudentStatuses(project string) error {
	// Get all instances
	instanceAPI := NewInstanceAPI()
	instances, err := instanceAPI.ListInstances(project)
	if err != nil {
		return fmt.Errorf("failed to get instances: %w", err)
	}

	// Update status for each student
	var errors []string
	for _, instance := range instances {
		if instance.Username != "" {
			err = api.UpdateStudentStatus(instance.Username, project)
			if err != nil {
				errors = append(errors, fmt.Sprintf("%s: %v", instance.Username, err))
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to update some statuses: %s", strings.Join(errors, "; "))
	}

	return nil
}