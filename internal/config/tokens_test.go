package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/scttfrdmn/lfr-tools/internal/utils"
)

func TestTokenManager(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "lfr-tokens-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create token manager with temp directory
	tm := &TokenManager{
		tokensDir: tempDir,
	}

	// Test token generation
	tokenString, token, err := tm.GenerateToken(
		"test-project",
		"alice",
		"12345",
		"student",
		[]string{"connect"},
		"test-bucket",
		time.Now().Add(24*time.Hour),
	)

	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	if tokenString == "" {
		t.Error("expected non-empty token string")
	}

	if token == nil {
		t.Fatal("expected non-nil token")
	}

	if token.Project != "test-project" {
		t.Errorf("expected project 'test-project', got %s", token.Project)
	}

	if token.Username != "alice" {
		t.Errorf("expected username 'alice', got %s", token.Username)
	}

	if token.StudentID != "12345" {
		t.Errorf("expected student ID '12345', got %s", token.StudentID)
	}

	if token.Role != "student" {
		t.Errorf("expected role 'student', got %s", token.Role)
	}

	// Test token saving
	err = tm.SaveToken("test-project", "alice", token)
	if err != nil {
		t.Fatalf("failed to save token: %v", err)
	}

	// Verify file was created
	tokenFile := filepath.Join(tempDir, "test-project-alice.json")
	if _, err := os.Stat(tokenFile); os.IsNotExist(err) {
		t.Error("token file was not created")
	}

	// Test token loading
	loadedToken, err := tm.LoadToken("test-project", "alice")
	if err != nil {
		t.Fatalf("failed to load token: %v", err)
	}

	if loadedToken.Project != token.Project {
		t.Errorf("loaded token project mismatch: expected %s, got %s", token.Project, loadedToken.Project)
	}

	if loadedToken.Username != token.Username {
		t.Errorf("loaded token username mismatch: expected %s, got %s", token.Username, loadedToken.Username)
	}
}

func TestActivateToken(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "lfr-tokens-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create token manager with temp directory
	tm := &TokenManager{
		tokensDir: tempDir,
	}

	// Test token activation
	tokenString := "test-project-alice-abc123"
	studentID := "12345"

	err = tm.ActivateToken(tokenString, studentID)
	if err != nil {
		t.Fatalf("failed to activate token: %v", err)
	}

	// Verify token was saved
	token, err := tm.LoadToken("test-project", "alice")
	if err != nil {
		t.Fatalf("failed to load activated token: %v", err)
	}

	if token.StudentID != studentID {
		t.Errorf("expected student ID %s, got %s", studentID, token.StudentID)
	}

	if token.Fingerprint == nil {
		t.Error("expected machine fingerprint to be set")
	}
}

func TestValidateToken(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "lfr-tokens-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create token manager with temp directory
	tm := &TokenManager{
		tokensDir: tempDir,
	}

	// Create a valid token
	fingerprint, err := utils.GenerateMachineFingerprint()
	if err != nil {
		t.Fatalf("failed to generate fingerprint: %v", err)
	}

	token := &StudentToken{
		Project:     "test-project",
		Username:    "alice",
		StudentID:   "12345",
		Role:        "student",
		Fingerprint: fingerprint,
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(24 * time.Hour),
	}

	// Save token
	err = tm.SaveToken("test-project", "alice", token)
	if err != nil {
		t.Fatalf("failed to save token: %v", err)
	}

	// Test validation (should pass)
	err = tm.ValidateToken("test-project", "alice")
	if err != nil {
		t.Errorf("expected token validation to pass, got error: %v", err)
	}

	// Test with expired token
	expiredToken := *token
	expiredToken.ExpiresAt = time.Now().Add(-1 * time.Hour)
	err = tm.SaveToken("test-project", "expired", &expiredToken)
	if err != nil {
		t.Fatalf("failed to save expired token: %v", err)
	}

	err = tm.ValidateToken("test-project", "expired")
	if err == nil {
		t.Error("expected validation to fail for expired token")
	}
}

func TestListTokens(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "lfr-tokens-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create token manager with temp directory
	tm := &TokenManager{
		tokensDir: tempDir,
	}

	// Create and save multiple tokens
	tokens := []*StudentToken{
		{
			Project:   "cs101",
			Username:  "alice",
			StudentID: "12345",
			Role:      "student",
			CreatedAt: time.Now(),
			ExpiresAt: time.Now().Add(24 * time.Hour),
		},
		{
			Project:   "cs101",
			Username:  "bob",
			StudentID: "67890",
			Role:      "student",
			CreatedAt: time.Now(),
			ExpiresAt: time.Now().Add(24 * time.Hour),
		},
	}

	for _, token := range tokens {
		err = tm.SaveToken(token.Project, token.Username, token)
		if err != nil {
			t.Fatalf("failed to save token for %s: %v", token.Username, err)
		}
	}

	// Test listing tokens
	listedTokens, err := tm.ListTokens()
	if err != nil {
		t.Fatalf("failed to list tokens: %v", err)
	}

	if len(listedTokens) != len(tokens) {
		t.Errorf("expected %d tokens, got %d", len(tokens), len(listedTokens))
	}

	// Verify tokens are correctly loaded
	for _, expectedToken := range tokens {
		found := false
		for _, listedToken := range listedTokens {
			if listedToken.Username == expectedToken.Username && listedToken.Project == expectedToken.Project {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("token for %s not found in list", expectedToken.Username)
		}
	}
}