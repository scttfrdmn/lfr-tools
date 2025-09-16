package cmd

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/scttfrdmn/lfr-tools/internal/aws"
	"github.com/scttfrdmn/lfr-tools/internal/config"
	"github.com/scttfrdmn/lfr-tools/internal/types"
)

var sshCmd = &cobra.Command{
	Use:   "ssh",
	Short: "Manage SSH access to Lightsail instances",
	Long:  `Simplify SSH access to Lightsail instances with automatic key management and connection helpers.`,
}

var sshConnectCmd = &cobra.Command{
	Use:   "connect [username]",
	Short: "Connect to a user's Lightsail instance via SSH",
	Long: `Connect to a user's Lightsail instance using SSH. Automatically handles key download,
connection setup, and provides a seamless SSH experience.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		username := args[0]
		project, _ := cmd.Flags().GetString("project")
		keyPath, _ := cmd.Flags().GetString("key-path")

		return connectSSH(cmd.Context(), username, project, keyPath)
	},
}

var sshKeysCmd = &cobra.Command{
	Use:   "keys",
	Short: "Manage SSH keys for Lightsail instances",
	Long:  `Download, list, and manage SSH keys for Lightsail instances.`,
}

var sshKeysDownloadCmd = &cobra.Command{
	Use:   "download [username]",
	Short: "Download SSH key for a user's instance",
	Long: `Download the SSH private key for a user's Lightsail instance and save it to
the local SSH directory with proper permissions.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		username := args[0]
		project, _ := cmd.Flags().GetString("project")
		outputPath, _ := cmd.Flags().GetString("output")

		return downloadSSHKey(cmd.Context(), username, project, outputPath)
	},
}

var sshKeysListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available SSH keys",
	Long:  `List all available SSH keys for Lightsail instances, showing which instances they can access.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		project, _ := cmd.Flags().GetString("project")

		if project != "" {
			fmt.Printf("Listing SSH keys for project: %s\n", project)
		} else {
			fmt.Println("Listing all SSH keys")
		}

		// TODO: Implement SSH key listing logic
		return fmt.Errorf("SSH key listing not yet implemented")
	},
}

var sshConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Generate SSH config entries",
	Long: `Generate SSH config entries for easy access to Lightsail instances. This creates
proper SSH config entries with hostnames, users, and key paths.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		project, _ := cmd.Flags().GetString("project")
		outputPath, _ := cmd.Flags().GetString("output")

		fmt.Printf("Generating SSH config")
		if project != "" {
			fmt.Printf(" for project: %s", project)
		}
		if outputPath != "" {
			fmt.Printf(" to file: %s", outputPath)
		}
		fmt.Println()

		// TODO: Implement SSH config generation logic
		return fmt.Errorf("SSH config generation not yet implemented")
	},
}

var sshTunnelCmd = &cobra.Command{
	Use:   "tunnel [username] [local_port:remote_port]",
	Short: "Create SSH tunnel to instance",
	Long: `Create an SSH tunnel to a user's Lightsail instance for secure access to services
running on the instance (e.g., Jupyter notebooks, web servers).`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		username := args[0]
		portMapping := args[1]
		project, _ := cmd.Flags().GetString("project")

		fmt.Printf("Creating SSH tunnel to %s instance\n", username)
		fmt.Printf("Port mapping: %s\n", portMapping)
		if project != "" {
			fmt.Printf("Project: %s\n", project)
		}

		// TODO: Implement SSH tunnel logic
		return fmt.Errorf("SSH tunnel not yet implemented")
	},
}

func init() {
	rootCmd.AddCommand(sshCmd)

	sshCmd.AddCommand(sshConnectCmd)
	sshCmd.AddCommand(sshKeysCmd)
	sshCmd.AddCommand(sshConfigCmd)
	sshCmd.AddCommand(sshTunnelCmd)

	sshKeysCmd.AddCommand(sshKeysDownloadCmd)
	sshKeysCmd.AddCommand(sshKeysListCmd)

	// Connect command flags
	sshConnectCmd.Flags().StringP("project", "p", "", "Filter by project name")
	sshConnectCmd.Flags().StringP("key-path", "k", "", "Path to SSH private key")

	// Key download command flags
	sshKeysDownloadCmd.Flags().StringP("project", "p", "", "Filter by project name")
	sshKeysDownloadCmd.Flags().StringP("output", "o", "", "Output path for the key file")

	// Key list command flags
	sshKeysListCmd.Flags().StringP("project", "p", "", "Filter by project name")

	// Config command flags
	sshConfigCmd.Flags().StringP("project", "p", "", "Filter by project name")
	sshConfigCmd.Flags().StringP("output", "o", "", "Output path for SSH config (default: stdout)")

	// Tunnel command flags
	sshTunnelCmd.Flags().StringP("project", "p", "", "Filter by project name")
}

