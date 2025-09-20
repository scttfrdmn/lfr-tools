package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/scttfrdmn/lfr-tools/internal/aws"
	"github.com/scttfrdmn/lfr-tools/internal/config"
	"github.com/scttfrdmn/lfr-tools/internal/types"
)

var dcvCmd = &cobra.Command{
	Use:   "dcv",
	Short: "Manage NICE DCV connections to Lightsail instances",
	Long:  `Configure and connect to Lightsail instances using NICE DCV for remote desktop access with optimized settings.`,
}

var dcvConnectCmd = &cobra.Command{
	Use:   "connect [username]",
	Short: "Connect to a user's instance via NICE DCV",
	Long: `Launch a NICE DCV connection to a user's Lightsail instance with optimized settings
for better performance and reliability. Automatically handles authentication and connection setup.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		username := args[0]
		project, _ := cmd.Flags().GetString("project")
		quality, _ := cmd.Flags().GetString("quality")
		fullscreen, _ := cmd.Flags().GetBool("fullscreen")

		return connectDCV(cmd.Context(), username, project, quality, fullscreen)
	},
}

var dcvConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure NICE DCV settings for instances",
	Long: `Configure NICE DCV settings on Lightsail instances to optimize performance
and reliability. Sets up authentication, display settings, and performance parameters.`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		project, _ := cmd.Flags().GetString("project")
		users, _ := cmd.Flags().GetStringSlice("users")
		quality, _ := cmd.Flags().GetString("quality")
		maxSessions, _ := cmd.Flags().GetInt("max-sessions")

		if len(users) > 0 {
			fmt.Printf("Configuring NICE DCV for users: %v\n", users)
		} else {
			fmt.Println("Configuring NICE DCV for all instances")
		}

		if project != "" {
			fmt.Printf("Project filter: %s\n", project)
		}
		fmt.Printf("Default quality: %s\n", quality)
		fmt.Printf("Max sessions per instance: %d\n", maxSessions)

		// TODO: Implement NICE DCV configuration logic
		return fmt.Errorf("NICE DCV configuration not yet implemented")
	},
}

var dcvStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check NICE DCV status on instances",
	Long: `Check the status of NICE DCV services on Lightsail instances, including
service health, active sessions, and configuration status.`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		project, _ := cmd.Flags().GetString("project")
		user, _ := cmd.Flags().GetString("user")

		fmt.Println("Checking NICE DCV status")
		if project != "" {
			fmt.Printf("Project filter: %s\n", project)
		}
		if user != "" {
			fmt.Printf("User filter: %s\n", user)
		}

		// TODO: Implement NICE DCV status check logic
		return fmt.Errorf("NICE DCV status check not yet implemented")
	},
}

var dcvSessionsCmd = &cobra.Command{
	Use:   "sessions",
	Short: "Manage NICE DCV sessions",
	Long:  `List, terminate, and manage active NICE DCV sessions on instances.`,
}

var dcvSessionsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List active NICE DCV sessions",
	Long:  `List all active NICE DCV sessions across instances with session details.`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		project, _ := cmd.Flags().GetString("project")
		user, _ := cmd.Flags().GetString("user")

		fmt.Println("Listing active NICE DCV sessions")
		if project != "" {
			fmt.Printf("Project filter: %s\n", project)
		}
		if user != "" {
			fmt.Printf("User filter: %s\n", user)
		}

		// TODO: Implement DCV session listing logic
		return fmt.Errorf("DCV session listing not yet implemented")
	},
}

var dcvSessionsTerminateCmd = &cobra.Command{
	Use:   "terminate [session-id]",
	Short: "Terminate a NICE DCV session",
	Long:  `Terminate a specific NICE DCV session by ID. This will disconnect the user.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		sessionID := args[0]
		force, _ := cmd.Flags().GetBool("force")

		if force {
			fmt.Printf("Force terminating DCV session: %s\n", sessionID)
		} else {
			fmt.Printf("Terminating DCV session: %s\n", sessionID)
		}

		// TODO: Implement DCV session termination logic
		return fmt.Errorf("DCV session termination not yet implemented")
	},
}

var dcvOptimizeCmd = &cobra.Command{
	Use:   "optimize",
	Short: "Optimize NICE DCV settings for performance",
	Long: `Apply performance optimizations to NICE DCV configurations based on
instance type, network conditions, and usage patterns.`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		project, _ := cmd.Flags().GetString("project")
		users, _ := cmd.Flags().GetStringSlice("users")
		profile, _ := cmd.Flags().GetString("profile")

		fmt.Printf("Optimizing NICE DCV settings with profile: %s\n", profile)
		if len(users) > 0 {
			fmt.Printf("Target users: %v\n", users)
		}
		if project != "" {
			fmt.Printf("Project filter: %s\n", project)
		}

		// TODO: Implement DCV optimization logic
		return fmt.Errorf("DCV optimization not yet implemented")
	},
}

