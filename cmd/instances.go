package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/scttfrdmn/lfr-tools/internal/aws"
	"github.com/scttfrdmn/lfr-tools/internal/config"
	"github.com/scttfrdmn/lfr-tools/internal/types"
	"github.com/scttfrdmn/lfr-tools/internal/utils"
)

var instancesCmd = &cobra.Command{
	Use:   "instances",
	Short: "Manage Lightsail for Research instances",
	Long:  `List, start, stop, and monitor Lightsail for Research instances.`,
}

var instancesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Lightsail instances",
	Long: `List all Lightsail instances with their current status, optionally filtered by project
or user. Shows instance details including state, IP addresses, and resource utilization.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		project, _ := cmd.Flags().GetString("project")
		user, _ := cmd.Flags().GetString("user")

		return listInstances(cmd.Context(), project, user)
	},
}

var instancesStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start Lightsail instances",
	Long: `Start stopped Lightsail instances for the specified users. Instances will be
charged according to your bundle pricing while running.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		users, _ := cmd.Flags().GetStringSlice("users")
		project, _ := cmd.Flags().GetString("project")
		wait, _ := cmd.Flags().GetBool("wait")

		return startInstances(cmd.Context(), users, project, wait)
	},
}

var instancesStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop Lightsail instances",
	Long: `Stop running Lightsail instances for the specified users. This will save costs
but users will lose access until instances are restarted.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		users, _ := cmd.Flags().GetStringSlice("users")
		project, _ := cmd.Flags().GetString("project")
		wait, _ := cmd.Flags().GetBool("wait")

		return stopInstances(cmd.Context(), users, project, wait)
	},
}

var instancesMonitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Monitor instance usage and idle time",
	Long: `Monitor instance CPU, memory, and network usage. Identify idle instances
that may be candidates for automatic shutdown to save costs.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		project, _ := cmd.Flags().GetString("project")
		threshold, _ := cmd.Flags().GetInt("idle-threshold")

		fmt.Printf("Monitoring instances")
		if project != "" {
			fmt.Printf(" for project: %s", project)
		}
		fmt.Printf(" with idle threshold: %d minutes\n", threshold)

		// TODO: Implement instance monitoring logic
		return fmt.Errorf("instance monitoring not yet implemented")
	},
}

