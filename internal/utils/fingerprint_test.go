package utils

import (
	"testing"
)

func TestGenerateMachineFingerprint(t *testing.T) {
	fingerprint, err := GenerateMachineFingerprint()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if fingerprint == nil {
		t.Fatal("expected non-nil fingerprint")
	}

	if fingerprint.Hash == "" {
		t.Error("expected non-empty hash")
	}

	if fingerprint.Platform == "" {
		t.Error("expected non-empty platform")
	}

	if fingerprint.Hostname == "" {
		t.Error("expected non-empty hostname")
	}

	// Hash should be consistent
	fingerprint2, err := GenerateMachineFingerprint()
	if err != nil {
		t.Fatalf("expected no error on second generation, got %v", err)
	}

	if fingerprint.Hash != fingerprint2.Hash {
		t.Error("expected consistent hash generation")
	}
}

func TestValidateMachineFingerprint(t *testing.T) {
	// Generate a fingerprint
	fingerprint, err := GenerateMachineFingerprint()
	if err != nil {
		t.Fatalf("failed to generate fingerprint: %v", err)
	}

	// Validate against itself (should pass)
	valid, err := ValidateMachineFingerprint(fingerprint)
	if err != nil {
		t.Fatalf("validation failed: %v", err)
	}

	if !valid {
		t.Error("expected fingerprint to validate against itself")
	}

	// Test with different fingerprint (should fail)
	differentFingerprint := &MachineFingerprint{
		Hash:     "different-hash",
		Platform: "different-platform",
		Hostname: "different-hostname",
	}

	valid, err = ValidateMachineFingerprint(differentFingerprint)
	if err != nil {
		t.Fatalf("validation failed: %v", err)
	}

	if valid {
		t.Error("expected validation to fail for different fingerprint")
	}
}

func TestGetPrimaryMACAddress(t *testing.T) {
	mac, err := getPrimaryMACAddress()
	if err != nil {
		// This might fail in some test environments, so just log
		t.Logf("MAC address detection failed (expected in some test environments): %v", err)
		return
	}

	if mac == "" {
		t.Error("expected non-empty MAC address")
	}

	// Basic format check (should contain colons for MAC address)
	if !contains(mac, ":") && !contains(mac, "-") {
		t.Errorf("MAC address format seems invalid: %s", mac)
	}
}

func TestGetPlatformSpecificID(t *testing.T) {
	id, err := getPlatformSpecificID()

	// This might fail on some platforms, which is okay
	if err != nil {
		t.Logf("Platform-specific ID failed (expected on some platforms): %v", err)
		return
	}

	if id == "" {
		t.Error("expected non-empty platform ID")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
		 len(substr) == 0 ||
		 (len(substr) <= len(s) && containsHelper(s, substr)))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}