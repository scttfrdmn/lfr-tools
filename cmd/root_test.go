package cmd

import (
	"testing"
)

func TestRootCmdExists(t *testing.T) {
	if rootCmd == nil {
		t.Error("rootCmd should not be nil")
	}

	if rootCmd.Use != "lfr" {
		t.Errorf("expected rootCmd.Use to be 'lfr', got %s", rootCmd.Use)
	}
}

func TestVersionInfo(t *testing.T) {
	// Test that version variables exist and can be set
	originalVersion := version
	originalCommit := commit
	originalDate := date

	version = "v1.0.0"
	commit = "abc123"
	date = "2023-01-01T00:00:00Z"

	if version != "v1.0.0" {
		t.Errorf("expected version 'v1.0.0', got %s", version)
	}

	if commit != "abc123" {
		t.Errorf("expected commit 'abc123', got %s", commit)
	}

	if date != "2023-01-01T00:00:00Z" {
		t.Errorf("expected date '2023-01-01T00:00:00Z', got %s", date)
	}

	// Restore original values
	version = originalVersion
	commit = originalCommit
	date = originalDate
}

func TestInitConfig(t *testing.T) {
	// Test that initConfig doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("initConfig panicked: %v", r)
		}
	}()

	initConfig()
}