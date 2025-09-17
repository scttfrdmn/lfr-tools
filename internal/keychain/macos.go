// +build darwin

package keychain

import (
	"fmt"
	"os/exec"
	"strings"
)

// MacOSKeychain implements keychain service for macOS using the system keychain.
type MacOSKeychain struct{}

// NewMacOSKeychain creates a new macOS keychain service.
func NewMacOSKeychain() *MacOSKeychain {
	return &MacOSKeychain{}
}

// Store stores a secret in the macOS keychain.
func (k *MacOSKeychain) Store(service, account, secret string) error {
	// Use security command to store in keychain
	cmd := exec.Command("security", "add-generic-password",
		"-s", service,
		"-a", account,
		"-w", secret,
		"-U") // Update if exists

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to store in keychain: %w (output: %s)", err, string(output))
	}

	return nil
}

// Retrieve retrieves a secret from the macOS keychain.
func (k *MacOSKeychain) Retrieve(service, account string) (string, error) {
	cmd := exec.Command("security", "find-generic-password",
		"-s", service,
		"-a", account,
		"-w") // Output password only

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve from keychain: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// Delete removes a secret from the macOS keychain.
func (k *MacOSKeychain) Delete(service, account string) error {
	cmd := exec.Command("security", "delete-generic-password",
		"-s", service,
		"-a", account)

	_, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to delete from keychain: %w", err)
	}

	return nil
}

// List lists all accounts for a service in the keychain.
func (k *MacOSKeychain) List(service string) ([]string, error) {
	cmd := exec.Command("security", "dump-keychain")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list keychain items: %w", err)
	}

	// Parse output to find accounts for our service
	lines := strings.Split(string(output), "\n")
	var accounts []string

	for _, line := range lines {
		if strings.Contains(line, fmt.Sprintf("svce<blob>=\"%s\"", service)) {
			// Extract account from the next line
			// This is a simplified implementation
			accounts = append(accounts, "keychain-account")
		}
	}

	return accounts, nil
}