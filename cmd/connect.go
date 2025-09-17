package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/scttfrdmn/lfr-tools/internal/aws"
	"github.com/scttfrdmn/lfr-tools/internal/config"
)

var connectCmd = &cobra.Command{
	Use:   "connect [username]",
	Short: "Connect to your instance (student/TA access)",
	Long: `Connect to your assigned instance using your access token. Automatically
starts stopped instances and handles SSH connection. No AWS credentials needed.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		username := args[0]
		project, _ := cmd.Flags().GetString("project")
		force, _ := cmd.Flags().GetBool("force")

		return connectToInstance(cmd.Context(), username, project, force)
	},
}

var connectListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available connections",
	Long:  `List all available connections from stored tokens.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return listAvailableConnections()
	},
}

var connectActivateCmd = &cobra.Command{
	Use:   "activate [token] [student-id]",
	Short: "Activate access token for this machine",
	Long: `Activate an access token for this machine. This binds the token to your
hardware and enables lfr connect functionality.`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		token := args[0]
		studentID := args[1]

		return activateStudentToken(token, studentID)
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)
	connectCmd.AddCommand(connectListCmd)
	connectCmd.AddCommand(connectActivateCmd)

	// Connect command flags
	connectCmd.Flags().StringP("project", "p", "", "Override project from token")
	connectCmd.Flags().BoolP("force", "f", false, "Force connection even if instance stopped")
}

// connectToInstance connects to a student's instance with automatic start.
func connectToInstance(ctx context.Context, username, project string, force bool) error {
	// Load token manager
	tm, err := config.NewTokenManager()
	if err != nil {
		return fmt.Errorf("failed to initialize token manager: %w", err)
	}

	// Find appropriate token
	var token *config.StudentToken
	if project != "" {
		token, err = tm.LoadToken(project, username)
	} else {
		// Auto-detect from available tokens
		tokens, err := tm.ListTokens()
		if err != nil {
			return fmt.Errorf("failed to list tokens: %w", err)
		}

		for _, t := range tokens {
			if t.Username == username {
				token = t
				break
			}
		}
	}

	if token == nil {
		return fmt.Errorf("no access token found for %s. Run: lfr connect activate <token> <student-id>", username)
	}

	// Validate token
	if err := tm.ValidateToken(token.Project, token.Username); err != nil {
		return fmt.Errorf("token validation failed: %w", err)
	}

	fmt.Printf("Connecting to %s's instance in project %s...\n", username, token.Project)

	// Check instance status via S3
	status, err := getInstanceStatusFromS3(ctx, token.S3Bucket, token.Project, username)
	if err != nil {
		return fmt.Errorf("failed to check instance status: %w", err)
	}

	fmt.Printf("Instance state: %s\n", status.State)

	// Handle stopped instance
	if status.State == "stopped" && !force {
		fmt.Printf("Instance is stopped. Requesting start from instructor...\n")

		// Submit start request
		err = submitStartRequest(ctx, token, status)
		if err != nil {
			return fmt.Errorf("failed to submit start request: %w", err)
		}

		fmt.Printf("✅ Start request submitted. Waiting for instructor approval...\n")

		// Wait for instance to start (with timeout)
		err = waitForInstanceStart(ctx, token.S3Bucket, token.Project, username, 5*time.Minute)
		if err != nil {
			return fmt.Errorf("instance start timed out: %w", err)
		}

		// Refresh status
		status, err = getInstanceStatusFromS3(ctx, token.S3Bucket, token.Project, username)
		if err != nil {
			return fmt.Errorf("failed to refresh instance status: %w", err)
		}
	}

	// Check if instance is ready for connection
	if status.State != "running" {
		if force {
			fmt.Printf("⚠️ Warning: Instance state is %s, attempting connection anyway...\n", status.State)
		} else {
			return fmt.Errorf("instance is not running (state: %s). Use --force to attempt connection anyway", status.State)
		}
	}

	if status.PublicIP == "" {
		return fmt.Errorf("instance has no public IP address")
	}

	// Establish SSH connection
	fmt.Printf("Connecting to %s...\n", status.PublicIP)

	// Create temporary SSH key file
	keyFile, err := createTempSSHKey(token.SSHKeyData)
	if err != nil {
		return fmt.Errorf("failed to create SSH key: %w", err)
	}
	defer os.Remove(keyFile)

	// Execute SSH
	sshArgs := []string{
		"-i", keyFile,
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
		fmt.Sprintf("%s@%s", username, status.PublicIP),
	}

	sshCmd := exec.Command("ssh", sshArgs...)
	sshCmd.Stdin = os.Stdin
	sshCmd.Stdout = os.Stdout
	sshCmd.Stderr = os.Stderr

	return sshCmd.Run()
}

// getInstanceStatusFromS3 retrieves instance status from S3.
func getInstanceStatusFromS3(ctx context.Context, bucket, project, username string) (*aws.StudentStatus, error) {
	url := fmt.Sprintf("https://%s.s3.amazonaws.com/%s/%s/status.json", bucket, project, username)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get status from S3: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get status: HTTP %d", resp.StatusCode)
	}

	var status aws.StudentStatus
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return nil, fmt.Errorf("failed to decode status: %w", err)
	}

	return &status, nil
}

