package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
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

		if project != "" {
			fmt.Printf("Listing instances for project: %s\n", project)
		}
		if user != "" {
			fmt.Printf("Listing instances for user: %s\n", user)
		}
		if project == "" && user == "" {
			fmt.Println("Listing all instances")
		}

		// TODO: Implement instance listing logic
		return fmt.Errorf("instance listing not yet implemented")
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

		fmt.Printf("Starting instances for users: %v\n", users)
		if project != "" {
			fmt.Printf("Project filter: %s\n", project)
		}

		// TODO: Implement instance start logic
		return fmt.Errorf("instance start not yet implemented")
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

		fmt.Printf("Stopping instances for users: %v\n", users)
		if project != "" {
			fmt.Printf("Project filter: %s\n", project)
		}

		// TODO: Implement instance stop logic
		return fmt.Errorf("instance stop not yet implemented")
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

func init() {
	rootCmd.AddCommand(instancesCmd)

	instancesCmd.AddCommand(instancesListCmd)
	instancesCmd.AddCommand(instancesStartCmd)
	instancesCmd.AddCommand(instancesStopCmd)
	instancesCmd.AddCommand(instancesMonitorCmd)

	// List command flags
	instancesListCmd.Flags().StringP("project", "p", "", "Filter by project name")
	instancesListCmd.Flags().StringP("user", "u", "", "Filter by username")

	// Start command flags
	instancesStartCmd.Flags().StringSliceP("users", "u", []string{}, "Comma-separated list of usernames (required)")
	instancesStartCmd.Flags().StringP("project", "p", "", "Filter by project name")
	instancesStartCmd.MarkFlagRequired("users")

	// Stop command flags
	instancesStopCmd.Flags().StringSliceP("users", "u", []string{}, "Comma-separated list of usernames (required)")
	instancesStopCmd.Flags().StringP("project", "p", "", "Filter by project name")
	instancesStopCmd.MarkFlagRequired("users")

	// Monitor command flags
	instancesMonitorCmd.Flags().StringP("project", "p", "", "Filter by project name")
	instancesMonitorCmd.Flags().IntP("idle-threshold", "t", 120, "Idle threshold in minutes")
}