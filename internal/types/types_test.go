package types

import (
	"testing"
	"time"
)

func TestProjectValidation(t *testing.T) {
	tests := []struct {
		name    string
		project Project
		valid   bool
	}{
		{
			name: "valid project",
			project: Project{
				Name:      "test-project",
				Blueprint: "ubuntu_22_04",
				Bundle:    "nano_2_0",
				Region:    "us-east-1",
				CreatedAt: time.Now(),
			},
			valid: true,
		},
		{
			name: "empty name",
			project: Project{
				Name:      "",
				Blueprint: "ubuntu_22_04",
				Bundle:    "nano_2_0",
				Region:    "us-east-1",
				CreatedAt: time.Now(),
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := tt.project.Name != ""
			if valid != tt.valid {
				t.Errorf("expected valid=%v, got valid=%v", tt.valid, valid)
			}
		})
	}
}

func TestUserValidation(t *testing.T) {
	tests := []struct {
		name  string
		user  User
		valid bool
	}{
		{
			name: "valid user",
			user: User{
				Username:     "test-user",
				Project:      "test-project",
				InstanceARN:  "arn:aws:lightsail:us-east-1:123456789012:Instance/test-instance",
				InstanceName: "test-user-ubuntu_22_04",
				CreatedAt:    time.Now(),
			},
			valid: true,
		},
		{
			name: "empty username",
			user: User{
				Username:     "",
				Project:      "test-project",
				InstanceARN:  "arn:aws:lightsail:us-east-1:123456789012:Instance/test-instance",
				InstanceName: "test-user-ubuntu_22_04",
				CreatedAt:    time.Now(),
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := tt.user.Username != ""
			if valid != tt.valid {
				t.Errorf("expected valid=%v, got valid=%v", tt.valid, valid)
			}
		})
	}
}

func TestInstanceStates(t *testing.T) {
	validStates := []string{"running", "stopped", "pending", "stopping", "starting", "terminating"}

	for _, state := range validStates {
		t.Run("state_"+state, func(t *testing.T) {
			instance := Instance{
				Name:  "test-instance",
				State: state,
			}

			if instance.State != state {
				t.Errorf("expected state %s, got %s", state, instance.State)
			}
		})
	}
}

func TestGroupWithPolicies(t *testing.T) {
	group := Group{
		Name:        "test-group",
		Policies:    []string{"policy1", "policy2"},
		Description: "Test group",
		CreatedAt:   time.Now(),
	}

	if len(group.Policies) != 2 {
		t.Errorf("expected 2 policies, got %d", len(group.Policies))
	}

	expectedPolicies := map[string]bool{"policy1": true, "policy2": true}
	for _, policy := range group.Policies {
		if !expectedPolicies[policy] {
			t.Errorf("unexpected policy: %s", policy)
		}
	}
}