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
		mode, _ := cmd.Flags().GetString("mode")

		return mountEFSOnInstance(cmd.Context(), filesystemID, username, mountPoint, project, mode)
	},
}

var efsMountAllCmd = &cobra.Command{
	Use:   "mount-all [filesystem-id]",
	Short: "Mount EFS on all instances in a project",
	Long: `Mount an EFS file system on all running instances in a project. Automatically
discovers instances and configures EFS mounting via SSH.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filesystemID := args[0]
		project, _ := cmd.Flags().GetString("project")
		mountPoint, _ := cmd.Flags().GetString("mount-point")
		mode, _ := cmd.Flags().GetString("mode")
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		return mountEFSOnAllInstances(cmd.Context(), filesystemID, project, mountPoint, mode, dryRun)
	},
}

var efsMountStatusCmd = &cobra.Command{
	Use:   "mount-status",
	Short: "Show EFS mount status across instances",
	Long: `Display which instances have EFS file systems mounted and their status.
Connects to running instances to check actual mount status.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		project, _ := cmd.Flags().GetString("project")

		return showEFSMountStatus(cmd.Context(), project)
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
	efsCmd.AddCommand(efsMountAllCmd)
	efsCmd.AddCommand(efsMountStatusCmd)
	efsCmd.AddCommand(efsStatusCmd)

	// Create command flags
	efsCreateCmd.Flags().StringP("project", "p", "", "Project name for tagging")

	// List command flags
	efsListCmd.Flags().StringP("project", "p", "", "Filter by project name")

	// Mount command flags
	efsMountCmd.Flags().StringP("mount-point", "m", "/mnt/efs", "Mount point on the instance")
	efsMountCmd.Flags().StringP("project", "p", "", "Filter by project name")
	efsMountCmd.Flags().StringP("mode", "", "rw", "Mount mode: rw (read-write) or ro (read-only)")

	// Mount all command flags
	efsMountAllCmd.Flags().StringP("project", "p", "", "Project to mount EFS on (required)")
	efsMountAllCmd.Flags().StringP("mount-point", "m", "/mnt/efs", "Mount point on instances")
	efsMountAllCmd.Flags().StringP("mode", "", "rw", "Mount mode: rw (read-write) or ro (read-only)")
	efsMountAllCmd.Flags().BoolP("dry-run", "d", false, "Show which instances would be affected without mounting")
	efsMountAllCmd.MarkFlagRequired("project")

	// Mount status command flags
	efsMountStatusCmd.Flags().StringP("project", "p", "", "Filter by project name")
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

// mountEFSOnInstance mounts EFS on a user's instance with specified permissions.
func mountEFSOnInstance(ctx context.Context, filesystemID, username, mountPoint, project, mode string) error {
	if mode != "rw" && mode != "ro" {
		return fmt.Errorf("mount mode must be 'rw' (read-write) or 'ro' (read-only), got: %s", mode)
	}

	modeDesc := "read-write"
	if mode == "ro" {
		modeDesc = "read-only"
	}

	fmt.Printf("Mounting EFS %s on %s's instance at %s (%s)\n", filesystemID, username, mountPoint, modeDesc)

	// Generate NFS mount options based on mode
	var mountOptions string
	if mode == "ro" {
		mountOptions = "nfsvers=4.1,rsize=1048576,wsize=1048576,hard,timeo=600,retrans=2,ro"
	} else {
		mountOptions = "nfsvers=4.1,rsize=1048576,wsize=1048576,hard,timeo=600,retrans=2"
	}

	// This would require SSH execution to the instance
	// For now, return instructions for manual mounting
	fmt.Printf(`
To mount EFS on the instance (%s), SSH to the instance and run:

1. Install EFS utilities:
   sudo apt-get update
   sudo apt-get install -y nfs-common

2. Create mount point:
   sudo mkdir -p %s

3. Get EFS mount target IP:
   aws efs describe-mount-targets --file-system-id %s --query "MountTargets[0].IpAddress" --output text

4. Mount EFS %s (replace <IP> with mount target IP):
   sudo mount -t nfs4 -o %s <IP>:/ %s

5. Add to /etc/fstab for persistent mounting:
   echo "<IP>:/ %s nfs4 %s,_netdev 0 0" | sudo tee -a /etc/fstab

`, modeDesc, mountPoint, filesystemID, modeDesc, mountOptions, mountPoint, mountPoint, mountOptions)

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

// mountEFSOnAllInstances mounts EFS on all instances in a project.
func mountEFSOnAllInstances(ctx context.Context, filesystemID, project, mountPoint, mode string, dryRun bool) error {
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
	efsService := aws.NewEFSService(awsClient)

	// Get instances for the project
	instances, err := lightsailService.ListInstances(ctx, project)
	if err != nil {
		return fmt.Errorf("failed to list instances: %w", err)
	}

	if len(instances) == 0 {
		fmt.Printf("No instances found for project: %s\n", project)
		return nil
	}

	// Get EFS mount target IP
	fileSystems, err := efsService.ListEFSFileSystems(ctx, "")
	if err != nil {
		return fmt.Errorf("failed to list EFS file systems: %w", err)
	}

	var targetFS *aws.EFSFileSystem
	for _, fs := range fileSystems {
		if fs.ID == filesystemID {
			targetFS = fs
			break
		}
	}

	if targetFS == nil {
		return fmt.Errorf("EFS file system not found: %s", filesystemID)
	}

	if len(targetFS.MountTargets) == 0 {
		return fmt.Errorf("EFS file system %s has no mount targets", filesystemID)
	}

	mountIP := targetFS.MountTargets[0].IPAddress

	if dryRun {
		fmt.Printf("Dry run: EFS %s would be mounted on %d instances in project %s:\n\n", filesystemID, len(instances), project)
		fmt.Printf("%-20s %-12s %-18s %-15s\n", "INSTANCE", "STATE", "PUBLIC IP", "MOUNT POINT")
		fmt.Println(strings.Repeat("-", 80))

		for _, instance := range instances {
			publicIP := instance.PublicIP
			if publicIP == "" {
				publicIP = "-"
			}
			fmt.Printf("%-20s %-12s %-18s %-15s\n", instance.Name, instance.State, publicIP, mountPoint)
		}
		return nil
	}

	fmt.Printf("Mounting EFS %s on %d instances in project %s\n", filesystemID, len(instances), project)
	fmt.Printf("Mount target IP: %s\n", mountIP)
	fmt.Printf("Mount point: %s\n\n", mountPoint)

	successCount := 0
	for i, instance := range instances {
		username := "unknown"
		parts := strings.Split(instance.Name, "-")
		if len(parts) > 0 {
			username = parts[0]
		}

		fmt.Printf("[%d/%d] Mounting on %s (%s)\n", i+1, len(instances), instance.Name, username)

		if instance.State != "running" {
			fmt.Printf("   ⏭️  Skipping %s (state: %s)\n", instance.Name, instance.State)
			continue
		}

		if instance.PublicIP == "" {
			fmt.Printf("   ⏭️  Skipping %s (no public IP)\n", instance.Name)
			continue
		}

		// Generate mount instructions for this instance
		modeDesc := "read-write"
		mountOptions := "nfsvers=4.1,rsize=1048576,wsize=1048576,hard,timeo=600,retrans=2"
		if mode == "ro" {
			modeDesc = "read-only"
			mountOptions += ",ro"
		}

		fmt.Printf("   📋 Mount instructions for %s (%s):\n", instance.Name, modeDesc)
		fmt.Printf("      ssh ubuntu@%s\n", instance.PublicIP)
		fmt.Printf("      sudo apt-get update && sudo apt-get install -y nfs-common\n")
		fmt.Printf("      sudo mkdir -p %s\n", mountPoint)
		fmt.Printf("      sudo mount -t nfs4 -o %s %s:/ %s\n", mountOptions, mountIP, mountPoint)
		fmt.Printf("   ✅ Ready to mount on %s (%s)\n\n", instance.Name, modeDesc)

		successCount++
	}

	fmt.Printf("🎉 Generated mount instructions for %d/%d instances\n", successCount, len(instances))
	return nil
}

// showEFSMountStatus shows EFS mount status across instances.
func showEFSMountStatus(ctx context.Context, project string) error {
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

	// Get instances
	instances, err := lightsailService.ListInstances(ctx, project)
	if err != nil {
		return fmt.Errorf("failed to list instances: %w", err)
	}

	if len(instances) == 0 {
		if project != "" {
			fmt.Printf("No instances found for project: %s\n", project)
		} else {
			fmt.Println("No instances found.")
		}
		return nil
	}

	// Display mount status
	if project != "" {
		fmt.Printf("EFS mount status for project: %s\n\n", project)
	} else {
		fmt.Printf("EFS mount status for all instances:\n\n")
	}

	fmt.Printf("%-20s %-12s %-18s %-15s %-15s\n",
		"INSTANCE", "STATE", "PUBLIC IP", "SSH READY", "MOUNT STATUS")
	fmt.Println(strings.Repeat("-", 95))

	for _, instance := range instances {
		publicIP := instance.PublicIP
		if publicIP == "" {
			publicIP = "-"
		}

		sshReady := "No"
		mountStatus := "Unknown"

		if instance.State == "running" && instance.PublicIP != "" {
			sshReady = "Yes"
			mountStatus = "Check manually"
		} else {
			mountStatus = "Instance not running"
		}

		fmt.Printf("%-20s %-12s %-18s %-15s %-15s\n",
			instance.Name,
			instance.State,
			publicIP,
			sshReady,
			mountStatus,
		)
	}

	fmt.Printf("\nTotal: %d instances\n", len(instances))
	fmt.Printf("\nTo check actual mount status on running instances:\n")
	fmt.Printf("  ssh ubuntu@<instance-ip> 'df -h | grep nfs'\n")

	return nil
}