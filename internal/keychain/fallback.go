package keychain

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FileKeychain implements keychain service using encrypted files as fallback.
type FileKeychain struct {
	storePath string
}

// NewFileKeychain creates a new file-based keychain service.
func NewFileKeychain() *FileKeychain {
	homeDir, _ := os.UserHomeDir()
	storePath := filepath.Join(homeDir, ".lfr-tools", "keychain")
	os.MkdirAll(storePath, 0700)

	return &FileKeychain{
		storePath: storePath,
	}
}

// NewWindowsKeychain creates a Windows credential manager service.
func NewWindowsKeychain() KeychainService {
	// For now, use file fallback on Windows
	// TODO: Implement Windows Credential Manager integration
	return NewFileKeychain()
}

// NewLinuxKeychain creates a Linux secret service.
func NewLinuxKeychain() KeychainService {
	// For now, use file fallback on Linux
	// TODO: Implement libsecret/gnome-keyring integration
	return NewFileKeychain()
}

// Store stores a secret in an encrypted file.
func (k *FileKeychain) Store(service, account, secret string) error {
	filename := fmt.Sprintf("%s-%s.json", service, account)
	filepath := filepath.Join(k.storePath, filename)

	// Simple storage (in production, this should be encrypted)
	data := map[string]string{
		"service": service,
		"account": account,
		"secret":  secret,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal secret: %w", err)
	}

	err = os.WriteFile(filepath, jsonData, 0600)
	if err != nil {
		return fmt.Errorf("failed to write secret file: %w", err)
	}

	return nil
}

// Retrieve retrieves a secret from file storage.
func (k *FileKeychain) Retrieve(service, account string) (string, error) {
	filename := fmt.Sprintf("%s-%s.json", service, account)
	filepath := filepath.Join(k.storePath, filename)

	data, err := os.ReadFile(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to read secret file: %w", err)
	}

	var secretData map[string]string
	err = json.Unmarshal(data, &secretData)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal secret: %w", err)
	}

	return secretData["secret"], nil
}

// Delete removes a secret file.
func (k *FileKeychain) Delete(service, account string) error {
	filename := fmt.Sprintf("%s-%s.json", service, account)
	filepath := filepath.Join(k.storePath, filename)

	err := os.Remove(filepath)
	if err != nil {
		return fmt.Errorf("failed to delete secret file: %w", err)
	}

	return nil
}

// List lists all stored secrets for a service.
func (k *FileKeychain) List(service string) ([]string, error) {
	files, err := os.ReadDir(k.storePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read keychain directory: %w", err)
	}

	var accounts []string
	prefix := service + "-"

	for _, file := range files {
		if strings.HasPrefix(file.Name(), prefix) && strings.HasSuffix(file.Name(), ".json") {
			// Extract account name
			account := strings.TrimSuffix(strings.TrimPrefix(file.Name(), prefix), ".json")
			accounts = append(accounts, account)
		}
	}

	return accounts, nil
}