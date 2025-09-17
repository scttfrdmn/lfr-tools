// Package utils provides hardware fingerprinting for non-transferable tokens.
package utils

import (
	"crypto/sha256"
	"fmt"
	"net"
	"os"
	"runtime"
	"strings"
)

// MachineFingerprint represents a unique machine identifier.
type MachineFingerprint struct {
	Hash     string `json:"hash"`
	Platform string `json:"platform"`
	Hostname string `json:"hostname"`
	Generated string `json:"generated"`
}

// GenerateMachineFingerprint creates a unique fingerprint for the current machine.
func GenerateMachineFingerprint() (*MachineFingerprint, error) {
	components := []string{}

	// Get hostname
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}
	components = append(components, "hostname:"+hostname)

	// Get primary MAC address
	macAddr, err := getPrimaryMACAddress()
	if err != nil {
		macAddr = "unknown"
	}
	components = append(components, "mac:"+macAddr)

	// Get platform info
	platform := runtime.GOOS + "-" + runtime.GOARCH
	components = append(components, "platform:"+platform)

	// Get user info (adds user-specific binding)
	userInfo, err := os.UserHomeDir()
	if err != nil {
		userInfo = "unknown"
	}
	components = append(components, "user:"+userInfo)

	// Additional platform-specific identifiers
	if platformID, err := getPlatformSpecificID(); err == nil {
		components = append(components, "platform-id:"+platformID)
	}

	// Create hash
	combined := strings.Join(components, "|")
	hash := sha256.Sum256([]byte(combined))
	hashStr := fmt.Sprintf("%x", hash)

	return &MachineFingerprint{
		Hash:     hashStr,
		Platform: platform,
		Hostname: hostname,
		Generated: fmt.Sprintf("%d", len(components)),
	}, nil
}

// ValidateMachineFingerprint checks if current machine matches the fingerprint.
func ValidateMachineFingerprint(expected *MachineFingerprint) (bool, error) {
	current, err := GenerateMachineFingerprint()
	if err != nil {
		return false, fmt.Errorf("failed to generate current fingerprint: %w", err)
	}

	return current.Hash == expected.Hash, nil
}

// getPrimaryMACAddress gets the MAC address of the primary network interface.
func getPrimaryMACAddress() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", fmt.Errorf("failed to get network interfaces: %w", err)
	}

	for _, iface := range interfaces {
		// Skip loopback and non-up interfaces
		if iface.Flags&net.FlagLoopback != 0 || iface.Flags&net.FlagUp == 0 {
			continue
		}

		// Skip virtual interfaces
		if strings.Contains(strings.ToLower(iface.Name), "docker") ||
			strings.Contains(strings.ToLower(iface.Name), "veth") ||
			strings.Contains(strings.ToLower(iface.Name), "br-") {
			continue
		}

		// Return first physical interface MAC
		if len(iface.HardwareAddr) > 0 {
			return iface.HardwareAddr.String(), nil
		}
	}

	return "", fmt.Errorf("no suitable network interface found")
}

// getPlatformSpecificID gets platform-specific machine identifiers.
func getPlatformSpecificID() (string, error) {
	switch runtime.GOOS {
	case "darwin":
		return getMacOSMachineID()
	case "windows":
		return getWindowsMachineID()
	case "linux":
		return getLinuxMachineID()
	default:
		return "", fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}

// getMacOSMachineID gets macOS-specific machine identifier.
func getMacOSMachineID() (string, error) {
	// Use hardware UUID (available via system_profiler)
	// This is a simplified version - real implementation would use system calls
	return "macos-placeholder", nil
}

// getWindowsMachineID gets Windows-specific machine identifier.
func getWindowsMachineID() (string, error) {
	// Use machine GUID or similar
	return "windows-placeholder", nil
}

// getLinuxMachineID gets Linux-specific machine identifier.
func getLinuxMachineID() (string, error) {
	// Try /etc/machine-id first
	if data, err := os.ReadFile("/etc/machine-id"); err == nil {
		return strings.TrimSpace(string(data)), nil
	}

	// Fallback to /var/lib/dbus/machine-id
	if data, err := os.ReadFile("/var/lib/dbus/machine-id"); err == nil {
		return strings.TrimSpace(string(data)), nil
	}

	return "linux-unknown", nil
}