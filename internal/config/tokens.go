// Package config handles token management for educational access.
package config

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/scttfrdmn/lfr-tools/internal/utils"
)

// StudentToken represents a hardware-tied access token for students.
type StudentToken struct {
	Project         string                    `json:"project"`
	Username        string                    `json:"username"`
	StudentID       string                    `json:"student_id"`
	Role            string                    `json:"role"` // student, ta, professor
	Permissions     []string                  `json:"permissions"`
	S3Bucket        string                    `json:"s3_bucket"`
	Fingerprint     *utils.MachineFingerprint `json:"machine_fingerprint,omitempty"`
	SSHKeyData      string                    `json:"ssh_key_data"`
	CreatedAt       time.Time                 `json:"created_at"`
	ExpiresAt       time.Time                 `json:"expires_at"`
	AccessStartDate time.Time                 `json:"access_start_date,omitempty"`
	AccessEndDate   time.Time                 `json:"access_end_date,omitempty"`
	TokenHash       string                    `json:"token_hash"`
}

// TokenManager manages student access tokens.
type TokenManager struct {
	tokensDir string
}

// NewTokenManager creates a new token manager.
func NewTokenManager() (*TokenManager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	tokensDir := filepath.Join(homeDir, ".lfr-tools", "tokens")
	if err := os.MkdirAll(tokensDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create tokens directory: %w", err)
	}

	return &TokenManager{
		tokensDir: tokensDir,
	}, nil
}

// GenerateToken creates a new student access token.
func (tm *TokenManager) GenerateToken(project, username, studentID, role string, permissions []string, s3Bucket string, expiresAt time.Time) (string, *StudentToken, error) {
	// Generate random token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", nil, fmt.Errorf("failed to generate token: %w", err)
	}

	tokenString := fmt.Sprintf("%s-%s-%s", project, username, base64.URLEncoding.EncodeToString(tokenBytes)[:16])

	// Create token object
	token := &StudentToken{
		Project:     project,
		Username:    username,
		StudentID:   studentID,
		Role:        role,
		Permissions: permissions,
		S3Bucket:    s3Bucket,
		CreatedAt:   time.Now(),
		ExpiresAt:   expiresAt,
		TokenHash:   fmt.Sprintf("%x", sha256.Sum256([]byte(tokenString))),
	}

	return tokenString, token, nil
}

// ActivateToken binds a token to the current machine.
func (tm *TokenManager) ActivateToken(tokenString, studentID string) error {
	// Parse token to get basic info
	// This is a simplified version - real implementation would validate token format
	parts := splitTokenString(tokenString)
	if len(parts) < 3 {
		return fmt.Errorf("invalid token format")
	}

	project := parts[0]
	username := parts[1]

	// Generate machine fingerprint
	fingerprint, err := utils.GenerateMachineFingerprint()
	if err != nil {
		return fmt.Errorf("failed to generate machine fingerprint: %w", err)
	}

	// Create token file
	token := &StudentToken{
		Project:     project,
		Username:    username,
		StudentID:   studentID,
		Role:        "student",
		Permissions: []string{"connect"},
		Fingerprint: fingerprint,
		CreatedAt:   time.Now(),
		TokenHash:   fmt.Sprintf("%x", sha256.Sum256([]byte(tokenString))),
	}

	return tm.SaveToken(project, username, token)
}

// SaveToken saves a token to local storage.
func (tm *TokenManager) SaveToken(project, username string, token *StudentToken) error {
	filename := fmt.Sprintf("%s-%s.json", project, username)
	filepath := filepath.Join(tm.tokensDir, filename)

	data, err := json.MarshalIndent(token, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal token: %w", err)
	}

	if err := os.WriteFile(filepath, data, 0600); err != nil {
		return fmt.Errorf("failed to save token: %w", err)
	}

	return nil
}

// LoadToken loads a token from local storage.
func (tm *TokenManager) LoadToken(project, username string) (*StudentToken, error) {
	filename := fmt.Sprintf("%s-%s.json", project, username)
	filepath := filepath.Join(tm.tokensDir, filename)

	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read token file: %w", err)
	}

	var token StudentToken
	if err := json.Unmarshal(data, &token); err != nil {
		return nil, fmt.Errorf("failed to unmarshal token: %w", err)
	}

	return &token, nil
}

// ValidateToken validates a token and checks machine binding.
func (tm *TokenManager) ValidateToken(project, username string) error {
	token, err := tm.LoadToken(project, username)
	if err != nil {
		return fmt.Errorf("token not found: %w", err)
	}

	// Check expiration
	if time.Now().After(token.ExpiresAt) {
		return fmt.Errorf("token expired on %s", token.ExpiresAt.Format("2006-01-02"))
	}

	// Check access window
	now := time.Now()
	if !token.AccessStartDate.IsZero() && now.Before(token.AccessStartDate) {
		return fmt.Errorf("access not yet available (starts %s)", token.AccessStartDate.Format("2006-01-02 15:04"))
	}
	if !token.AccessEndDate.IsZero() && now.After(token.AccessEndDate) {
		return fmt.Errorf("access expired on %s", token.AccessEndDate.Format("2006-01-02 15:04"))
	}

	// Check machine fingerprint (if bound)
	if token.Fingerprint != nil {
		valid, err := utils.ValidateMachineFingerprint(token.Fingerprint)
		if err != nil {
			return fmt.Errorf("failed to validate machine fingerprint: %w", err)
		}
		if !valid {
			return fmt.Errorf("token is bound to a different machine")
		}
	}

	return nil
}

// ListTokens lists all stored tokens.
func (tm *TokenManager) ListTokens() ([]*StudentToken, error) {
	files, err := os.ReadDir(tm.tokensDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read tokens directory: %w", err)
	}

	var tokens []*StudentToken
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		data, err := os.ReadFile(filepath.Join(tm.tokensDir, file.Name()))
		if err != nil {
			continue // Skip unreadable files
		}

		var token StudentToken
		if err := json.Unmarshal(data, &token); err != nil {
			continue // Skip invalid tokens
		}

		tokens = append(tokens, &token)
	}

	return tokens, nil
}

// splitTokenString splits a token string into components.
func splitTokenString(tokenString string) []string {
	return strings.Split(tokenString, "-")
}