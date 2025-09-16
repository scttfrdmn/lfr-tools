package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/scttfrdmn/lfr-tools/internal/aws"
	"github.com/scttfrdmn/lfr-tools/internal/config"
)

var efsCmd = &cobra.Command{
	Use:   "efs",
	Short: "Manage EFS shared storage for Lightsail instances",
	Long:  `Create, mount, and manage EFS file systems accessible from Lightsail instances via VPC peering.`,
}

var efsSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Set up VPC peering for EFS integration",
	Long: `Enable VPC peering between Lightsail and your default VPC to allow EFS access.
This is a one-time setup required before using EFS with Lightsail instances.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return setupEFSIntegration(cmd.Context())
	},
}

var efsCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create an EFS file system for Lightsail access",
	Long: `Create an EFS file system with proper configuration for Lightsail instances.
Automatically sets up mount targets and security groups for seamless access.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		project, _ := cmd.Flags().GetString("project")

		return createEFSFileSystem(cmd.Context(), name, project)
	},
}

var efsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List EFS file systems",
	Long:  `List all EFS file systems with their mount targets and access information.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		project, _ := cmd.Flags().GetString("project")

		return listEFSFileSystems(cmd.Context(), project)
	},
}

var efsMountCmd = &cobra.Command{
	Use:   "mount [filesystem-id] [username]",
	Short: "Mount EFS on a user's Lightsail instance",
	Long: `Mount an EFS file system on a user's Lightsail instance. This command connects
to the instance via SSH and configures the EFS mount with proper permissions.`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		filesystemID := args[0]
		username := args[1]
		mountPoint, _ := cmd.Flags().GetString("mount-point")
		project, _ := cmd.Flags().GetString("project")

		return mountEFSOnInstance(cmd.Context(), filesystemID, username, mountPoint, project)
	},
}

var efsStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check EFS and VPC peering status",
	Long:  `Check the status of VPC peering and EFS file systems, including mount target health.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return checkEFSStatus(cmd.Context())
	},
}

func init() {
	rootCmd.AddCommand(efsCmd)

	efsCmd.AddCommand(efsSetupCmd)
	efsCmd.AddCommand(efsCreateCmd)
	efsCmd.AddCommand(efsListCmd)
	efsCmd.AddCommand(efsMountCmd)
	efsCmd.AddCommand(efsStatusCmd)

	// Create command flags
	efsCreateCmd.Flags().StringP("project", "p", "", "Project name for tagging")

	// List command flags
	efsListCmd.Flags().StringP("project", "p", "", "Filter by project name")

	// Mount command flags
	efsMountCmd.Flags().StringP("mount-point", "m", "/mnt/efs", "Mount point on the instance")
	efsMountCmd.Flags().StringP("project", "p", "", "Filter by project name")
}

// setupEFSIntegration enables VPC peering for EFS integration.
func setupEFSIntegration(ctx context.Context) error {
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

	efsService := aws.NewEFSService(awsClient)

	fmt.Println("Setting up EFS integration with Lightsail...")

	// Check if VPC peering is already enabled
	isPeered, err := efsService.IsVPCPeered(ctx)
	if err != nil {
		return fmt.Errorf("failed to check VPC peering status: %w", err)
	}

	if isPeered {
		fmt.Println("✅ VPC peering is already enabled")
	} else {
		fmt.Println("Enabling VPC peering...")
		err = efsService.EnableVPCPeering(ctx)
		if err != nil {
			return fmt.Errorf("failed to enable VPC peering: %w", err)
		}
		fmt.Println("✅ VPC peering enabled")
	}

	// Check default VPC
	vpcID, err := efsService.GetDefaultVPC(ctx)
	if err != nil {
		return fmt.Errorf("failed to get default VPC: %w", err)
	}

	fmt.Printf("✅ Default VPC found: %s\n", vpcID)
	fmt.Println("🎉 EFS integration setup completed!")
	fmt.Println("\nYou can now create EFS file systems with: lfr efs create <name>")

	return nil
}

// createEFSFileSystem creates an EFS file system.
func createEFSFileSystem(ctx context.Context, name, project string) error {
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

	efsService := aws.NewEFSService(awsClient)

	// Check VPC peering
	isPeered, err := efsService.IsVPCPeered(ctx)
	if err != nil {
		return fmt.Errorf("failed to check VPC peering status: %w", err)
	}

	if !isPeered {
		return fmt.Errorf("VPC peering is not enabled. Run 'lfr efs setup' first")
	}

	fmt.Printf("Creating EFS file system: %s\n", name)
	if project != "" {
		fmt.Printf("Project: %s\n", project)
	}

	// Create file system
	fileSystem, err := efsService.CreateEFSFileSystem(ctx, name, project)
	if err != nil {
		return fmt.Errorf("failed to create EFS file system: %w", err)
	}

	fmt.Printf("✅ EFS file system created: %s\n", fileSystem.ID)

	// Create mount targets
	fmt.Println("Creating mount targets...")
	mountTargets, err := efsService.CreateMountTargets(ctx, fileSystem.ID)
	if err != nil {
		return fmt.Errorf("failed to create mount targets: %w", err)
	}

	fmt.Printf("✅ Created %d mount targets\n", len(mountTargets))

	for _, mt := range mountTargets {
		fmt.Printf("  Mount target: %s (IP: %s, AZ: %s)\n", mt.ID, mt.IPAddress, mt.AvailabilityZone)
	}

	fmt.Println("🎉 EFS file system setup completed!")
	fmt.Printf("\nMount on instances with: lfr efs mount %s <username>\n", fileSystem.ID)

	return nil
}