func init() {
	rootCmd.AddCommand(dcvCmd)

	dcvCmd.AddCommand(dcvConnectCmd)
	dcvCmd.AddCommand(dcvConfigCmd)
	dcvCmd.AddCommand(dcvStatusCmd)
	dcvCmd.AddCommand(dcvSessionsCmd)
	dcvCmd.AddCommand(dcvOptimizeCmd)

	dcvSessionsCmd.AddCommand(dcvSessionsListCmd)
	dcvSessionsCmd.AddCommand(dcvSessionsTerminateCmd)

	// Connect command flags
	dcvConnectCmd.Flags().StringP("project", "p", "", "Filter by project name")
	dcvConnectCmd.Flags().StringP("quality", "q", "medium", "Connection quality (low, medium, high, lossless)")
	dcvConnectCmd.Flags().BoolP("fullscreen", "f", false, "Launch in fullscreen mode")

	// Config command flags
	dcvConfigCmd.Flags().StringP("project", "p", "", "Filter by project name")
	dcvConfigCmd.Flags().StringSliceP("users", "u", []string{}, "Target specific users")
	dcvConfigCmd.Flags().StringP("quality", "q", "medium", "Default quality setting")
	dcvConfigCmd.Flags().IntP("max-sessions", "m", 1, "Maximum concurrent sessions per instance")

	// Status command flags
	dcvStatusCmd.Flags().StringP("project", "p", "", "Filter by project name")
	dcvStatusCmd.Flags().StringP("user", "u", "", "Filter by username")

	// Sessions list command flags
	dcvSessionsListCmd.Flags().StringP("project", "p", "", "Filter by project name")
	dcvSessionsListCmd.Flags().StringP("user", "u", "", "Filter by username")

	// Sessions terminate command flags
	dcvSessionsTerminateCmd.Flags().BoolP("force", "f", false, "Force terminate without confirmation")

	// Optimize command flags
	dcvOptimizeCmd.Flags().StringP("project", "p", "", "Filter by project name")
	dcvOptimizeCmd.Flags().StringSliceP("users", "u", []string{}, "Target specific users")
	dcvOptimizeCmd.Flags().StringP("profile", "P", "balanced", "Optimization profile (performance, balanced, bandwidth-saver)")
}

// connectDCV implements NICE DCV connection with real SSH execution
func connectDCV(ctx context.Context, username, project, quality string, fullscreen bool) error {
	fmt.Printf("Connecting to %s's instance via NICE DCV\n", username)
	if project != "" {
		fmt.Printf("Project: %s\n", project)
	}
	fmt.Printf("Quality: %s\n", quality)

	// Load configuration
	_, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Create AWS client
	awsClient, err := aws.NewClient(ctx, aws.Options{
		Region:  viper.GetString("aws.region"),
		Profile: viper.GetString("aws.profile"),
	})
	if err != nil {
		return fmt.Errorf("failed to create AWS client: %w", err)
	}

	lightsailService := aws.NewLightsailService(awsClient)

	// Find user's instance
	instances, err := lightsailService.ListInstances(ctx, project)
	if err != nil {
		return fmt.Errorf("failed to list instances: %w", err)
	}

	var targetInstance *types.Instance
	for _, instance := range instances {
		if strings.HasPrefix(instance.Name, username+"-") {
			targetInstance = instance
			break
		}
	}

	if targetInstance == nil {
		return fmt.Errorf("no instance found for user: %s", username)
	}

	if targetInstance.State != instanceStateRunning {
		return fmt.Errorf("instance %s is not running (state: %s). Start it first", targetInstance.Name, targetInstance.State)
	}

	if targetInstance.PublicIP == "" {
		return fmt.Errorf("instance %s has no public IP", targetInstance.Name)
	}

	fmt.Printf("Target instance: %s (%s)\n", targetInstance.Name, targetInstance.PublicIP)

	// Configure and start DCV session via SSH
	fmt.Printf("Configuring DCV session...\n")

	// SSH to instance and setup DCV
	sshCommand := fmt.Sprintf(`
		# Ensure DCV server is installed and running
		if ! command -v dcv &> /dev/null; then
			echo "Installing DCV server..."
			wget -q https://d1uj6qtbmh3dt5.cloudfront.net/nice-dcv-ubuntu2204-x86_64.tgz
			tar -xzf nice-dcv-ubuntu2204-x86_64.tgz
			cd nice-dcv-*-ubuntu2204-x86_64
			sudo apt install -y ./nice-dcv-server_*.deb
			sudo apt install -y ./nice-xdcv_*.deb
		fi

		# Configure DCV server
		sudo systemctl enable dcvserver
		sudo systemctl start dcvserver

		# Create DCV session for user
		dcv create-session %s-session --user ubuntu --type=virtual

		# Configure session quality
		dcv set-session-parameter %s-session frame-rate %s

		echo "DCV session ready on port 8443"
	`, username, username, getDCVFrameRate(quality))

	// Execute via SSH
	err = executeSSHCommand(ctx, targetInstance.PublicIP, sshCommand)
	if err != nil {
		return fmt.Errorf("failed to configure DCV: %w", err)
	}

	// Launch DCV viewer
	dcvURL := fmt.Sprintf("dcv://%s:8443/%s-session", targetInstance.PublicIP, username)

	fmt.Printf("✅ DCV session created successfully!\n")
	fmt.Printf("Session ID: %s-session\n", username)
	fmt.Printf("Connection URL: %s\n", dcvURL)

	// Launch platform-specific DCV viewer
	return launchDCVViewer(dcvURL, quality, fullscreen)
}

