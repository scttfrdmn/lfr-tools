package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var groupsCmd = &cobra.Command{
	Use:   "groups",
	Short: "Manage IAM groups for Lightsail for Research",
	Long:  `Create, remove, and list IAM groups with custom policies for organized user management.`,
}

var groupsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an IAM group with policies",
	Long: `Create an IAM group and attach the specified policies to it. This allows for
organized user management and consistent permission sets.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		policies, _ := cmd.Flags().GetStringSlice("policies")
		description, _ := cmd.Flags().GetString("description")

		fmt.Printf("Creating group: %s\n", name)
		fmt.Printf("Description: %s\n", description)
		fmt.Printf("Policies: %v\n", policies)

		// TODO: Implement group creation logic
		return fmt.Errorf("group creation not yet implemented")
	},
}

var groupsRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove an IAM group",
	Long: `Remove an IAM group and detach all policies. Users in the group will lose
the group's permissions but will not be deleted.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")

		fmt.Printf("Removing group: %s\n", name)

		// TODO: Implement group removal logic
		return fmt.Errorf("group removal not yet implemented")
	},
}

var groupsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all IAM groups",
	Long:  `List all IAM groups with their attached policies and member counts.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Listing all groups")

		// TODO: Implement group listing logic
		return fmt.Errorf("group listing not yet implemented")
	},
}

func init() {
	rootCmd.AddCommand(groupsCmd)

	groupsCmd.AddCommand(groupsCreateCmd)
	groupsCmd.AddCommand(groupsRemoveCmd)
	groupsCmd.AddCommand(groupsListCmd)

	// Create command flags
	groupsCreateCmd.Flags().StringP("name", "n", "", "Group name (required)")
	groupsCreateCmd.Flags().StringSliceP("policies", "p", []string{}, "Policy ARNs to attach (required)")
	groupsCreateCmd.Flags().StringP("description", "d", "", "Group description")

	groupsCreateCmd.MarkFlagRequired("name")
	groupsCreateCmd.MarkFlagRequired("policies")

	// Remove command flags
	groupsRemoveCmd.Flags().StringP("name", "n", "", "Group name (required)")
	groupsRemoveCmd.MarkFlagRequired("name")
}