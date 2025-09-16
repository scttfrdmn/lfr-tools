// Package utils provides common utility functions.
package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

const (
	upperCase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lowerCase = "abcdefghijklmnopqrstuvwxyz"
	digits    = "0123456789"
	allChars  = upperCase + lowerCase + digits
)

// GeneratePassword generates a secure password that meets AWS requirements.
// Based on the original script's password generation logic.
func GeneratePassword() (string, error) {
	// Generate two 8-character segments separated by a dash
	// This ensures we have uppercase, lowercase, and digits
	segment1, err := generateSegment(8)
	if err != nil {
		return "", fmt.Errorf("failed to generate first password segment: %w", err)
	}

	segment2, err := generateSegment(8)
	if err != nil {
		return "", fmt.Errorf("failed to generate second password segment: %w", err)
	}

	return segment1 + "-" + segment2, nil
}

// generateSegment generates a random string segment with mixed case and digits.
func generateSegment(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("length must be positive")
	}

	result := make([]byte, length)

	for i := 0; i < length; i++ {
		// Choose random character from all possible characters
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(allChars))))
		if err != nil {
			return "", fmt.Errorf("failed to generate random character: %w", err)
		}

		result[i] = allChars[randomIndex.Int64()]
	}

	return string(result), nil
}

// ValidatePassword checks if a password meets AWS requirements.
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}

	hasUpper := false
	hasLower := false
	hasDigit := false

	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasDigit = true
		}
	}

	if !hasUpper {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}
	if !hasDigit {
		return fmt.Errorf("password must contain at least one digit")
	}

	return nil
}