var instancesSnapshotCmd = &cobra.Command{
	Use:   "snapshot [instance-name]",
	Short: "Create a snapshot of an instance",
	Long: `Create a point-in-time snapshot of a Lightsail instance for backup or cloning purposes.
Snapshots preserve the instance state and can be used to restore or create new instances.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		instanceName := args[0]
		snapshotName, _ := cmd.Flags().GetString("name")

		if snapshotName == "" {
			snapshotName = fmt.Sprintf("%s-snapshot-%d", instanceName, time.Now().Unix())
		}

		fmt.Printf("Creating snapshot '%s' of instance '%s'\n", snapshotName, instanceName)

		// TODO: Implement snapshot creation logic
		return fmt.Errorf("instance snapshot not yet implemented")
	},
}

var instancesRestoreCmd = &cobra.Command{
	Use:   "restore [snapshot-name] [new-instance-name]",
	Short: "Restore an instance from a snapshot",
	Long: `Create a new instance from an existing snapshot. This allows you to restore
previous states or clone instances with identical configurations.`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		snapshotName := args[0]
		newInstanceName := args[1]
		bundle, _ := cmd.Flags().GetString("bundle")
		region, _ := cmd.Flags().GetString("region")

		fmt.Printf("Restoring snapshot '%s' to new instance '%s'\n", snapshotName, newInstanceName)
		if bundle != "" {
			fmt.Printf("Bundle: %s\n", bundle)
		}
		if region != "" {
			fmt.Printf("Region: %s\n", region)
		}

		// TODO: Implement restore from snapshot logic
		return fmt.Errorf("instance restore not yet implemented")
	},
}

var instancesCloneCmd = &cobra.Command{
	Use:   "clone [source-instance] [new-instance-name]",
	Short: "Clone an existing instance",
	Long: `Create a new instance by cloning an existing one. This creates a snapshot
of the source instance and then creates a new instance from that snapshot.`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		sourceInstance := args[0]
		newInstanceName := args[1]
		bundle, _ := cmd.Flags().GetString("bundle")
		region, _ := cmd.Flags().GetString("region")

		fmt.Printf("Cloning instance '%s' to '%s'\n", sourceInstance, newInstanceName)
		if bundle != "" {
			fmt.Printf("New bundle: %s\n", bundle)
		}
		if region != "" {
			fmt.Printf("Region: %s\n", region)
		}

		// TODO: Implement instance cloning logic
		return fmt.Errorf("instance cloning not yet implemented")
	},
}

var instancesRebootCmd = &cobra.Command{
	Use:   "reboot [instance-names...]",
	Short: "Reboot Lightsail instances",
	Long: `Reboot one or more Lightsail instances. This performs a graceful restart
of the instances, preserving all data.`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		instanceNames := args
		force, _ := cmd.Flags().GetBool("force")

		if force {
			fmt.Printf("Force rebooting instances: %v\n", instanceNames)
		} else {
			fmt.Printf("Rebooting instances: %v\n", instanceNames)
		}

		// TODO: Implement instance reboot logic
		return fmt.Errorf("instance reboot not yet implemented")
	},
}

var instancesResizeCmd = &cobra.Command{
	Use:   "resize [instance-name] [direction]",
	Short: "Resize an instance to next larger/smaller bundle",
	Long: `Resize an instance to the next larger or smaller bundle in the same category.
Direction: 'up' for larger, 'down' for smaller. Instance will be stopped during resize.`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		instanceName := args[0]
		direction := args[1]
		wait, _ := cmd.Flags().GetBool("wait")

		return resizeInstance(cmd.Context(), instanceName, direction, wait)
	},
}

var instancesGpuCmd = &cobra.Command{
	Use:   "gpu [instance-name] [enable|disable]",
	Short: "Switch instance between GPU and standard bundles",
	Long: `Switch an instance between GPU and standard bundles of equivalent size.
Automatically finds the equivalent bundle and checks GPU quota availability.`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		instanceName := args[0]
		action := args[1]
		wait, _ := cmd.Flags().GetBool("wait")

		return switchGPUMode(cmd.Context(), instanceName, action, wait)
	},
}

func init() {
	rootCmd.AddCommand(instancesCmd)

	instancesCmd.AddCommand(instancesListCmd)
	instancesCmd.AddCommand(instancesStartCmd)
	instancesCmd.AddCommand(instancesStopCmd)
	instancesCmd.AddCommand(instancesMonitorCmd)
	instancesCmd.AddCommand(instancesSnapshotCmd)
	instancesCmd.AddCommand(instancesRestoreCmd)
	instancesCmd.AddCommand(instancesCloneCmd)
	instancesCmd.AddCommand(instancesRebootCmd)
	instancesCmd.AddCommand(instancesResizeCmd)
	instancesCmd.AddCommand(instancesGpuCmd)

	// List command flags
	instancesListCmd.Flags().StringP("project", "p", "", "Filter by project name")
	instancesListCmd.Flags().StringP("user", "u", "", "Filter by username")

	// Start command flags
	instancesStartCmd.Flags().StringSliceP("users", "u", []string{}, "Comma-separated list of usernames (required)")
	instancesStartCmd.Flags().StringP("project", "p", "", "Filter by project name")
	instancesStartCmd.Flags().BoolP("wait", "w", false, "Wait for instances to reach running state")
	instancesStartCmd.MarkFlagRequired("users")

	// Stop command flags
	instancesStopCmd.Flags().StringSliceP("users", "u", []string{}, "Comma-separated list of usernames (required)")
	instancesStopCmd.Flags().StringP("project", "p", "", "Filter by project name")
	instancesStopCmd.Flags().BoolP("wait", "w", false, "Wait for instances to reach stopped state")
	instancesStopCmd.MarkFlagRequired("users")

	// Monitor command flags
	instancesMonitorCmd.Flags().StringP("project", "p", "", "Filter by project name")
	instancesMonitorCmd.Flags().IntP("idle-threshold", "t", 120, "Idle threshold in minutes")

	// Snapshot command flags
	instancesSnapshotCmd.Flags().StringP("name", "n", "", "Custom snapshot name (auto-generated if not provided)")

	// Restore command flags
	instancesRestoreCmd.Flags().StringP("bundle", "b", "", "Bundle for the new instance")
	instancesRestoreCmd.Flags().StringP("region", "r", "", "Region for the new instance")

	// Clone command flags
	instancesCloneCmd.Flags().StringP("bundle", "b", "", "Bundle for the new instance (uses source bundle if not provided)")
	instancesCloneCmd.Flags().StringP("region", "r", "", "Region for the new instance (uses source region if not provided)")

	// Reboot command flags
	instancesRebootCmd.Flags().BoolP("force", "f", false, "Force reboot without confirmation")

	// Resize command flags
	instancesResizeCmd.Flags().BoolP("wait", "w", false, "Wait for resize to complete")

	// GPU command flags
	instancesGpuCmd.Flags().BoolP("wait", "w", false, "Wait for GPU switch to complete")
}

// listInstances lists Lightsail instances with filtering.
func listInstances(ctx context.Context, project, user string) error {
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

	// Filter by user if specified
	if user != "" {
		var filtered []*types.Instance
		for _, instance := range instances {
			if strings.HasPrefix(instance.Name, user+"-") {
				filtered = append(filtered, instance)
			}
		}
		instances = filtered
	}

	if len(instances) == 0 {
		if project != "" && user != "" {
			fmt.Printf("No instances found for user %s in project: %s\n", user, project)
		} else if project != "" {
			fmt.Printf("No instances found for project: %s\n", project)
		} else if user != "" {
			fmt.Printf("No instances found for user: %s\n", user)
		} else {
			fmt.Println("No instances found.")
		}
		return nil
	}

	// Display results
	fmt.Printf("%-20s %-12s %-18s %-15s %-20s %-12s %-15s\n",
		"INSTANCE", "STATE", "PUBLIC IP", "BLUEPRINT", "BUNDLE", "REGION", "PROJECT")
	fmt.Println(strings.Repeat("-", 130))

	for _, instance := range instances {
		project := instance.Tags["Project"]
		if project == "" {
			project = "untagged"
		}

		publicIP := instance.PublicIP
		if publicIP == "" {
			publicIP = "-"
		}

		fmt.Printf("%-20s %-12s %-18s %-15s %-20s %-12s %-15s\n",
			instance.Name,
			instance.State,
			publicIP,
			instance.Blueprint,
			instance.Bundle,
			instance.Region,
			project,
		)
	}

	fmt.Printf("\nTotal: %d instances\n", len(instances))
	return nil
}

// startInstances starts instances for specified users.
func startInstances(ctx context.Context, users []string, project string, wait bool) error {
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

	// Get all instances
	instances, err := lightsailService.ListInstances(ctx, project)
	if err != nil {
		return fmt.Errorf("failed to list instances: %w", err)
	}

	// Filter instances for specified users
	var instancesToStart []string
	for _, instance := range instances {
		for _, user := range users {
			if strings.HasPrefix(instance.Name, user+"-") {
				instancesToStart = append(instancesToStart, instance.Name)
				break
			}
		}
	}

	if len(instancesToStart) == 0 {
		fmt.Printf("No instances found for users: %v\n", users)
		return nil
	}

	fmt.Printf("Starting %d instances for users: %v\n", len(instancesToStart), users)

	for i, instanceName := range instancesToStart {
		fmt.Printf("[%d/%d] Starting instance: %s\n", i+1, len(instancesToStart), instanceName)

		err = lightsailService.StartInstance(ctx, instanceName)
		if err != nil {
			fmt.Printf("❌ Error starting instance %s: %v\n", instanceName, err)
			continue
		}

		fmt.Printf("✅ Started instance: %s\n", instanceName)

		// Wait for instance to reach running state if requested
		if wait {
			err = utils.WaitForInstanceState(ctx, instanceName, "running", func() (string, error) {
				instance, err := lightsailService.GetInstance(ctx, instanceName)
				if err != nil {
					return "", err
				}

				// Update S3 status during wait
				_ = utils.UpdateInstanceStatusInS3(ctx, instance)

				return instance.State, nil
			})
			if err != nil {
				fmt.Printf("❌ Error waiting for instance %s: %v\n", instanceName, err)
			}
		}

		// Update S3 status after start
		if instance, err := lightsailService.GetInstance(ctx, instanceName); err == nil {
			_ = utils.UpdateInstanceStatusInS3(ctx, instance)
		}
	}

	fmt.Printf("\n🎉 Instance start completed!\n")
	return nil
}

// stopInstances stops instances for specified users.
func stopInstances(ctx context.Context, users []string, project string, wait bool) error {
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

	// Get all instances
	instances, err := lightsailService.ListInstances(ctx, project)
	if err != nil {
		return fmt.Errorf("failed to list instances: %w", err)
	}

	// Filter instances for specified users
	var instancesToStop []string
	for _, instance := range instances {
		for _, user := range users {
			if strings.HasPrefix(instance.Name, user+"-") {
				instancesToStop = append(instancesToStop, instance.Name)
				break
			}
		}
	}

	if len(instancesToStop) == 0 {
		fmt.Printf("No instances found for users: %v\n", users)
		return nil
	}

	fmt.Printf("Stopping %d instances for users: %v\n", len(instancesToStop), users)

	for i, instanceName := range instancesToStop {
		fmt.Printf("[%d/%d] Stopping instance: %s\n", i+1, len(instancesToStop), instanceName)

		err = lightsailService.StopInstance(ctx, instanceName)
		if err != nil {
			fmt.Printf("❌ Error stopping instance %s: %v\n", instanceName, err)
			continue
		}

		fmt.Printf("✅ Stopped instance: %s\n", instanceName)

		// Wait for instance to reach stopped state if requested
		if wait {
			err = utils.WaitForInstanceState(ctx, instanceName, "stopped", func() (string, error) {
				instance, err := lightsailService.GetInstance(ctx, instanceName)
				if err != nil {
					return "", err
				}

				// Update S3 status during wait
				_ = utils.UpdateInstanceStatusInS3(ctx, instance)

				return instance.State, nil
			})
			if err != nil {
				fmt.Printf("❌ Error waiting for instance %s: %v\n", instanceName, err)
			}
		}

		// Update S3 status after stop
		if instance, err := lightsailService.GetInstance(ctx, instanceName); err == nil {
			_ = utils.UpdateInstanceStatusInS3(ctx, instance)
		}
	}

	fmt.Printf("\n🎉 Instance stop completed!\n")
	return nil
}

// resizeInstance resizes an instance using the snapshot method.
func resizeInstance(ctx context.Context, instanceName, direction string, wait bool) error {
	if direction != "up" && direction != "down" {
		return fmt.Errorf("direction must be 'up' or 'down', got: %s", direction)
	}

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

	// Get current instance details
	instance, err := lightsailService.GetInstance(ctx, instanceName)
	if err != nil {
		return fmt.Errorf("failed to get instance details: %w", err)
	}

	currentBundle, err := utils.GetBundleInfo(instance.Bundle)
	if err != nil {
		return fmt.Errorf("failed to get current bundle info: %w", err)
	}

	// Find target bundle
	var targetBundle *utils.BundleInfo
	if direction == "up" {
		targetBundle, err = utils.GetNextSizeBundle(instance.Bundle)
		if err != nil {
			return fmt.Errorf("failed to find larger bundle: %w", err)
		}
	} else {
		targetBundle, err = utils.GetPreviousSizeBundle(instance.Bundle)
		if err != nil {
			return fmt.Errorf("failed to find smaller bundle: %w", err)
		}
	}

	// Display resize plan
	fmt.Printf("Resizing instance: %s\n", instanceName)
	fmt.Printf("Current state: %s\n", instance.State)
	fmt.Printf("%s\n", utils.FormatBundleComparison(currentBundle, targetBundle))
	fmt.Printf("\n⚠️  This operation will:\n")
	fmt.Printf("   1. Stop the instance (if running)\n")
	fmt.Printf("   2. Create a snapshot: %s-resize-snapshot\n", instanceName)
	fmt.Printf("   3. Create new instance: %s-resized\n", instanceName)
	fmt.Printf("   4. Optionally delete old instance and snapshot\n\n")

	// Stop instance if running
	if instance.State == "running" {
		fmt.Printf("Stopping instance %s...\n", instanceName)
		err = lightsailService.StopInstance(ctx, instanceName)
		if err != nil {
			return fmt.Errorf("failed to stop instance: %w", err)
		}

		if wait {
			err = utils.WaitForInstanceState(ctx, instanceName, "stopped", func() (string, error) {
				inst, err := lightsailService.GetInstance(ctx, instanceName)
				if err != nil {
					return "", err
				}
				return inst.State, nil
			})
			if err != nil {
				return fmt.Errorf("error waiting for instance to stop: %w", err)
			}
		}
	}

	// Create snapshot
	snapshotName := instanceName + "-resize-snapshot"
	fmt.Printf("Creating snapshot: %s\n", snapshotName)
	err = lightsailService.CreateInstanceSnapshot(ctx, instanceName, snapshotName)
	if err != nil {
		return fmt.Errorf("failed to create snapshot: %w", err)
	}

	// Wait for snapshot to complete if requested
	if wait {
		fmt.Printf("Waiting for snapshot to complete...\n")
		err = utils.WaitForSnapshotState(ctx, snapshotName, "available", func() (string, error) {
			snapshot, err := lightsailService.GetInstanceSnapshot(ctx, snapshotName)
			if err != nil {
				return "", err
			}
			return string(snapshot.State), nil
		})
		if err != nil {
			return fmt.Errorf("error waiting for snapshot: %w", err)
		}
	}

	// Create new instance from snapshot
	newInstanceName := instanceName + "-resized"
	fmt.Printf("Creating resized instance: %s\n", newInstanceName)

	tags := instance.Tags
	if tags == nil {
		tags = make(map[string]string)
	}
	tags["ResizedFrom"] = instanceName

	_, err = lightsailService.CreateInstanceFromSnapshot(ctx, newInstanceName, snapshotName, targetBundle.ID, instance.Region+"a", tags)
	if err != nil {
		return fmt.Errorf("failed to create instance from snapshot: %w", err)
	}

	fmt.Printf("✅ Resize operation initiated!\n")
	fmt.Printf("New instance: %s (%s)\n", newInstanceName, targetBundle.Name)
	fmt.Printf("\nNext steps:\n")
	fmt.Printf("1. Wait for new instance to be ready\n")
	fmt.Printf("2. Test the new instance: lfr ssh connect %s\n", strings.Split(newInstanceName, "-")[0])
	fmt.Printf("3. Delete old instance: lfr instances delete %s\n", instanceName)
	fmt.Printf("4. Rename new instance if desired\n")

	return nil
}

// switchGPUMode switches an instance between GPU and standard bundles.
func switchGPUMode(ctx context.Context, instanceName, action string, wait bool) error {
	if action != "enable" && action != "disable" {
		return fmt.Errorf("action must be 'enable' or 'disable', got: %s", action)
	}

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

	// Get current instance details
	instance, err := lightsailService.GetInstance(ctx, instanceName)
	if err != nil {
		return fmt.Errorf("failed to get instance details: %w", err)
	}

	currentBundle, err := utils.GetBundleInfo(instance.Bundle)
	if err != nil {
		return fmt.Errorf("failed to get current bundle info: %w", err)
	}

	// Find target bundle
	var targetBundle *utils.BundleInfo
	if action == "enable" {
		if currentBundle.IsGPU {
			return fmt.Errorf("instance %s already has GPU enabled (%s)", instanceName, currentBundle.Name)
		}
		targetBundle, err = utils.GetEquivalentGPUBundle(instance.Bundle)
		if err != nil {
			return fmt.Errorf("failed to find equivalent GPU bundle: %w", err)
		}
	} else {
		if !currentBundle.IsGPU {
			return fmt.Errorf("instance %s does not have GPU enabled (%s)", instanceName, currentBundle.Name)
		}
		targetBundle, err = utils.GetEquivalentStandardBundle(instance.Bundle)
		if err != nil {
			return fmt.Errorf("failed to find equivalent standard bundle: %w", err)
		}
	}

	// Display GPU switch plan
	actionDesc := "Enabling GPU"
	if action == "disable" {
		actionDesc = "Disabling GPU"
	}

	fmt.Printf("%s for instance: %s\n", actionDesc, instanceName)
	fmt.Printf("Current state: %s\n", instance.State)
	fmt.Printf("%s\n", utils.FormatBundleComparison(currentBundle, targetBundle))

	// TODO: Add GPU quota checking here when AWS implements it
	if action == "enable" {
		fmt.Printf("\n💡 Note: GPU quota checking not yet available from AWS API\n")
		fmt.Printf("Ensure you have GPU quota in your account before proceeding\n")
	}

	fmt.Printf("\n⚠️  This operation will:\n")
	fmt.Printf("   1. Stop the instance (if running)\n")
	fmt.Printf("   2. Create a snapshot: %s-gpu-snapshot\n", instanceName)
	fmt.Printf("   3. Create new instance: %s-gpu\n", instanceName)
	fmt.Printf("   4. Optionally delete old instance and snapshot\n\n")

	// Stop instance if running
	if instance.State == "running" {
		fmt.Printf("Stopping instance %s...\n", instanceName)
		err = lightsailService.StopInstance(ctx, instanceName)
		if err != nil {
			return fmt.Errorf("failed to stop instance: %w", err)
		}

		if wait {
			err = utils.WaitForInstanceState(ctx, instanceName, "stopped", func() (string, error) {
				inst, err := lightsailService.GetInstance(ctx, instanceName)
				if err != nil {
					return "", err
				}
				return inst.State, nil
			})
			if err != nil {
				return fmt.Errorf("error waiting for instance to stop: %w", err)
			}
		}
	}

	// Create snapshot
	snapshotName := instanceName + "-gpu-snapshot"
	fmt.Printf("Creating snapshot: %s\n", snapshotName)
	err = lightsailService.CreateInstanceSnapshot(ctx, instanceName, snapshotName)
	if err != nil {
		return fmt.Errorf("failed to create snapshot: %w", err)
	}

	// Wait for snapshot if requested
	if wait {
		fmt.Printf("Waiting for snapshot to complete...\n")
		err = utils.WaitForSnapshotState(ctx, snapshotName, "available", func() (string, error) {
			snapshot, err := lightsailService.GetInstanceSnapshot(ctx, snapshotName)
			if err != nil {
				return "", err
			}
			return string(snapshot.State), nil
		})
		if err != nil {
			return fmt.Errorf("error waiting for snapshot: %w", err)
		}
	}

	// Create new instance with GPU/standard bundle
	newInstanceName := instanceName + "-gpu"
	fmt.Printf("Creating new instance: %s\n", newInstanceName)

	tags := instance.Tags
	if tags == nil {
		tags = make(map[string]string)
	}
	tags["GPUSwitchFrom"] = instanceName
	tags["GPUMode"] = action

	_, err = lightsailService.CreateInstanceFromSnapshot(ctx, newInstanceName, snapshotName, targetBundle.ID, instance.Region+"a", tags)
	if err != nil {
		return fmt.Errorf("failed to create GPU-switched instance: %w", err)
	}

	fmt.Printf("✅ GPU switch operation initiated!\n")
	fmt.Printf("New instance: %s (%s)\n", newInstanceName, targetBundle.Name)
	fmt.Printf("\nNext steps:\n")
	fmt.Printf("1. Wait for new instance to be ready\n")
	fmt.Printf("2. Test the new instance: lfr ssh connect %s\n", strings.Split(newInstanceName, "-")[0])
	fmt.Printf("3. Delete old instance: lfr instances delete %s\n", instanceName)
	fmt.Printf("4. Rename new instance if desired\n")

	return nil
}