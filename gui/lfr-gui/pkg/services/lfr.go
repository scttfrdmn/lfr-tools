// Package services provides GUI backend services for LFR Tools
package services

import (
	"fmt"

	"github.com/scttfrdmn/lfr-tools/pkg/api"
)

// LFRService provides the main service interface for the GUI
type LFRService struct {
	instanceAPI *api.InstanceAPI
}

// NewLFRService creates a new LFR service instance
func NewLFRService() *LFRService {
	return &LFRService{
		instanceAPI: api.NewInstanceAPI(),
	}
}

// Type aliases for frontend compatibility
type UserInfo = api.UserInfo
type InstanceInfo = api.InstanceInfo
type ProjectInfo = api.ProjectInfo

// GetUserRole delegates to the API
func (s *LFRService) GetUserRole() (*UserInfo, error) {
	return s.instanceAPI.GetUserRole()
}

// ListInstances delegates to the API
func (s *LFRService) ListInstances(project string) ([]*InstanceInfo, error) {
	return s.instanceAPI.ListInstances(project)
}

// StartInstance delegates to the API
func (s *LFRService) StartInstance(instanceName string) error {
	return s.instanceAPI.StartInstance(instanceName)
}

// StopInstance delegates to the API
func (s *LFRService) StopInstance(instanceName string) error {
	return s.instanceAPI.StopInstance(instanceName)
}

// GetProjectInfo delegates to the API
func (s *LFRService) GetProjectInfo(project string) (*ProjectInfo, error) {
	return s.instanceAPI.GetProjectInfo(project)
}

// ConnectToInstance returns SSH connection information
func (s *LFRService) ConnectToInstance(username, project string) (string, error) {
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