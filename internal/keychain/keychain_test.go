package keychain

import (
	"os"
	"testing"
)

func TestNewKeychainService(t *testing.T) {
	service, err := NewKeychainService()
	if err != nil {
		t.Fatalf("failed to create keychain service: %v", err)
	}

	if service == nil {
		t.Fatal("expected non-nil keychain service")
	}
}

func TestNewTokenStore(t *testing.T) {
	store, err := NewTokenStore()
	if err != nil {
		t.Fatalf("failed to create token store: %v", err)
	}

	if store == nil {
		t.Fatal("expected non-nil token store")
	}

	if store.keychain == nil {
		t.Error("expected non-nil keychain")
	}

	if store.service == "" {
		t.Error("expected non-empty service name")
	}
}

func TestFileKeychainOperations(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "keychain-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Warning: failed to cleanup temp dir: %v", err)
		}
	}()

	// Create file keychain with temp directory
	keychain := &FileKeychain{
		storePath: tempDir,
	}

	// Test store operation
	err = keychain.Store("test-service", "test-account", "test-secret")
	if err != nil {
		t.Fatalf("failed to store secret: %v", err)
	}

	// Test retrieve operation
	secret, err := keychain.Retrieve("test-service", "test-account")
	if err != nil {
		t.Fatalf("failed to retrieve secret: %v", err)
	}

	if secret != "test-secret" {
		t.Errorf("expected 'test-secret', got %s", secret)
	}

	// Test list operation
	accounts, err := keychain.List("test-service")
	if err != nil {
		t.Fatalf("failed to list accounts: %v", err)
	}

	if len(accounts) != 1 {
		t.Errorf("expected 1 account, got %d", len(accounts))
	}

	if accounts[0] != "test-account" {
		t.Errorf("expected 'test-account', got %s", accounts[0])
	}

	// Test delete operation
	err = keychain.Delete("test-service", "test-account")
	if err != nil {
		t.Fatalf("failed to delete secret: %v", err)
	}

	// Verify deletion
	_, err = keychain.Retrieve("test-service", "test-account")
	if err == nil {
		t.Error("expected error when retrieving deleted secret")
	}
}

func TestTokenStoreOperations(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "token-store-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Warning: failed to cleanup temp dir: %v", err)
		}
	}()

	// Create token store with file keychain
	store := &TokenStore{
		keychain: &FileKeychain{storePath: tempDir},
		service:  "lfr-tools-test",
	}

	// Test store token
	err = store.StoreToken("cs101", "alice", "test-token-data")
	if err != nil {
		t.Fatalf("failed to store token: %v", err)
	}

	// Test retrieve token
	tokenData, err := store.RetrieveToken("cs101", "alice")
	if err != nil {
		t.Fatalf("failed to retrieve token: %v", err)
	}

	if tokenData != "test-token-data" {
		t.Errorf("expected 'test-token-data', got %s", tokenData)
	}

	// Test list tokens
	tokens, err := store.ListTokens()
	if err != nil {
		t.Fatalf("failed to list tokens: %v", err)
	}

	if len(tokens) != 1 {
		t.Errorf("expected 1 token, got %d", len(tokens))
	}

	// Test delete token
	err = store.DeleteToken("cs101", "alice")
	if err != nil {
		t.Fatalf("failed to delete token: %v", err)
	}

	// Verify deletion
	_, err = store.RetrieveToken("cs101", "alice")
	if err == nil {
		t.Error("expected error when retrieving deleted token")
	}
}

func TestKeychainServiceInterface(t *testing.T) {
	// Test that FileKeychain implements Service interface
	var _ Service = &FileKeychain{}

	// Test that we can create the interface
	service := NewFileKeychain()
	if service == nil {
		t.Error("expected non-nil service")
	}
}