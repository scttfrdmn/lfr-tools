// Package services provides DCV (NICE Desktop Cloud Visualization) integration
package services

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/scttfrdmn/lfr-tools/pkg/api"
)

// DCVService provides NICE DCV remote desktop functionality
type DCVService struct {
	instanceAPI *api.InstanceAPI
	connections map[string]*DCVConnection
	mu          sync.RWMutex
}

// DCVConnection represents an active DCV connection
type DCVConnection struct {
	ID           string
	Username     string
	Project      string
	InstanceName string
	PublicIP     string
	DCVPort      int
	SessionID    string
	Connected    bool
	CreatedAt    time.Time
}

// DCVConnectionInfo represents DCV connection information for frontend
type DCVConnectionInfo struct {
	ID          string `json:"id"`
	Connected   bool   `json:"connected"`
	PublicIP    string `json:"public_ip"`
	Username    string `json:"username"`
	SessionID   string `json:"session_id"`
	DCVPort     int    `json:"dcv_port"`
	ViewerURL   string `json:"viewer_url"`
}

// DCVQuality represents DCV connection quality settings
type DCVQuality string

const (
	DCVQualityLow      DCVQuality = "low"
	DCVQualityMedium   DCVQuality = "medium"
	DCVQualityHigh     DCVQuality = "high"
	DCVQualityLossless DCVQuality = "lossless"
)

// NewDCVService creates a new DCV service
func NewDCVService() *DCVService {
	return &DCVService{
		instanceAPI: api.NewInstanceAPI(),
		connections: make(map[string]*DCVConnection),
	}
}

// ConnectDCV initiates a DCV connection to a user's instance
func (d *DCVService) ConnectDCV(username, project string, quality DCVQuality) (*DCVConnectionInfo, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	// Get user's instance
	instances, err := d.instanceAPI.ListInstances(project)
	if err != nil {
		return nil, fmt.Errorf("failed to list instances: %w", err)
	}

	var targetInstance *api.InstanceInfo
	for _, instance := range instances {
		if instance.Username == username {
			targetInstance = instance
			break
		}
	}

	if targetInstance == nil {
		return nil, fmt.Errorf("no instance found for user: %s", username)
	}

	// Check if instance is running
	if targetInstance.State != "running" {
		return nil, fmt.Errorf("instance is %s, not running. Start the instance first", targetInstance.State)
	}

	if targetInstance.PublicIP == "" {
		return nil, fmt.Errorf("instance has no public IP")
	}

	// Create connection ID
	connectionID := fmt.Sprintf("dcv-%s-%s-%d", project, username, time.Now().Unix())

	// Configure DCV session
	sessionID := fmt.Sprintf("%s-session", username)
	dcvPort := 8443 // Standard DCV port

	// Check DCV server status on instance
	dcvStatus, err := d.checkDCVStatus(targetInstance.PublicIP, username)
	if err != nil {
		return nil, fmt.Errorf("failed to check DCV status: %w", err)
	}

	if !dcvStatus.Running {
		log.Printf("DCV server not running on %s, attempting to start", targetInstance.PublicIP)
		err = d.startDCVServer(targetInstance.PublicIP, username, sessionID)
		if err != nil {
			return nil, fmt.Errorf("failed to start DCV server: %w", err)
		}
	}

	// Create DCV connection
	connection := &DCVConnection{
		ID:           connectionID,
		Username:     username,
		Project:      project,
		InstanceName: targetInstance.Name,
		PublicIP:     targetInstance.PublicIP,
		DCVPort:      dcvPort,
		SessionID:    sessionID,
		Connected:    true,
		CreatedAt:    time.Now(),
	}

	d.connections[connectionID] = connection

	// Create viewer URL for web-based DCV access
	viewerURL := fmt.Sprintf("https://%s:%d/#%s", targetInstance.PublicIP, dcvPort, sessionID)

	log.Printf("DCV connection established: %s@%s:%d", username, targetInstance.PublicIP, dcvPort)

	return &DCVConnectionInfo{
		ID:        connectionID,
		Connected: true,
		PublicIP:  targetInstance.PublicIP,
		Username:  username,
		SessionID: sessionID,
		DCVPort:   dcvPort,
		ViewerURL: viewerURL,
	}, nil
}

