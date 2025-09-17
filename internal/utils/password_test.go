package utils

import (
	"fmt"
	"strings"
	"testing"
)

func TestGeneratePassword(t *testing.T) {
	password, err := GeneratePassword()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if password == "" {
		t.Fatal("expected non-empty password")
	}

	// Check format (should be segment1-segment2)
	parts := strings.Split(password, "-")
	if len(parts) != 2 {
		t.Errorf("expected password format 'segment1-segment2', got %s", password)
	}

	// Each segment should be 8 characters
	for i, part := range parts {
		if len(part) != 8 {
			t.Errorf("expected segment %d to be 8 characters, got %d: %s", i+1, len(part), part)
		}
	}

	// Password should meet AWS requirements
	err = ValidatePassword(password)
	if err != nil {
		t.Errorf("generated password failed validation: %v", err)
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		valid    bool
		errorMsg string
	}{
		{
			name:     "valid password",
			password: "Abc123def-Xyz789ghi",
			valid:    true,
		},
		{
			name:     "too short",
			password: "Abc123",
			valid:    false,
			errorMsg: "at least 8 characters",
		},
		{
			name:     "no uppercase",
			password: "abc123def456",
			valid:    false,
			errorMsg: "uppercase letter",
		},
		{
			name:     "no lowercase",
			password: "ABC123DEF456",
			valid:    false,
			errorMsg: "lowercase letter",
		},
		{
			name:     "no digits",
			password: "AbcDefGhiJkl",
			valid:    false,
			errorMsg: "digit",
		},
		{
			name:     "meets all requirements",
			password: "Password123",
			valid:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.password)

			if tt.valid {
				if err != nil {
					t.Errorf("expected password to be valid, got error: %v", err)
				}
			} else {
				if err == nil {
					t.Error("expected password to be invalid, got no error")
				} else if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error to contain '%s', got: %v", tt.errorMsg, err)
				}
			}
		})
	}
}

func TestGenerateSegment(t *testing.T) {
	tests := []struct {
		length int
		valid  bool
	}{
		{8, true},
		{16, true},
		{1, true},
		{0, false},
		{-1, false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("length_%d", tt.length), func(t *testing.T) {
			segment, err := generateSegment(tt.length)

			if tt.valid {
				if err != nil {
					t.Errorf("expected no error for length %d, got: %v", tt.length, err)
					return
				}

				if len(segment) != tt.length {
					t.Errorf("expected segment length %d, got %d: %s", tt.length, len(segment), segment)
				}

				// Check that all characters are from the expected set
				for _, char := range segment {
					if !strings.ContainsRune(allChars, char) {
						t.Errorf("unexpected character in segment: %c", char)
					}
				}
			} else {
				if err == nil {
					t.Errorf("expected error for invalid length %d", tt.length)
				}
			}
		})
	}
}

func TestPasswordUniqueness(t *testing.T) {
	// Generate multiple passwords and ensure they're unique
	passwords := make(map[string]bool)
	numPasswords := 100

	for i := 0; i < numPasswords; i++ {
		password, err := GeneratePassword()
		if err != nil {
			t.Fatalf("failed to generate password %d: %v", i, err)
		}

		if passwords[password] {
			t.Errorf("duplicate password generated: %s", password)
		}

		passwords[password] = true
	}

	if len(passwords) != numPasswords {
		t.Errorf("expected %d unique passwords, got %d", numPasswords, len(passwords))
	}
}