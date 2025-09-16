package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/scttfrdmn/lfr-tools/internal/aws"
	"github.com/scttfrdmn/lfr-tools/internal/config"
	"github.com/scttfrdmn/lfr-tools/internal/utils"
)

var volumesCmd = &cobra.Command{
	Use:   "volumes",
	Short: "Manage Lightsail block storage volumes",
	Long:  `Create, attach, detach, and manage Lightsail block storage volumes for persistent data.`,
}

var volumesCreateCmd = &cobra.Command{
	Use:   "create [name] [size-gb]",
	Short: "Create a new block storage volume",
	Long: `Create a new SSD block storage volume that can be attached to Lightsail instances.
Volume size must be between 8 GB and 16 TB (16384 GB).`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		sizeGB := args[1]
		project, _ := cmd.Flags().GetString("project")
		zone, _ := cmd.Flags().GetString("zone")
		wait, _ := cmd.Flags().GetBool("wait")

		return createVolume(cmd.Context(), name, sizeGB, zone, project, wait)
	},
}

var volumesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List block storage volumes",
	Long:  `List all block storage volumes with their attachment status and details.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		project, _ := cmd.Flags().GetString("project")

		return listVolumes(cmd.Context(), project)
	},
}

var volumesAttachCmd = &cobra.Command{
	Use:   "attach [volume-name] [instance-name] [device-path]",
	Short: "Attach a volume to an instance",
	Long: `Attach a block storage volume to a Lightsail instance at the specified device path.
Common device paths: /dev/xvdf, /dev/xvdg, /dev/xvdh, etc.`,
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		volumeName := args[0]
		instanceName := args[1]
		devicePath := args[2]
		wait, _ := cmd.Flags().GetBool("wait")

		return attachVolume(cmd.Context(), volumeName, instanceName, devicePath, wait)
	},
}

var volumesDetachCmd = &cobra.Command{
	Use:   "detach [volume-name]",
	Short: "Detach a volume from its instance",
	Long: `Detach a block storage volume from its attached instance. The instance should be
stopped before detaching to prevent data loss.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		volumeName := args[0]
		wait, _ := cmd.Flags().GetBool("wait")

		return detachVolume(cmd.Context(), volumeName, wait)
	},
}