// LaunchDCVViewer launches the appropriate DCV viewer for the platform
func (d *DCVService) LaunchDCVViewer(connectionID string, quality DCVQuality) error {
	d.mu.RLock()
	connection, exists := d.connections[connectionID]
	d.mu.RUnlock()

	if !exists {
		return fmt.Errorf("DCV connection not found: %s", connectionID)
	}

	// Launch platform-specific DCV viewer
	switch runtime.GOOS {
	case "darwin":
		return d.launchDCVViewerMacOS(connection, quality)
	case "windows":
		return d.launchDCVViewerWindows(connection, quality)
	case "linux":
		return d.launchDCVViewerLinux(connection, quality)
	default:
		return d.launchDCVViewerWeb(connection)
	}
}

// launchDCVViewerMacOS launches DCV viewer on macOS
func (d *DCVService) launchDCVViewerMacOS(connection *DCVConnection, quality DCVQuality) error {
	// Check if DCV viewer is installed
	dcvPath := "/Applications/DCV Viewer.app/Contents/MacOS/DCV Viewer"
	if _, err := exec.LookPath(dcvPath); err != nil {
		// Fall back to web viewer
		return d.launchDCVViewerWeb(connection)
	}

	// Build DCV connection URL
	dcvURL := fmt.Sprintf("dcv://%s:%d/%s", connection.PublicIP, connection.DCVPort, connection.SessionID)

	// Launch DCV viewer with quality settings
	qualityArgs := d.getDCVQualityArgs(quality)
	args := append([]string{dcvURL}, qualityArgs...)

	cmd := exec.Command(dcvPath, args...)
	return cmd.Start()
}

// launchDCVViewerWindows launches DCV viewer on Windows
func (d *DCVService) launchDCVViewerWindows(connection *DCVConnection, quality DCVQuality) error {
	// Check if DCV viewer is installed
	dcvPath := "C:\\Program Files\\NICE\\DCV\\Client\\bin\\dcvviewer.exe"
	if _, err := exec.LookPath(dcvPath); err != nil {
		// Fall back to web viewer
		return d.launchDCVViewerWeb(connection)
	}

	// Build DCV connection URL
	dcvURL := fmt.Sprintf("dcv://%s:%d/%s", connection.PublicIP, connection.DCVPort, connection.SessionID)

	// Launch DCV viewer
	qualityArgs := d.getDCVQualityArgs(quality)
	args := append([]string{dcvURL}, qualityArgs...)

	cmd := exec.Command(dcvPath, args...)
	return cmd.Start()
}

// launchDCVViewerLinux launches DCV viewer on Linux
func (d *DCVService) launchDCVViewerLinux(connection *DCVConnection, quality DCVQuality) error {
	// Check if DCV viewer is installed
	if _, err := exec.LookPath("dcvviewer"); err != nil {
		// Fall back to web viewer
		return d.launchDCVViewerWeb(connection)
	}

	// Build DCV connection URL
	dcvURL := fmt.Sprintf("dcv://%s:%d/%s", connection.PublicIP, connection.DCVPort, connection.SessionID)

	// Launch DCV viewer
	qualityArgs := d.getDCVQualityArgs(quality)
	args := append([]string{dcvURL}, qualityArgs...)

	cmd := exec.Command("dcvviewer", args...)
	return cmd.Start()
}

// launchDCVViewerWeb opens DCV in web browser
func (d *DCVService) launchDCVViewerWeb(connection *DCVConnection) error {
	viewerURL := fmt.Sprintf("https://%s:%d/#%s", connection.PublicIP, connection.DCVPort, connection.SessionID)

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", viewerURL)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", viewerURL)
	case "linux":
		cmd = exec.Command("xdg-open", viewerURL)
	default:
		return fmt.Errorf("unsupported platform for web browser launch")
	}

	return cmd.Start()
}