// executeSSHCommand executes a command on an instance via SSH
func executeSSHCommand(_ context.Context, publicIP, command string) error {
	// Get SSH key path
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	keyPath := filepath.Join(homeDir, ".ssh", "lfr-tools", "LightsailDefaultKey.pem")

	// Execute SSH command
	// #nosec G204 - SSH command execution is intentional for DCV setup
	cmd := exec.Command("ssh",
		"-i", keyPath,
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
		"-o", "ConnectTimeout=10",
		fmt.Sprintf("ubuntu@%s", publicIP),
		command)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("SSH command failed: %w (output: %s)", err, string(output))
	}

	fmt.Printf("SSH execution completed: %s\n", string(output))
	return nil
}

// getDCVFrameRate returns appropriate frame rate for quality setting
func getDCVFrameRate(quality string) string {
	switch quality {
	case "low":
		return "15"
	case "medium":
		return "30"
	case "high":
		return "60"
	case "lossless":
		return "60"
	default:
		return "30"
	}
}

const (
	platformDarwin  = "darwin"
	platformWindows = "windows"
	platformLinux   = "linux"
)

// launchDCVViewer launches the appropriate DCV viewer
func launchDCVViewer(dcvURL, _ string, fullscreen bool) error {
	// Try native DCV viewer first
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case platformDarwin:
		dcvPath := "/Applications/DCV Viewer.app/Contents/MacOS/DCV Viewer"
		if _, err := os.Stat(dcvPath); err == nil {
			args := []string{dcvURL}
			if fullscreen {
				args = append(args, "--full-screen")
			}
			// #nosec G204 - DCV viewer command execution is intentional
			cmd = exec.Command(dcvPath, args...)
		}
	case platformWindows:
		dcvPath := "C:\\Program Files\\NICE\\DCV\\Client\\bin\\dcvviewer.exe"
		if _, err := os.Stat(dcvPath); err == nil {
			args := []string{dcvURL}
			if fullscreen {
				args = append(args, "--full-screen")
			}
			// #nosec G204 - DCV viewer command execution is intentional
			cmd = exec.Command(dcvPath, args...)
		}
	case platformLinux:
		if _, err := exec.LookPath("dcvviewer"); err == nil {
			args := []string{dcvURL}
			if fullscreen {
				args = append(args, "--full-screen")
			}
			// #nosec G204 - DCV viewer command execution is intentional
			cmd = exec.Command("dcvviewer", args...)
		}
	}

	// Launch native viewer if available
	if cmd != nil {
		fmt.Printf("Launching native DCV viewer...\n")
		err := cmd.Start()
		if err != nil {
			fmt.Printf("Failed to launch native viewer: %v\n", err)
		} else {
			fmt.Printf("✅ DCV viewer launched successfully!\n")
			return nil
		}
	}

	// Fall back to web browser
	fmt.Printf("Native DCV viewer not found, opening in web browser...\n")

	// Convert dcv:// URL to https:// for web access
	webURL := strings.Replace(dcvURL, "dcv://", "https://", 1)

	var browserCmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		// #nosec G204 - Browser command execution is intentional for DCV web access
		browserCmd = exec.Command("open", webURL)
	case platformWindows:
		// #nosec G204 - Browser command execution is intentional for DCV web access
		browserCmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", webURL)
	case platformLinux:
		// #nosec G204 - Browser command execution is intentional for DCV web access
		browserCmd = exec.Command("xdg-open", webURL)
	default:
		return fmt.Errorf("unsupported platform for DCV viewer launch")
	}

	err := browserCmd.Start()
	if err != nil {
		return fmt.Errorf("failed to open DCV in browser: %w", err)
	}

	fmt.Printf("✅ DCV web viewer opened in browser!\n")
	fmt.Printf("If the page doesn't load, the DCV server may still be starting.\n")
	fmt.Printf("Wait 1-2 minutes and refresh the page.\n")

	return nil
}