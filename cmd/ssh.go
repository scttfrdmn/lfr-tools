package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
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

		fmt.Printf("Connecting to instance for user: %s\n", username)
		if project != "" {
			fmt.Printf("Project: %s\n", project)
		}
		if keyPath != "" {
			fmt.Printf("Using key: %s\n", keyPath)
		}

		// TODO: Implement SSH connection logic
		return fmt.Errorf("SSH connect not yet implemented")
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

		fmt.Printf("Downloading SSH key for user: %s\n", username)
		if project != "" {
			fmt.Printf("Project: %s\n", project)
		}
		if outputPath != "" {
			fmt.Printf("Output path: %s\n", outputPath)
		}

		// TODO: Implement SSH key download logic
		return fmt.Errorf("SSH key download not yet implemented")
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