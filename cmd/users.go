package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "Manage IAM users for Lightsail for Research",
	Long:  `Create, remove, and list IAM users with appropriate Lightsail for Research permissions.`,
}

var usersCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create IAM users and Lightsail instances",
	Long: `Create IAM users with auto-generated passwords and provision Lightsail instances
for each user with appropriate access controls.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		project, _ := cmd.Flags().GetString("project")
		blueprint, _ := cmd.Flags().GetString("blueprint")
		bundle, _ := cmd.Flags().GetString("bundle")
		region, _ := cmd.Flags().GetString("region")
		users, _ := cmd.Flags().GetStringSlice("users")

		fmt.Printf("Creating users for project: %s\n", project)
		fmt.Printf("Blueprint: %s, Bundle: %s, Region: %s\n", blueprint, bundle, region)
		fmt.Printf("Users: %v\n", users)

		// TODO: Implement user creation logic
		return fmt.Errorf("user creation not yet implemented")
	},
}

var usersRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove IAM users and their Lightsail instances",
	Long: `Remove IAM users and their associated Lightsail instances. This action is irreversible
and will delete all data on the instances.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		project, _ := cmd.Flags().GetString("project")
		users, _ := cmd.Flags().GetStringSlice("users")
		all, _ := cmd.Flags().GetBool("all")

		fmt.Printf("Removing users from project: %s\n", project)
		if all {
			fmt.Println("Removing ALL users in project")
		} else {
			fmt.Printf("Users to remove: %v\n", users)
		}

		// TODO: Implement user removal logic
		return fmt.Errorf("user removal not yet implemented")
	},
}

var usersListCmd = &cobra.Command{
	Use:   "list",
	Short: "List IAM users and their instances",
	Long:  `List all IAM users and their associated Lightsail instances, optionally filtered by project.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		project, _ := cmd.Flags().GetString("project")

		if project != "" {
			fmt.Printf("Listing users for project: %s\n", project)
		} else {
			fmt.Println("Listing all users")
		}

		// TODO: Implement user listing logic
		return fmt.Errorf("user listing not yet implemented")
	},
}

func init() {
	rootCmd.AddCommand(usersCmd)

	usersCmd.AddCommand(usersCreateCmd)
	usersCmd.AddCommand(usersRemoveCmd)
	usersCmd.AddCommand(usersListCmd)

	// Create command flags
	usersCreateCmd.Flags().StringP("project", "p", "", "Project name (required)")
	usersCreateCmd.Flags().StringP("blueprint", "b", "", "Lightsail blueprint ID (required)")
	usersCreateCmd.Flags().String("bundle", "", "Lightsail bundle ID (required)")
	usersCreateCmd.Flags().StringP("region", "r", "", "AWS region (required)")
	usersCreateCmd.Flags().StringSliceP("users", "u", []string{}, "Comma-separated list of usernames (required)")

	usersCreateCmd.MarkFlagRequired("project")
	usersCreateCmd.MarkFlagRequired("blueprint")
	usersCreateCmd.MarkFlagRequired("bundle")
	usersCreateCmd.MarkFlagRequired("region")
	usersCreateCmd.MarkFlagRequired("users")

	// Remove command flags
	usersRemoveCmd.Flags().StringP("project", "p", "", "Project name (required)")
	usersRemoveCmd.Flags().StringSliceP("users", "u", []string{}, "Comma-separated list of usernames to remove")
	usersRemoveCmd.Flags().Bool("all", false, "Remove all users in the project")

	usersRemoveCmd.MarkFlagRequired("project")

	// List command flags
	usersListCmd.Flags().StringP("project", "p", "", "Filter by project name")
}