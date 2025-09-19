// Package services provides SSH proxy services for embedded terminal
package services

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/scttfrdmn/lfr-tools/pkg/api"
)

// SSHProxyService provides WebSocket SSH proxy functionality
type SSHProxyService struct {
	instanceAPI  *api.InstanceAPI
	connections  map[string]*SSHConnection
	mu           sync.RWMutex
}

// SSHConnection represents an active SSH connection
type SSHConnection struct {
	ID           string
	Username     string
	Project      string
	InstanceName string
	PublicIP     string
	Client       *ssh.Client
	Session      *ssh.Session
	Connected    bool
	CreatedAt    time.Time
}

// SSHConnectionInfo represents connection information for frontend
type SSHConnectionInfo struct {
	ID        string `json:"id"`
	Connected bool   `json:"connected"`
	PublicIP  string `json:"public_ip"`
	Username  string `json:"username"`
}

// NewSSHProxyService creates a new SSH proxy service
func NewSSHProxyService() *SSHProxyService {
	return &SSHProxyService{
		instanceAPI: api.NewInstanceAPI(),
		connections: make(map[string]*SSHConnection),
	}
}

// InitiateSSHConnection starts an SSH connection to a user's instance
func (s *SSHProxyService) InitiateSSHConnection(username, project string) (*SSHConnectionInfo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Get user's instance
	instances, err := s.instanceAPI.ListInstances(project)
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
		return nil, fmt.Errorf("instance is %s, not running", targetInstance.State)
	}

	if targetInstance.PublicIP == "" {
		return nil, fmt.Errorf("instance has no public IP")
	}

	// Create connection ID
	connectionID := fmt.Sprintf("%s-%s-%d", project, username, time.Now().Unix())

	// Get SSH key
	keyPath, err := s.getSSHKeyPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get SSH key: %w", err)
	}

	// Create SSH client configuration
	sshConfig, err := s.createSSHConfig(keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create SSH config: %w", err)
	}

	// Establish SSH connection
	client, err := ssh.Dial("tcp", targetInstance.PublicIP+":22", sshConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect via SSH: %w", err)
	}

	// Create SSH session
	session, err := client.NewSession()
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to create SSH session: %w", err)
	}

	// Store connection
	connection := &SSHConnection{
		ID:           connectionID,
		Username:     username,
		Project:      project,
		InstanceName: targetInstance.Name,
		PublicIP:     targetInstance.PublicIP,
		Client:       client,
		Session:      session,
		Connected:    true,
		CreatedAt:    time.Now(),
	}

	s.connections[connectionID] = connection

	log.Printf("SSH connection established: %s@%s", username, targetInstance.PublicIP)

	return &SSHConnectionInfo{
		ID:        connectionID,
		Connected: true,
		PublicIP:  targetInstance.PublicIP,
		Username:  username,
	}, nil
}

// ExecuteCommand executes a command via SSH and returns the output
func (s *SSHProxyService) ExecuteCommand(connectionID, command string) (string, error) {
	s.mu.RLock()
	connection, exists := s.connections[connectionID]
	s.mu.RUnlock()

	if !exists {
		return "", fmt.Errorf("connection not found: %s", connectionID)
	}

	if !connection.Connected {
		return "", fmt.Errorf("connection is not active")
	}

	// Create new session for command execution
	session, err := connection.Client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	// Execute command and capture output
	output, err := session.CombinedOutput(command)
	if err != nil {
		return string(output), fmt.Errorf("command failed: %w", err)
	}

	return string(output), nil
}

// StartInteractiveSession starts an interactive SSH session
func (s *SSHProxyService) StartInteractiveSession(connectionID string) error {
	s.mu.RLock()
	connection, exists := s.connections[connectionID]
	s.mu.RUnlock()

	if !exists {
		return fmt.Errorf("connection not found: %s", connectionID)
	}

	// Set up pseudo-terminal
	err := connection.Session.RequestPty("xterm-256color", 80, 24, ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	})
	if err != nil {
		return fmt.Errorf("failed to request pty: %w", err)
	}

	// Start shell
	err = connection.Session.Shell()
	if err != nil {
		return fmt.Errorf("failed to start shell: %w", err)
	}

	return nil
}

// CloseConnection closes an SSH connection
func (s *SSHProxyService) CloseConnection(connectionID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	connection, exists := s.connections[connectionID]
	if !exists {
		return fmt.Errorf("connection not found: %s", connectionID)
	}

	// Close SSH session and client
	if connection.Session != nil {
		connection.Session.Close()
	}
	if connection.Client != nil {
		connection.Client.Close()
	}

	connection.Connected = false
	delete(s.connections, connectionID)

	log.Printf("SSH connection closed: %s", connectionID)
	return nil
}

// ListConnections returns all active SSH connections
func (s *SSHProxyService) ListConnections() ([]*SSHConnectionInfo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var connections []*SSHConnectionInfo
	for _, conn := range s.connections {
		connections = append(connections, &SSHConnectionInfo{
			ID:        conn.ID,
			Connected: conn.Connected,
			PublicIP:  conn.PublicIP,
			Username:  conn.Username,
		})
	}

	return connections, nil
}

// getSSHKeyPath gets the path to the SSH private key
func (s *SSHProxyService) getSSHKeyPath() (string, error) {
	// Try default LFR Tools key location
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	keyPath := fmt.Sprintf("%s/.ssh/lfr-tools/LightsailDefaultKey.pem", homeDir)

	// Check if key exists
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		return "", fmt.Errorf("SSH key not found at %s. Run 'lfr ssh keys download' first", keyPath)
	}

	return keyPath, nil
}

// createSSHConfig creates SSH client configuration
func (s *SSHProxyService) createSSHConfig(keyPath string) (*ssh.ClientConfig, error) {
	// Read private key
	keyData, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read SSH key: %w", err)
	}

	// Parse private key
	signer, err := ssh.ParsePrivateKey(keyData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse SSH key: %w", err)
	}

	// Create SSH client configuration
	config := &ssh.ClientConfig{
		User: "ubuntu", // Default user for LfR instances
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // For educational use
		Timeout:         10 * time.Second,
	}

	return config, nil
}

// StartInstanceIfNeeded starts an instance if it's not running
func (s *SSHProxyService) StartInstanceIfNeeded(username, project string) (*api.InstanceInfo, error) {
	instances, err := s.instanceAPI.ListInstances(project)
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

	// Start instance if not running
	if targetInstance.State != "running" {
		log.Printf("Starting instance %s for user %s", targetInstance.Name, username)

		err = s.instanceAPI.StartInstance(targetInstance.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to start instance: %w", err)
		}

		// Wait for instance to be running (simplified)
		for i := 0; i < 60; i++ { // Wait up to 5 minutes
			time.Sleep(5 * time.Second)

			// Refresh instance status
			instances, err = s.instanceAPI.ListInstances(project)
			if err != nil {
				continue
			}

			for _, instance := range instances {
				if instance.Username == username {
					if instance.State == "running" && instance.PublicIP != "" {
						targetInstance = instance
						log.Printf("Instance %s is now running with IP %s", instance.Name, instance.PublicIP)
						return targetInstance, nil
					}
					break
				}
			}
		}

		return nil, fmt.Errorf("instance failed to start within timeout")
	}

	return targetInstance, nil
}

// ConnectStudent provides simplified connection for students
func (s *SSHProxyService) ConnectStudent(username, project string) (*SSHConnectionInfo, error) {
	// First, ensure instance is running
	instance, err := s.StartInstanceIfNeeded(username, project)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare instance: %w", err)
	}

	// Then establish SSH connection
	return s.InitiateSSHConnection(username, project)
}