// getDCVQualityArgs returns DCV quality arguments
func (d *DCVService) getDCVQualityArgs(quality DCVQuality) []string {
	switch quality {
	case DCVQualityLow:
		return []string{"--quality", "low", "--fps", "15"}
	case DCVQualityMedium:
		return []string{"--quality", "medium", "--fps", "30"}
	case DCVQualityHigh:
		return []string{"--quality", "high", "--fps", "60"}
	case DCVQualityLossless:
		return []string{"--quality", "lossless", "--fps", "60"}
	default:
		return []string{"--quality", "medium", "--fps", "30"}
	}
}

// DCVStatus represents DCV server status
type DCVStatus struct {
	Running     bool   `json:"running"`
	Port        int    `json:"port"`
	Sessions    int    `json:"sessions"`
	Version     string `json:"version"`
	Error       string `json:"error,omitempty"`
}

// checkDCVStatus checks if DCV server is running on the instance
func (d *DCVService) checkDCVStatus(publicIP, username string) (*DCVStatus, error) {
	// This would normally SSH to the instance and check DCV status
	// For now, return a mock status

	// In production, this would execute:
	// ssh ubuntu@$publicIP "sudo systemctl status dcvserver && dcv list-sessions"

	return &DCVStatus{
		Running:  true, // Assume DCV is configured on LfR instances
		Port:     8443,
		Sessions: 1,
		Version:  "2023.1",
	}, nil
}

// startDCVServer starts the DCV server on the instance
func (d *DCVService) startDCVServer(publicIP, username, sessionID string) error {
	// This would normally SSH to the instance and start DCV
	// For now, return success

	// In production, this would execute:
	// ssh ubuntu@$publicIP "sudo systemctl start dcvserver && dcv create-session $sessionID --user ubuntu"

	log.Printf("Started DCV server on %s for session %s", publicIP, sessionID)
	return nil
}

// CloseDCVConnection closes a DCV connection
func (d *DCVService) CloseDCVConnection(connectionID string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	connection, exists := d.connections[connectionID]
	if !exists {
		return fmt.Errorf("DCV connection not found: %s", connectionID)
	}

	// Close DCV session on instance
	err := d.closeDCVSession(connection.PublicIP, connection.SessionID)
	if err != nil {
		log.Printf("Failed to close DCV session: %v", err)
		// Don't return error - still remove from our tracking
	}

	connection.Connected = false
	delete(d.connections, connectionID)

	log.Printf("DCV connection closed: %s", connectionID)
	return nil
}

// closeDCVSession closes the DCV session on the instance
func (d *DCVService) closeDCVSession(publicIP, sessionID string) error {
	// This would normally SSH to the instance and close the DCV session
	// ssh ubuntu@$publicIP "dcv close-session $sessionID"

	log.Printf("Closed DCV session %s on %s", sessionID, publicIP)
	return nil
}

// ListDCVConnections returns all active DCV connections
func (d *DCVService) ListDCVConnections() ([]*DCVConnectionInfo, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	var connections []*DCVConnectionInfo
	for _, conn := range d.connections {
		viewerURL := fmt.Sprintf("https://%s:%d/#%s", conn.PublicIP, conn.DCVPort, conn.SessionID)

		connections = append(connections, &DCVConnectionInfo{
			ID:        conn.ID,
			Connected: conn.Connected,
			PublicIP:  conn.PublicIP,
			Username:  conn.Username,
			SessionID: conn.SessionID,
			DCVPort:   conn.DCVPort,
			ViewerURL: viewerURL,
		})
	}

	return connections, nil
}

// GetDCVConnectionInfo returns information about a specific DCV connection
func (d *DCVService) GetDCVConnectionInfo(connectionID string) (*DCVConnectionInfo, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	connection, exists := d.connections[connectionID]
	if !exists {
		return nil, fmt.Errorf("DCV connection not found: %s", connectionID)
	}

	viewerURL := fmt.Sprintf("https://%s:%d/#%s", connection.PublicIP, connection.DCVPort, connection.SessionID)

	return &DCVConnectionInfo{
		ID:        connection.ID,
		Connected: connection.Connected,
		PublicIP:  connection.PublicIP,
		Username:  connection.Username,
		SessionID: connection.SessionID,
		DCVPort:   connection.DCVPort,
		ViewerURL: viewerURL,
	}, nil
}