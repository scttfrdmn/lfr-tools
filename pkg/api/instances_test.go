package api

import (
	"testing"
)

func TestNewInstanceAPI(t *testing.T) {
	api := NewInstanceAPI()
	if api == nil {
		t.Fatal("expected non-nil instance API")
	}

	if api.ctx == nil {
		t.Error("expected non-nil context")
	}
}

func TestExtractUsernameFromInstanceName(t *testing.T) {
	tests := []struct {
		instanceName string
		expected     string
	}{
		{"alice-ubuntu_22_04", "alice"},
		{"bob-gpu", "bob"},
		{"charlie-test-instance", "charlie"},
		{"single", "single"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.instanceName, func(t *testing.T) {
			result := extractUsernameFromInstanceName(tt.instanceName)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestUserInfoStructure(t *testing.T) {
	userInfo := &UserInfo{
		Role:        "student",
		Username:    "alice",
		Project:     "cs101",
		Permissions: []string{"connect"},
	}

	if userInfo.Role != "student" {
		t.Errorf("expected role 'student', got %s", userInfo.Role)
	}

	if userInfo.Username != "alice" {
		t.Errorf("expected username 'alice', got %s", userInfo.Username)
	}

	if len(userInfo.Permissions) != 1 {
		t.Errorf("expected 1 permission, got %d", len(userInfo.Permissions))
	}
}

const testInstanceStateRunning = "running"

func TestInstanceInfoStructure(t *testing.T) {
	instance := &InstanceInfo{
		Name:      "alice-ubuntu",
		State:     testInstanceStateRunning,
		PublicIP:  "1.2.3.4",
		Blueprint: "ubuntu_22_04",
		Bundle:    "app_standard_xl_1_0",
		Username:  "alice",
	}

	if instance.Name != "alice-ubuntu" {
		t.Errorf("expected name 'alice-ubuntu', got %s", instance.Name)
	}

	if instance.State != testInstanceStateRunning {
		t.Errorf("expected state 'running', got %s", instance.State)
	}

	if instance.Username != "alice" {
		t.Errorf("expected username 'alice', got %s", instance.Username)
	}
}

func TestProjectInfoStructure(t *testing.T) {
	project := &ProjectInfo{
		Name:          "cs101",
		StudentCount:  25,
		RunningCount:  12,
		BudgetUsed:    340.50,
		BudgetTotal:   500.00,
		DaysRemaining: 45,
	}

	if project.Name != "cs101" {
		t.Errorf("expected name 'cs101', got %s", project.Name)
	}

	if project.StudentCount != 25 {
		t.Errorf("expected student count 25, got %d", project.StudentCount)
	}

	if project.RunningCount > project.StudentCount {
		t.Error("running count cannot exceed student count")
	}

	if project.BudgetUsed < 0 {
		t.Error("budget used cannot be negative")
	}

	if project.DaysRemaining < 0 {
		t.Error("days remaining cannot be negative")
	}
}

func TestGetUserRoleWithoutAWS(t *testing.T) {
	api := NewInstanceAPI()

	// This test will likely fail without AWS credentials, which is expected
	userInfo, err := api.GetUserRole()
	if err != nil {
		// Expected failure without proper AWS setup
		t.Logf("GetUserRole failed as expected without AWS config: %v", err)
		return
	}

	// If it succeeds, validate the response
	if userInfo == nil {
		t.Error("expected non-nil user info")
		return
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

	if userInfo.Username == "" {
		t.Error("expected non-empty username")
	}

	if len(userInfo.Permissions) == 0 {
		t.Error("expected at least one permission")
	}
}

func TestListInstancesWithoutAWS(t *testing.T) {
	api := NewInstanceAPI()

	// This test will likely fail without AWS credentials, which is expected
	instances, err := api.ListInstances("")
	if err != nil {
		// Expected failure without proper AWS setup
		t.Logf("ListInstances failed as expected without AWS config: %v", err)
		return
	}

	// If it succeeds, validate the response
	if instances == nil {
		t.Log("instances slice is nil (acceptable when AWS unavailable)")
		return
	}

	// Validate instance structure
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