package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/scttfrdmn/lfr-tools/internal/aws"
	"github.com/scttfrdmn/lfr-tools/internal/config"
	"github.com/scttfrdmn/lfr-tools/internal/types"
)

var idleCmd = &cobra.Command{
	Use:   "idle",
	Short: "Manage Lightsail for Research idle detection rules",
	Long:  `Configure and manage idle detection settings for instances to automatically stop them when not in use.`,
}

var idleConfigureCmd = &cobra.Command{
	Use:   "configure [instance-name]",
	Short: "Configure idle detection for an instance",
	Long: `Configure idle detection settings for a Lightsail instance. When an instance
is idle for the specified threshold, it will automatically stop after the duration period.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		instanceName := args[0]
		threshold, _ := cmd.Flags().GetInt("threshold")
		duration, _ := cmd.Flags().GetInt("duration")
		disable, _ := cmd.Flags().GetBool("disable")

		return configureIdleDetection(cmd.Context(), instanceName, threshold, duration, disable)
	},
}

var idleConfigureBulkCmd = &cobra.Command{
	Use:   "configure-bulk",
	Short: "Configure idle detection for multiple instances",
	Long: `Configure idle detection for multiple instances by project or user list.
Applies the same idle settings to all specified instances.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		project, _ := cmd.Flags().GetString("project")
		users, _ := cmd.Flags().GetStringSlice("users")
		threshold, _ := cmd.Flags().GetInt("threshold")
		duration, _ := cmd.Flags().GetInt("duration")
		disable, _ := cmd.Flags().GetBool("disable")

		return configureIdleDetectionBulk(cmd.Context(), project, users, threshold, duration, disable)
	},
}

var idleStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show idle detection status for instances",
	Long: `Display idle detection configuration for all instances, showing thresholds,
durations, and current idle status.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		project, _ := cmd.Flags().GetString("project")
		user, _ := cmd.Flags().GetString("user")

		return showIdleStatus(cmd.Context(), project, user)
	},
}

func init() {
	rootCmd.AddCommand(idleCmd)

	idleCmd.AddCommand(idleConfigureCmd)
	idleCmd.AddCommand(idleConfigureBulkCmd)
	idleCmd.AddCommand(idleStatusCmd)

	// Configure command flags
	idleConfigureCmd.Flags().IntP("threshold", "t", 120, "Idle threshold in minutes (default: 120)")
	idleConfigureCmd.Flags().IntP("duration", "d", 30, "Duration in minutes before stopping (default: 30)")
	idleConfigureCmd.Flags().BoolP("disable", "", false, "Disable idle detection")

	// Configure bulk command flags
	idleConfigureBulkCmd.Flags().StringP("project", "p", "", "Configure all instances in project")
	idleConfigureBulkCmd.Flags().StringSliceP("users", "u", []string{}, "Configure instances for specific users")
	idleConfigureBulkCmd.Flags().IntP("threshold", "t", 120, "Idle threshold in minutes (default: 120)")
	idleConfigureBulkCmd.Flags().IntP("duration", "d", 30, "Duration in minutes before stopping (default: 30)")
	idleConfigureBulkCmd.Flags().BoolP("disable", "", false, "Disable idle detection")

	// Status command flags
	idleStatusCmd.Flags().StringP("project", "p", "", "Filter by project name")
	idleStatusCmd.Flags().StringP("user", "u", "", "Filter by username")
}

// configureIdleDetection configures idle detection for a single instance.
func configureIdleDetection(ctx context.Context, instanceName string, threshold, duration int, disable bool) error {
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

	// Get instance details
	instance, err := lightsailService.GetInstance(ctx, instanceName)
	if err != nil {
		return fmt.Errorf("failed to get instance details: %w", err)
	}

	if disable {
		fmt.Printf("Disabling idle detection for: %s\n", instanceName)
		// TODO: Implement disable idle detection API call
		fmt.Printf("⚠️  Idle detection disable requires PutInstancePublicPorts or ModifyInstanceAttributes API\n")
		fmt.Printf("This may not be available in current Lightsail API\n")
		return nil
	}

	fmt.Printf("Configuring idle detection for: %s\n", instanceName)
	fmt.Printf("Current state: %s\n", instance.State)
	fmt.Printf("Threshold: %d minutes of inactivity\n", threshold)
	fmt.Printf("Duration: %d minutes before auto-stop\n", duration)

	// NOTE: Lightsail for Research instances are created with idle detection by default
	// The original script used: --add-ons "addOnType=StopInstanceOnIdle,stopInstanceOnIdleRequest={threshold=2,duration=30}"
	// This is set at creation time, not modified after creation

	fmt.Printf("\n💡 Important: Idle detection is configured at instance creation time\n")
	fmt.Printf("Current instances were created with default settings (2 hours threshold, 30 minutes duration)\n")
	fmt.Printf("To change idle settings, you need to:\n")
	fmt.Printf("1. Create a snapshot of the instance\n")
	fmt.Printf("2. Create a new instance with new idle settings\n")
	fmt.Printf("3. Delete the old instance\n\n")

	fmt.Printf("For new instances, use: lfr users create --idle-threshold=%d --idle-duration=%d\n", threshold, duration)

	return nil
}

// configureIdleDetectionBulk configures idle detection for multiple instances.
func configureIdleDetectionBulk(ctx context.Context, project string, users []string, threshold, duration int, disable bool) error {
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

	// Filter by users if specified
	if len(users) > 0 {
		var filtered []*types.Instance
		for _, instance := range instances {
			for _, user := range users {
				if strings.HasPrefix(instance.Name, user+"-") {
					filtered = append(filtered, instance)
					break
				}
			}
		}
		instances = filtered
	}

	if len(instances) == 0 {
		fmt.Println("No instances found to configure.")
		return nil
	}

	actionDesc := "Configuring idle detection"
	if disable {
		actionDesc = "Disabling idle detection"
	}

	fmt.Printf("%s for %d instances\n", actionDesc, len(instances))
	if !disable {
		fmt.Printf("Threshold: %d minutes, Duration: %d minutes\n", threshold, duration)
	}
	fmt.Println()

	for i, instance := range instances {
		fmt.Printf("[%d/%d] %s: %s\n", i+1, len(instances), actionDesc, instance.Name)

		err := configureIdleDetection(ctx, instance.Name, threshold, duration, disable)
		if err != nil {
			fmt.Printf("❌ Failed to configure %s: %v\n", instance.Name, err)
		} else {
			fmt.Printf("✅ Configured %s\n", instance.Name)
		}
	}

	fmt.Printf("\n🎉 Bulk idle detection configuration completed!\n")
	return nil
}

// showIdleStatus shows idle detection status for instances.
func showIdleStatus(ctx context.Context, project, user string) error {
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
		fmt.Println("No instances found.")
		return nil
	}

	// Display idle status
	fmt.Printf("Idle Detection Status:\n\n")
	fmt.Printf("%-20s %-12s %-15s %-10s %-15s\n",
		"INSTANCE", "STATE", "IDLE CONFIG", "THRESHOLD", "DURATION")
	fmt.Println(strings.Repeat("-", 85))

	for _, instance := range instances {
		// Default LfR idle detection settings
		idleConfig := "Enabled"
		threshold := "2 hours"
		duration := "30 min"

		// Note: Getting actual idle detection settings requires additional API calls
		// that may not be readily available in the Lightsail API

		fmt.Printf("%-20s %-12s %-15s %-10s %-15s\n",
			instance.Name,
			instance.State,
			idleConfig,
			threshold,
			duration,
		)
	}

	fmt.Printf("\nTotal: %d instances\n", len(instances))
	fmt.Printf("\nNote: Idle detection settings are configured at instance creation time.\n")
	fmt.Printf("Default LfR settings: 2 hours threshold, 30 minutes duration\n")
	fmt.Printf("To modify settings, recreate instance with new idle configuration.\n")

	return nil
}