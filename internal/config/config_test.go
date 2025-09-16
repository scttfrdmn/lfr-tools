package config

import (
	"os"
	"testing"
	"path/filepath"
)

func TestLoadDefaultConfig(t *testing.T) {
	config, err := Load()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Test default values
	if config.AWS.Profile != "default" {
		t.Errorf("expected default profile 'default', got %s", config.AWS.Profile)
	}

	if config.AWS.Region != "us-east-1" {
		t.Errorf("expected default region 'us-east-1', got %s", config.AWS.Region)
	}

	if config.Defaults.Blueprint != "ubuntu_22_04" {
		t.Errorf("expected default blueprint 'ubuntu_22_04', got %s", config.Defaults.Blueprint)
	}

	if config.Defaults.Bundle != "nano_2_0" {
		t.Errorf("expected default bundle 'nano_2_0', got %s", config.Defaults.Bundle)
	}

	if config.Defaults.IdleThreshold != 120 {
		t.Errorf("expected default idle threshold 120, got %d", config.Defaults.IdleThreshold)
	}
}

func TestExpandPath(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "absolute path unchanged",
			path:     "/usr/local/bin",
			expected: "/usr/local/bin",
		},
		{
			name:     "relative path unchanged",
			path:     "relative/path",
			expected: "relative/path",
		},
		{
			name:     "empty path unchanged",
			path:     "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := expandPath(tt.path)
			if tt.name != "tilde expansion" && result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestExpandPathWithTilde(t *testing.T) {
	homeDir := mustGetHomeDir()

	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "tilde only",
			path:     "~",
			expected: homeDir,
		},
		{
			name:     "tilde with path",
			path:     "~/.ssh/config",
			expected: filepath.Join(homeDir, ".ssh/config"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := expandPath(tt.path)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestMustGetHomeDir(t *testing.T) {
	// This should not panic
	homeDir := mustGetHomeDir()
	if homeDir == "" {
		t.Error("expected non-empty home directory")
	}

	// Verify it's a real directory
	if _, err := os.Stat(homeDir); os.IsNotExist(err) {
		t.Errorf("home directory does not exist: %s", homeDir)
	}
}

func TestConfigStructure(t *testing.T) {
	config := &Config{
		AWS: AWSConfig{
			Profile: "test-profile",
			Region:  "us-west-2",
		},
		Defaults: DefaultsConfig{
			Blueprint:     "test-blueprint",
			Bundle:        "test-bundle",
			IdleThreshold: 60,
		},
		SSH: SSHConfig{
			KeyPath:    "/test/keys",
			ConfigPath: "/test/config",
		},
		Debug: true,
	}

	// Verify all fields are set correctly
	if config.AWS.Profile != "test-profile" {
		t.Errorf("expected profile 'test-profile', got %s", config.AWS.Profile)
	}

	if config.AWS.Region != "us-west-2" {
		t.Errorf("expected region 'us-west-2', got %s", config.AWS.Region)
	}

	if config.Defaults.Blueprint != "test-blueprint" {
		t.Errorf("expected blueprint 'test-blueprint', got %s", config.Defaults.Blueprint)
	}

	if config.Defaults.Bundle != "test-bundle" {
		t.Errorf("expected bundle 'test-bundle', got %s", config.Defaults.Bundle)
	}

	if config.Defaults.IdleThreshold != 60 {
		t.Errorf("expected idle threshold 60, got %d", config.Defaults.IdleThreshold)
	}

	if config.SSH.KeyPath != "/test/keys" {
		t.Errorf("expected key path '/test/keys', got %s", config.SSH.KeyPath)
	}

	if config.SSH.ConfigPath != "/test/config" {
		t.Errorf("expected config path '/test/config', got %s", config.SSH.ConfigPath)
	}

	if !config.Debug {
		t.Error("expected debug to be true")
	}
}