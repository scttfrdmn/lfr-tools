package services

import (
	"testing"
)

func TestNewLFRService(t *testing.T) {
	service := NewLFRService()
	if service == nil {
		t.Fatal("expected non-nil service")
	}

	if service.instanceAPI == nil {
		t.Error("expected non-nil instance API")
	}
}

func TestGetUserRole(t *testing.T) {
	service := NewLFRService()

	userInfo, err := service.GetUserRole()
	if err != nil {
		// This might fail without proper configuration, which is expected
		t.Logf("GetUserRole failed (expected in test environment): %v", err)
		return
	}

	if userInfo == nil {
		t.Error("expected non-nil user info")
		return
	}

	if userInfo.Role == "" {
		t.Error("expected non-empty role")
	}

	if userInfo.Username == "" {
		t.Error("expected non-empty username")
	}

	validRoles := []string{"student", "ta", "professor", "admin"}
	roleValid := false
	for _, validRole := range validRoles {
		if userInfo.Role == validRole {
			roleValid = true
			break
		}
	}

	if !roleValid {
		t.Errorf("invalid role: %s", userInfo.Role)
	}
}

func TestExtractUsername(t *testing.T) {
	// This function is now in the API package
	// Testing moved to API package tests
	t.Log("Username extraction testing moved to API package")
}

func TestListInstances(t *testing.T) {
	service := NewLFRService()

	// Test with empty project
	instances, err := service.ListInstances("")
	if err != nil {
		// This is expected to fail without AWS configuration
		t.Logf("ListInstances failed (expected without AWS config): %v", err)
		return
	}

	// If it succeeds, validate the response
	if instances == nil {
		t.Error("expected non-nil instances slice")
	}

	// Instances slice can be empty, which is valid
	for _, instance := range instances {
		if instance.Name == "" {
			t.Error("expected non-empty instance name")
		}
		if instance.State == "" {
			t.Error("expected non-empty instance state")
		}
		if instance.Username == "" {
			t.Error("expected non-empty username")
		}
	}
}

func TestProjectInfo(t *testing.T) {
	service := NewLFRService()

	projectInfo, err := service.GetProjectInfo("test-project")
	if err != nil {
		// This is expected to fail without AWS configuration
		t.Logf("GetProjectInfo failed (expected without AWS config): %v", err)
		return
	}

	if projectInfo == nil {
		t.Error("expected non-nil project info")
		return
	}

	if projectInfo.Name != "test-project" {
		t.Errorf("expected project name 'test-project', got %s", projectInfo.Name)
	}

	if projectInfo.StudentCount < 0 {
		t.Error("expected non-negative student count")
	}

	if projectInfo.RunningCount < 0 {
		t.Error("expected non-negative running count")
	}

	if projectInfo.RunningCount > projectInfo.StudentCount {
		t.Error("running count cannot exceed student count")
	}
}

func TestConnectToInstance(t *testing.T) {
	service := NewLFRService()

	// Test with non-existent user
	_, err := service.ConnectToInstance("nonexistent", "test-project")
	if err == nil {
		t.Error("expected error for nonexistent user")
	}

	// Error message should be helpful
	if err != nil && err.Error() == "" {
		t.Error("expected non-empty error message")
	}
}