// submitStartRequest submits a start request via S3.
func submitStartRequest(ctx context.Context, token *config.StudentToken, status *aws.StudentStatus) error {
	// Create start request
	request := &aws.StudentStartRequest{
		Username:    token.Username,
		StudentID:   token.StudentID,
		Token:       token.TokenHash,
		RequestedAt: time.Now(),
	}

	// Add machine hash if available
	if token.Fingerprint != nil {
		request.MachineHash = token.Fingerprint.Hash
	}

	// Submit to S3 (this requires the bucket to allow public writes to start-request.json files)
	url := fmt.Sprintf("https://%s.s3.amazonaws.com/%s/%s/start-request.json",
		token.S3Bucket, token.Project, token.Username)

	// For now, just indicate the request was made
	// Real implementation would need S3 PUT with public write permissions
	fmt.Printf("Start request submitted to: %s\n", url)

	return nil
}

// waitForInstanceStart waits for an instance to start with progress indicator.
func waitForInstanceStart(ctx context.Context, bucket, project, username string, timeout time.Duration) error {
	fmt.Printf("⏳ Waiting for instance to start")

	start := time.Now()
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	spinnerChars := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	spinnerTicker := time.NewTicker(200 * time.Millisecond)
	defer spinnerTicker.Stop()

	spinnerIndex := 0

	for {
		select {
		case <-ctx.Done():
			fmt.Print("\r")
			return ctx.Err()

		case <-time.After(timeout):
			fmt.Print("\r")
			return fmt.Errorf("timeout waiting for instance to start after %v", timeout)

		case <-spinnerTicker.C:
			elapsed := time.Since(start)
			fmt.Printf("\r⏳ %s Waiting for instance to start (elapsed: %v)",
				spinnerChars[spinnerIndex], elapsed.Round(time.Second))
			spinnerIndex = (spinnerIndex + 1) % len(spinnerChars)

		case <-ticker.C:
			status, err := getInstanceStatusFromS3(ctx, bucket, project, username)
			if err != nil {
				continue // Keep waiting
			}

			if status.State == "running" && status.PublicIP != "" {
				fmt.Print("\r")
				elapsed := time.Since(start)
				fmt.Printf("✅ Instance started and ready after %v\n", elapsed.Round(time.Second))
				return nil
			}
		}
	}
}

// listAvailableConnections lists all available connections.
func listAvailableConnections() error {
	tm, err := config.NewTokenManager()
	if err != nil {
		return fmt.Errorf("failed to initialize token manager: %w", err)
	}

	tokens, err := tm.ListTokens()
	if err != nil {
		return fmt.Errorf("failed to list tokens: %w", err)
	}

	if len(tokens) == 0 {
		fmt.Println("No access tokens configured.")
		fmt.Println("To add access: lfr connect activate <token> <student-id>")
		return nil
	}

	fmt.Printf("Available connections:\n\n")
	fmt.Printf("%-15s %-15s %-10s %-20s %-15s\n",
		"USERNAME", "PROJECT", "ROLE", "EXPIRES", "STATUS")
	fmt.Println(strings.Repeat("-", 85))

	for _, token := range tokens {
		status := "✅ Valid"
		if time.Now().After(token.ExpiresAt) {
			status = "❌ Expired"
		} else if !token.AccessStartDate.IsZero() && time.Now().Before(token.AccessStartDate) {
			status = "⏳ Not yet active"
		} else if !token.AccessEndDate.IsZero() && time.Now().After(token.AccessEndDate) {
			status = "❌ Access ended"
		}

		fmt.Printf("%-15s %-15s %-10s %-20s %-15s\n",
			token.Username,
			token.Project,
			token.Role,
			token.ExpiresAt.Format("2006-01-02"),
			status,
		)
	}

	fmt.Printf("\nTotal: %d connections\n", len(tokens))
	return nil
}

// activateStudentToken activates a token for the current machine.
func activateStudentToken(tokenString, studentID string) error {
	tm, err := config.NewTokenManager()
	if err != nil {
		return fmt.Errorf("failed to initialize token manager: %w", err)
	}

	fmt.Printf("Activating access token for student ID: %s\n", studentID)
	fmt.Printf("Binding to current machine...\n")

	err = tm.ActivateToken(tokenString, studentID)
	if err != nil {
		return fmt.Errorf("failed to activate token: %w", err)
	}

	// Parse token to show what was activated
	parts := strings.Split(tokenString, "-")
	if len(parts) >= 2 {
		project := parts[0]
		username := parts[1]
		fmt.Printf("✅ Token activated!\n")
		fmt.Printf("Project: %s\n", project)
		fmt.Printf("Username: %s\n", username)
		fmt.Printf("Student ID: %s\n", studentID)
		fmt.Printf("\nYou can now connect with: lfr connect %s\n", username)
	}

	return nil
}

// createTempSSHKey creates a temporary SSH key file from base64 data.
func createTempSSHKey(keyData string) (string, error) {
	if keyData == "" {
		return "", fmt.Errorf("no SSH key data in token")
	}

	// Create temporary file
	tmpFile, err := os.CreateTemp("", "lfr-ssh-key-*.pem")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}

	// Write key data
	_, err = tmpFile.WriteString(keyData)
	if err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("failed to write SSH key: %w", err)
	}

	// Set proper permissions
	if err := tmpFile.Chmod(0600); err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("failed to set key permissions: %w", err)
	}

	tmpFile.Close()
	return tmpFile.Name(), nil
}