var volumesDeleteCmd = &cobra.Command{
	Use:   "delete [volume-name]",
	Short: "Delete a block storage volume",
	Long: `Delete a block storage volume. The volume must be detached before deletion.
This action is irreversible and will permanently delete all data.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		volumeName := args[0]
		confirm, _ := cmd.Flags().GetBool("confirm")

		return deleteVolume(cmd.Context(), volumeName, confirm)
	},
}

func init() {
	rootCmd.AddCommand(volumesCmd)

	volumesCmd.AddCommand(volumesCreateCmd)
	volumesCmd.AddCommand(volumesListCmd)
	volumesCmd.AddCommand(volumesAttachCmd)
	volumesCmd.AddCommand(volumesDetachCmd)
	volumesCmd.AddCommand(volumesDeleteCmd)

	// Create command flags
	volumesCreateCmd.Flags().StringP("project", "p", "", "Project name for tagging")
	volumesCreateCmd.Flags().StringP("zone", "z", "", "Availability zone (defaults to region + 'a')")
	volumesCreateCmd.Flags().BoolP("wait", "w", false, "Wait for volume to become available")

	// List command flags
	volumesListCmd.Flags().StringP("project", "p", "", "Filter by project name")

	// Attach command flags
	volumesAttachCmd.Flags().BoolP("wait", "w", false, "Wait for volume to attach")

	// Detach command flags
	volumesDetachCmd.Flags().BoolP("wait", "w", false, "Wait for volume to detach")

	// Delete command flags
	volumesDeleteCmd.Flags().BoolP("confirm", "y", false, "Skip confirmation prompt")
}

// createVolume creates a new block storage volume.
func createVolume(ctx context.Context, name, sizeStr, zone, project string, wait bool) error {
	// Parse size
	var sizeGB int32
	_, err := fmt.Sscanf(sizeStr, "%d", &sizeGB)
	if err != nil {
		return fmt.Errorf("invalid size format: %s (expected integer GB)", sizeStr)
	}

	if sizeGB < 8 || sizeGB > 16384 {
		return fmt.Errorf("size must be between 8 GB and 16384 GB (16 TB), got %d GB", sizeGB)
	}

	// Load configuration
	_, err = config.Load()
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

	// Default availability zone
	if zone == "" {
		zone = awsClient.GetRegion() + "a"
	}

	fmt.Printf("Creating volume: %s (%d GB) in %s\n", name, sizeGB, zone)
	if project != "" {
		fmt.Printf("Project: %s\n", project)
	}

	// Create the volume
	disk, err := lightsailService.CreateDisk(ctx, name, sizeGB, zone, project)
	if err != nil {
		return fmt.Errorf("failed to create volume: %w", err)
	}

	fmt.Printf("✅ Volume created: %s\n", disk.Name)
	fmt.Printf("   Size: %d GB, IOPS: %d\n", disk.SizeGB, disk.IOPS)

	// Wait for availability if requested
	if wait {
		err = utils.WaitForDiskState(ctx, disk.Name, "available", func() (string, error) {
			updatedDisk, err := lightsailService.GetDisk(ctx, disk.Name)
			if err != nil {
				return "", err
			}
			return updatedDisk.State, nil
		})
		if err != nil {
			return fmt.Errorf("error waiting for volume: %w", err)
		}
	}

	return nil
}

// listVolumes lists all block storage volumes.
func listVolumes(ctx context.Context, project string) error {
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

	// Get volumes
	disks, err := lightsailService.ListDisks(ctx, project)
	if err != nil {
		return fmt.Errorf("failed to list volumes: %w", err)
	}

	if len(disks) == 0 {
		if project != "" {
			fmt.Printf("No volumes found for project: %s\n", project)
		} else {
			fmt.Println("No volumes found.")
		}
		return nil
	}

	// Display results
	if project != "" {
		fmt.Printf("Block storage volumes for project: %s\n\n", project)
	} else {
		fmt.Printf("All block storage volumes:\n\n")
	}

	fmt.Printf("%-20s %-12s %-8s %-8s %-20s %-12s %-15s\n",
		"VOLUME", "STATE", "SIZE", "IOPS", "ATTACHED TO", "ZONE", "PROJECT")
	fmt.Println(strings.Repeat("-", 120))

	for _, disk := range disks {
		project := disk.Tags["Project"]
		if project == "" {
			project = "untagged"
		}

		attachedTo := disk.AttachedTo
		if attachedTo == "" {
			attachedTo = "-"
		}

		fmt.Printf("%-20s %-12s %-8d %-8d %-20s %-12s %-15s\n",
			disk.Name,
			disk.State,
			disk.SizeGB,
			disk.IOPS,
			attachedTo,
			disk.AvailabilityZone,
			project,
		)
	}

	fmt.Printf("\nTotal: %d volumes\n", len(disks))
	return nil
}

// attachVolume attaches a volume to an instance.
func attachVolume(ctx context.Context, volumeName, instanceName, devicePath string, wait bool) error {
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

	fmt.Printf("Attaching volume %s to instance %s at %s\n", volumeName, instanceName, devicePath)

	err = lightsailService.AttachDisk(ctx, volumeName, instanceName, devicePath)
	if err != nil {
		return fmt.Errorf("failed to attach volume: %w", err)
	}

	fmt.Printf("✅ Volume attached: %s → %s\n", volumeName, instanceName)

	// Wait for attachment if requested
	if wait {
		err = utils.WaitForDiskState(ctx, volumeName, "in-use", func() (string, error) {
			disk, err := lightsailService.GetDisk(ctx, volumeName)
			if err != nil {
				return "", err
			}
			return disk.State, nil
		})
		if err != nil {
			return fmt.Errorf("error waiting for volume attachment: %w", err)
		}
	}

	return nil
}

// detachVolume detaches a volume from its instance.
func detachVolume(ctx context.Context, volumeName string, wait bool) error {
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

	fmt.Printf("Detaching volume: %s\n", volumeName)

	err = lightsailService.DetachDisk(ctx, volumeName)
	if err != nil {
		return fmt.Errorf("failed to detach volume: %w", err)
	}

	fmt.Printf("✅ Volume detached: %s\n", volumeName)

	// Wait for detachment if requested
	if wait {
		err = utils.WaitForDiskState(ctx, volumeName, "available", func() (string, error) {
			disk, err := lightsailService.GetDisk(ctx, volumeName)
			if err != nil {
				return "", err
			}
			return disk.State, nil
		})
		if err != nil {
			return fmt.Errorf("error waiting for volume detachment: %w", err)
		}
	}

	return nil
}

// deleteVolume deletes a block storage volume.
func deleteVolume(ctx context.Context, volumeName string, confirm bool) error {
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

	// Get volume info for confirmation
	disk, err := lightsailService.GetDisk(ctx, volumeName)
	if err != nil {
		return fmt.Errorf("failed to get volume info: %w", err)
	}

	if disk.AttachedTo != "" {
		return fmt.Errorf("volume %s is attached to %s. Detach it first", volumeName, disk.AttachedTo)
	}

	if !confirm {
		fmt.Printf("⚠️  Are you sure you want to delete volume %s (%d GB)?\n", volumeName, disk.SizeGB)
		fmt.Printf("This action is irreversible and will permanently delete all data.\n")
		fmt.Printf("Run with --confirm flag to proceed.\n")
		return nil
	}

	fmt.Printf("Deleting volume: %s (%d GB)\n", volumeName, disk.SizeGB)

	err = lightsailService.DeleteDisk(ctx, volumeName)
	if err != nil {
		return fmt.Errorf("failed to delete volume: %w", err)
	}

	fmt.Printf("✅ Volume deleted: %s\n", volumeName)
	return nil
}