// connectSSH connects to a user's instance via SSH.
func connectSSH(ctx context.Context, username, project, keyPath string) error {
	// Load configuration
	cfg, err := config.Load()
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

	// Find the user's instance
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

	if targetInstance.PublicIP == "" {
		return fmt.Errorf("instance %s has no public IP address", targetInstance.Name)
	}

	fmt.Printf("Connecting to %s's instance: %s\n", username, targetInstance.Name)
	fmt.Printf("Instance: %s (%s)\n", targetInstance.PublicIP, targetInstance.State)

	// Use custom key path if provided, otherwise try to download/use default
	var privateKeyPath string
	if keyPath != "" {
		privateKeyPath = keyPath
	} else {
		// Try to download and use the default key
		privateKeyPath = filepath.Join(cfg.SSH.KeyPath, "LightsailDefaultKey.pem")

		// Ensure the SSH key directory exists
		if err := os.MkdirAll(cfg.SSH.KeyPath, 0700); err != nil {
			return fmt.Errorf("failed to create SSH key directory: %w", err)
		}

		// Download the key if it doesn't exist
		if _, err := os.Stat(privateKeyPath); os.IsNotExist(err) {
			fmt.Println("Downloading SSH key...")
			keyContent, err := lightsailService.DownloadSSHKey(ctx, "")
			if err != nil {
				return fmt.Errorf("failed to download SSH key: %w", err)
			}

			// Check if key content is already decoded or needs base64 decoding
			var keyBytes []byte
			if strings.HasPrefix(keyContent, "-----BEGIN") {
				// Key is already in PEM format
				keyBytes = []byte(keyContent)
			} else {
				// Try to decode from base64
				decoded, err := base64.StdEncoding.DecodeString(keyContent)
				if err != nil {
					return fmt.Errorf("failed to decode SSH key (not base64): %w", err)
				}
				keyBytes = decoded
			}

			// Write key file with proper permissions
			err = os.WriteFile(privateKeyPath, keyBytes, 0600)
			if err != nil {
				return fmt.Errorf("failed to write SSH key file: %w", err)
			}

			fmt.Printf("SSH key saved to: %s\n", privateKeyPath)
		}
	}

	// Verify key file exists and has proper permissions
	if _, err := os.Stat(privateKeyPath); os.IsNotExist(err) {
		return fmt.Errorf("SSH key file not found: %s", privateKeyPath)
	}

	// Ensure proper permissions
	if err := os.Chmod(privateKeyPath, 0600); err != nil {
		return fmt.Errorf("failed to set SSH key permissions: %w", err)
	}

	// Execute SSH command
	sshArgs := []string{
		"-i", privateKeyPath,
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
		fmt.Sprintf("ubuntu@%s", targetInstance.PublicIP),
	}

	fmt.Printf("Executing: ssh %s\n\n", strings.Join(sshArgs, " "))

	sshCmdExec := exec.Command("ssh", sshArgs...)
	sshCmdExec.Stdin = os.Stdin
	sshCmdExec.Stdout = os.Stdout
	sshCmdExec.Stderr = os.Stderr

	return sshCmdExec.Run()
}

// downloadSSHKey downloads SSH key for a user's instance.
func downloadSSHKey(ctx context.Context, username, project, outputPath string) error {
	// Load configuration
	cfg, err := config.Load()
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

	// Download the key
	fmt.Printf("Downloading SSH key for user: %s\n", username)
	keyContent, err := lightsailService.DownloadSSHKey(ctx, "")
	if err != nil {
		return fmt.Errorf("failed to download SSH key: %w", err)
	}

	// Check if key content is already decoded or needs base64 decoding
	var keyBytes []byte
	if strings.HasPrefix(keyContent, "-----BEGIN") {
		// Key is already in PEM format
		keyBytes = []byte(keyContent)
	} else {
		// Try to decode from base64
		decoded, err := base64.StdEncoding.DecodeString(keyContent)
		if err != nil {
			return fmt.Errorf("failed to decode SSH key (not base64): %w", err)
		}
		keyBytes = decoded
	}

	// Determine output path
	var keyPath string
	if outputPath != "" {
		keyPath = filepath.Join(outputPath, "LightsailDefaultKey.pem")
		// Ensure output directory exists
		if err := os.MkdirAll(outputPath, 0700); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}
	} else {
		keyPath = filepath.Join(cfg.SSH.KeyPath, "LightsailDefaultKey.pem")
		// Ensure SSH key directory exists
		if err := os.MkdirAll(cfg.SSH.KeyPath, 0700); err != nil {
			return fmt.Errorf("failed to create SSH key directory: %w", err)
		}
	}

	// Write key file with proper permissions
	err = os.WriteFile(keyPath, keyBytes, 0600)
	if err != nil {
		return fmt.Errorf("failed to write SSH key file: %w", err)
	}

	fmt.Printf("✅ SSH key downloaded to: %s\n", keyPath)
	fmt.Printf("Key permissions set to 600 (owner read/write only)\n")

	return nil
}