// listEFSFileSystems lists EFS file systems.
func listEFSFileSystems(ctx context.Context, project string) error {
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

	efsService := aws.NewEFSService(awsClient)

	// List file systems
	fileSystems, err := efsService.ListEFSFileSystems(ctx, project)
	if err != nil {
		return fmt.Errorf("failed to list EFS file systems: %w", err)
	}

	if len(fileSystems) == 0 {
		if project != "" {
			fmt.Printf("No EFS file systems found for project: %s\n", project)
		} else {
			fmt.Println("No EFS file systems found.")
		}
		return nil
	}

	// Display results
	if project != "" {
		fmt.Printf("EFS file systems for project: %s\n\n", project)
	} else {
		fmt.Printf("All EFS file systems:\n\n")
	}

	fmt.Printf("%-20s %-15s %-20s %-15s %-15s %-15s\n",
		"NAME", "ID", "STATE", "PERFORMANCE", "THROUGHPUT", "MOUNT TARGETS")
	fmt.Println(strings.Repeat("-", 110))

	for _, fs := range fileSystems {
		fmt.Printf("%-20s %-15s %-20s %-15s %-15s %-15d\n",
			fs.Name,
			fs.ID,
			fs.State,
			fs.PerformanceMode,
			fs.ThroughputMode,
			len(fs.MountTargets),
		)
	}

	fmt.Printf("\nTotal: %d EFS file systems\n", len(fileSystems))
	return nil
}

// mountEFSOnInstance mounts EFS on a user's instance.
func mountEFSOnInstance(ctx context.Context, filesystemID, username, mountPoint, project string) error {
	fmt.Printf("Mounting EFS %s on %s's instance at %s\n", filesystemID, username, mountPoint)

	// This would require SSH execution to the instance
	// For now, return instructions for manual mounting
	fmt.Printf(`
To mount EFS on the instance, SSH to the instance and run:

1. Install EFS utilities:
   sudo apt-get update
   sudo apt-get install -y nfs-common

2. Create mount point:
   sudo mkdir -p %s

3. Get EFS mount target IP:
   aws efs describe-mount-targets --file-system-id %s --query "MountTargets[0].IpAddress" --output text

4. Mount EFS (replace <IP> with mount target IP):
   sudo mount -t nfs4 -o nfsvers=4.1,rsize=1048576,wsize=1048576,hard,timeo=600,retrans=2 <IP>:/ %s

5. Add to /etc/fstab for persistent mounting:
   echo "<IP>:/ %s nfs4 nfsvers=4.1,rsize=1048576,wsize=1048576,hard,timeo=600,retrans=2,_netdev 0 0" | sudo tee -a /etc/fstab

`, mountPoint, filesystemID, mountPoint, mountPoint)

	return nil
}

// checkEFSStatus checks EFS and VPC peering status.
func checkEFSStatus(ctx context.Context) error {
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

	efsService := aws.NewEFSService(awsClient)

	fmt.Printf("EFS Integration Status for region: %s\n\n", awsClient.GetRegion())

	// Check VPC peering status
	isPeered, err := efsService.IsVPCPeered(ctx)
	if err != nil {
		fmt.Printf("❌ Failed to check VPC peering: %v\n", err)
	} else {
		if isPeered {
			fmt.Println("✅ VPC peering: Enabled")
		} else {
			fmt.Println("❌ VPC peering: Disabled (run 'lfr efs setup')")
		}
	}

	// Check default VPC
	vpcID, err := efsService.GetDefaultVPC(ctx)
	if err != nil {
		fmt.Printf("❌ Default VPC: Not found (%v)\n", err)
	} else {
		fmt.Printf("✅ Default VPC: %s\n", vpcID)
	}

	// List EFS file systems
	fileSystems, err := efsService.ListEFSFileSystems(ctx, "")
	if err != nil {
		fmt.Printf("❌ Failed to list EFS file systems: %v\n", err)
	} else {
		fmt.Printf("📁 EFS file systems: %d found\n", len(fileSystems))

		for _, fs := range fileSystems {
			fmt.Printf("   %s (%s) - %s - %d mount targets\n",
				fs.Name, fs.ID, fs.State, len(fs.MountTargets))
		}
	}

	return nil
}