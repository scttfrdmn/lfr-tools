package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

// Test the overall command structure without trying to execute individual commands
func TestCommandStructure(t *testing.T) {
	tests := []struct {
		name     string
		cmd      *cobra.Command
		expected string
	}{
		{"users command", usersCmd, "users"},
		{"groups command", groupsCmd, "groups"},
		{"instances command", instancesCmd, "instances"},
		{"ssh command", sshCmd, "ssh"},
		{"dcv command", dcvCmd, "dcv"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.cmd == nil {
				t.Errorf("%s should not be nil", tt.name)
				return
			}

			if tt.cmd.Use != tt.expected {
				t.Errorf("expected %s.Use to be '%s', got '%s'", tt.name, tt.expected, tt.cmd.Use)
			}
		})
	}
}

func TestCommandSubcommands(t *testing.T) {
	tests := []struct {
		parentCmd    *cobra.Command
		parentName   string
		expectedSubs []string
	}{
		{
			usersCmd, "users",
			[]string{"create", "remove", "list"},
		},
		{
			groupsCmd, "groups",
			[]string{"create", "remove", "list"},
		},
		{
			instancesCmd, "instances",
			[]string{"list", "start", "stop", "monitor"},
		},
		{
			sshCmd, "ssh",
			[]string{"keys", "config"},
		},
		{
			dcvCmd, "dcv",
			[]string{"config", "status", "sessions", "optimize"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.parentName+"_subcommands", func(t *testing.T) {
			if tt.parentCmd == nil {
				t.Errorf("%s command should not be nil", tt.parentName)
				return
			}

			for _, expectedSub := range tt.expectedSubs {
				found := false
				for _, cmd := range tt.parentCmd.Commands() {
					if cmd.Use == expectedSub {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected subcommand '%s' not found in %s", expectedSub, tt.parentName)
				}
			}
		})
	}
}

func TestSSHKeysSubcommands(t *testing.T) {
	expectedSubs := []string{"list"}
	for _, expectedSub := range expectedSubs {
		found := false
		for _, cmd := range sshKeysCmd.Commands() {
			if cmd.Use == expectedSub {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected subcommand '%s' not found in ssh keys", expectedSub)
		}
	}
}

func TestDCVSessionsSubcommands(t *testing.T) {
	expectedSubs := []string{"list"}
	for _, expectedSub := range expectedSubs {
		found := false
		for _, cmd := range dcvSessionsCmd.Commands() {
			if cmd.Use == expectedSub {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected subcommand '%s' not found in dcv sessions", expectedSub)
		}
	}
}

func TestCommandFlags(t *testing.T) {
	// Test that key commands have their expected flags
	tests := []struct {
		cmd      *cobra.Command
		cmdName  string
		flagName string
		required bool
	}{
		{usersCreateCmd, "users create", "project", true},
		{usersCreateCmd, "users create", "blueprint", true},
		{usersCreateCmd, "users create", "bundle", true},
		{usersCreateCmd, "users create", "region", true},
		{usersCreateCmd, "users create", "users", true},
		{groupsCreateCmd, "groups create", "name", true},
		{groupsCreateCmd, "groups create", "policies", true},
		{instancesListCmd, "instances list", "project", false},
		{sshConnectCmd, "ssh connect", "project", false},
		{dcvConnectCmd, "dcv connect", "quality", false},
	}

	for _, tt := range tests {
		t.Run(tt.cmdName+"_"+tt.flagName, func(t *testing.T) {
			if tt.cmd == nil {
				t.Errorf("%s command should not be nil", tt.cmdName)
				return
			}

			flag := tt.cmd.Flags().Lookup(tt.flagName)
			if flag == nil {
				t.Errorf("expected flag '%s' to exist on command '%s'", tt.flagName, tt.cmdName)
			}
		})
	}
}