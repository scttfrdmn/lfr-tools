package utils

import (
	"testing"
	"time"

	"github.com/scttfrdmn/lfr-tools/internal/types"
)

func TestExtractUsernameFromInstance(t *testing.T) {
	tests := []struct {
		instanceName string
		expected     string
	}{
		{"alice-ubuntu_22_04", "alice"},
		{"bob-gpu", "bob"},
		{"charlie-test-instance", "charlie"},
		{"single", "single"},
		{"", ""},
		{"no-username-here", "no"},
	}

	for _, tt := range tests {
		t.Run(tt.instanceName, func(t *testing.T) {
			result := ExtractUsernameFromInstance(tt.instanceName)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestUpdateInstanceStatusInS3(t *testing.T) {
	// Test with S3 sync disabled (should not error)
	instance := &types.Instance{
		Name:      "alice-ubuntu",
		State:     "running",
		PublicIP:  "1.2.3.4",
		Tags:      map[string]string{"Project": "test-project"},
		CreatedAt: time.Now(),
	}

	// This should not error even without S3 configuration
	err := UpdateInstanceStatusInS3(nil, instance)
	if err != nil {
		t.Errorf("expected no error when S3 sync disabled, got: %v", err)
	}
}

func TestGetS3SyncConfig(t *testing.T) {
	// Test getting S3 sync config (should not panic)
	config := GetS3SyncConfig()

	// Should return a valid config struct
	if config.Bucket == "" && config.Enabled {
		t.Error("if enabled is true, bucket should not be empty")
	}
}

func TestS3SyncConfigStructure(t *testing.T) {
	config := S3SyncConfig{
		Enabled:  true,
		Bucket:   "test-bucket",
		Project:  "test-project",
		AutoSync: true,
	}

	if !config.Enabled {
		t.Error("expected Enabled to be true")
	}

	if config.Bucket != "test-bucket" {
		t.Errorf("expected bucket 'test-bucket', got %s", config.Bucket)
	}

	if config.Project != "test-project" {
		t.Errorf("expected project 'test-project', got %s", config.Project)
	}

	if !config.AutoSync {
		t.Error("expected AutoSync to be true")
	}
}

func TestUpdateMultipleInstancesInS3(t *testing.T) {
	instances := []*types.Instance{
		{
			Name:      "alice-ubuntu",
			State:     "running",
			PublicIP:  "1.2.3.4",
			Tags:      map[string]string{"Project": "test-project"},
			CreatedAt: time.Now(),
		},
		{
			Name:      "bob-ubuntu",
			State:     "stopped",
			PublicIP:  "",
			Tags:      map[string]string{"Project": "test-project"},
			CreatedAt: time.Now(),
		},
	}

	// This should not panic or error
	UpdateMultipleInstancesInS3(nil, instances)

	// Function should complete without error even with disabled S3 sync
}