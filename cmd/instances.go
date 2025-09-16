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
	Use:   "resize [instance-name] [new-bundle]",
	Short: "Resize an instance to a different bundle",
	Long: `Change the bundle (CPU, RAM, storage) of a Lightsail instance. The instance
will be stopped during the resize operation and then restarted.`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		instanceName := args[0]
		newBundle := args[1]

		fmt.Printf("Resizing instance '%s' to bundle '%s'\n", instanceName, newBundle)
		fmt.Println("Note: Instance will be stopped during resize operation")

		// TODO: Implement instance resize logic
		return fmt.Errorf("instance resize not yet implemented")
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
	fmt.Printf("%-20s %-15s %-20s %-15s %-15s %-15s %-12s\n",
		"INSTANCE", "STATE", "PUBLIC IP", "BLUEPRINT", "BUNDLE", "REGION", "PROJECT")
	fmt.Println(strings.Repeat("-", 125))

	for _, instance := range instances {
		project := instance.Tags["Project"]
		if project == "" {
			project = "untagged"
		}

		fmt.Printf("%-20s %-15s %-20s %-15s %-15s %-15s %-12s\n",
			instance.Name,
			instance.State,
			instance.PublicIP,
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
				return instance.State, nil
			})
			if err != nil {
				fmt.Printf("❌ Error waiting for instance %s: %v\n", instanceName, err)
			}
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
				return instance.State, nil
			})
			if err != nil {
				fmt.Printf("❌ Error waiting for instance %s: %v\n", instanceName, err)
			}
		}
	}

	fmt.Printf("\n🎉 Instance stop completed!\n")
	return nil
}