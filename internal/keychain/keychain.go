// Package keychain provides cross-platform secure credential storage.
package keychain

import (
	"fmt"
	"runtime"
)

// KeychainService provides secure credential storage.
type KeychainService interface {
	Store(service, account, secret string) error
	Retrieve(service, account string) (string, error)
	Delete(service, account string) error
	List(service string) ([]string, error)
}

// NewKeychainService creates a platform-appropriate keychain service.
func NewKeychainService() (KeychainService, error) {
	switch runtime.GOOS {
	case "darwin":
		return NewMacOSKeychain(), nil
	case "windows":
		return NewWindowsKeychain(), nil
	case "linux":
		return NewLinuxKeychain(), nil
	default:
		return NewFileKeychain(), nil // Fallback to file-based storage
	}
}

// TokenStore provides high-level token storage operations.
type TokenStore struct {
	keychain KeychainService
	service  string
}

// NewTokenStore creates a new token store.
func NewTokenStore() (*TokenStore, error) {
	keychain, err := NewKeychainService()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize keychain: %w", err)
	}

	return &TokenStore{
		keychain: keychain,
		service:  "lfr-tools",
	}, nil
}

// StoreToken stores a student token securely.
func (ts *TokenStore) StoreToken(project, username, tokenData string) error {
	account := fmt.Sprintf("%s-%s", project, username)
	return ts.keychain.Store(ts.service, account, tokenData)
}

// RetrieveToken retrieves a stored token.
func (ts *TokenStore) RetrieveToken(project, username string) (string, error) {
	account := fmt.Sprintf("%s-%s", project, username)
	return ts.keychain.Retrieve(ts.service, account)
}

// DeleteToken removes a stored token.
func (ts *TokenStore) DeleteToken(project, username string) error {
	account := fmt.Sprintf("%s-%s", project, username)
	return ts.keychain.Delete(ts.service, account)
}

// ListTokens lists all stored tokens.
func (ts *TokenStore) ListTokens() ([]string, error) {
	return ts.keychain.List(ts.service)
}