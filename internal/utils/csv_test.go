package utils

import (
	"os"
	"strings"
	"testing"
)

func TestParseUsersCSV(t *testing.T) {
	// Create temporary CSV file
	csvContent := `username,project,blueprint,bundle,groups
alice,test-project,ubuntu_22_04,app_standard_xl_1_0,group1;group2
bob,test-project,ubuntu_22_04,app_standard_2xl_1_0,group1
charlie,ml-project,ubuntu_22_04,gpu_nvidia_xl_1_0,ml-group`

	tmpFile, err := os.CreateTemp("", "test-users-*.csv")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(csvContent)
	if err != nil {
		t.Fatalf("failed to write CSV content: %v", err)
	}
	tmpFile.Close()

	// Parse the CSV
	users, err := ParseUsersCSV(tmpFile.Name())
	if err != nil {
		t.Fatalf("failed to parse CSV: %v", err)
	}

	if len(users) != 3 {
		t.Errorf("expected 3 users, got %d", len(users))
	}

	// Verify first user
	alice := users[0]
	if alice.Username != "alice" {
		t.Errorf("expected username 'alice', got %s", alice.Username)
	}

	if alice.Project != "test-project" {
		t.Errorf("expected project 'test-project', got %s", alice.Project)
	}

	if len(alice.Groups) != 2 {
		t.Errorf("expected 2 groups for alice, got %d", len(alice.Groups))
	}

	expectedGroups := map[string]bool{"group1": true, "group2": true}
	for _, group := range alice.Groups {
		if !expectedGroups[group] {
			t.Errorf("unexpected group for alice: %s", group)
		}
	}

	// Verify second user (single group)
	bob := users[1]
	if len(bob.Groups) != 1 || bob.Groups[0] != "group1" {
		t.Errorf("expected bob to have single group 'group1', got %v", bob.Groups)
	}

	// Verify third user (different project)
	charlie := users[2]
	if charlie.Project != "ml-project" {
		t.Errorf("expected charlie's project to be 'ml-project', got %s", charlie.Project)
	}
}

func TestParseUsersCSVErrors(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		errorMsg string
	}{
		{
			name:     "missing required column",
			content:  "username,project\nalice,test",
			errorMsg: "required column 'blueprint' not found",
		},
		{
			name:     "empty username",
			content:  "username,project,blueprint,bundle\n,test-project,ubuntu_22_04,app_standard_xl_1_0",
			errorMsg: "username cannot be empty",
		},
		{
			name:     "empty project",
			content:  "username,project,blueprint,bundle\nalice,,ubuntu_22_04,app_standard_xl_1_0",
			errorMsg: "project cannot be empty",
		},
		{
			name:     "mismatched columns",
			content:  "username,project,blueprint,bundle\nalice,test-project,ubuntu_22_04",
			errorMsg: "wrong number of fields",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpFile, err := os.CreateTemp("", "test-users-*.csv")
			if err != nil {
				t.Fatalf("failed to create temp file: %v", err)
			}
			defer os.Remove(tmpFile.Name())

			_, err = tmpFile.WriteString(tt.content)
			if err != nil {
				t.Fatalf("failed to write CSV content: %v", err)
			}
			tmpFile.Close()

			_, err = ParseUsersCSV(tmpFile.Name())
			if err == nil {
				t.Error("expected error but got none")
			} else if !strings.Contains(err.Error(), tt.errorMsg) {
				t.Errorf("expected error to contain '%s', got: %v", tt.errorMsg, err)
			}
		})
	}
}

func TestGenerateUsersCSVTemplate(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-template-*.csv")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	err = GenerateUsersCSVTemplate(tmpFile.Name())
	if err != nil {
		t.Fatalf("failed to generate template: %v", err)
	}

	// Verify file was created and has content
	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("failed to read generated template: %v", err)
	}

	contentStr := string(content)

	// Check for required headers
	if !strings.Contains(contentStr, "username") {
		t.Error("template should contain 'username' header")
	}

	if !strings.Contains(contentStr, "project") {
		t.Error("template should contain 'project' header")
	}

	if !strings.Contains(contentStr, "blueprint") {
		t.Error("template should contain 'blueprint' header")
	}

	if !strings.Contains(contentStr, "bundle") {
		t.Error("template should contain 'bundle' header")
	}

	// Check for sample data
	if !strings.Contains(contentStr, "alice") {
		t.Error("template should contain sample user 'alice'")
	}
}

func TestParseGroupsCSV(t *testing.T) {
	csvContent := `name,description,policies,project
researchers,Research team,arn:aws:iam::aws:policy/ReadOnlyAccess,research-project
admins,Admin team,arn:aws:iam::aws:policy/PowerUserAccess;arn:aws:iam::aws:policy/IAMReadOnlyAccess,admin-project`

	tmpFile, err := os.CreateTemp("", "test-groups-*.csv")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(csvContent)
	if err != nil {
		t.Fatalf("failed to write CSV content: %v", err)
	}
	tmpFile.Close()

	groups, err := ParseGroupsCSV(tmpFile.Name())
	if err != nil {
		t.Fatalf("failed to parse groups CSV: %v", err)
	}

	if len(groups) != 2 {
		t.Errorf("expected 2 groups, got %d", len(groups))
	}

	// Verify first group
	researchers := groups[0]
	if researchers.Name != "researchers" {
		t.Errorf("expected name 'researchers', got %s", researchers.Name)
	}

	if len(researchers.Policies) != 1 {
		t.Errorf("expected 1 policy for researchers, got %d", len(researchers.Policies))
	}

	// Verify second group (multiple policies)
	admins := groups[1]
	if len(admins.Policies) != 2 {
		t.Errorf("expected 2 policies for admins, got %d", len(admins.Policies))
	}
}