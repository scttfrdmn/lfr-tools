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

		return createGroup(cmd.Context(), name, description, policies)
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

// createGroup creates an IAM group with policies.
func createGroup(ctx context.Context, name, description string, policies []string) error {
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

	iamService := aws.NewIAMService(awsClient)

	fmt.Printf("Creating group: %s\n", name)
	if description != "" {
		fmt.Printf("Description: %s\n", description)
	}
	fmt.Printf("Policies: %v\n", policies)

	// Create group with policies
	group, err := iamService.CreateGroup(ctx, name, description, policies)
	if err != nil {
		return fmt.Errorf("failed to create group: %w", err)
	}

	fmt.Printf("✅ Group created: %s\n", group.Name)
	fmt.Printf("Attached policies: %v\n", group.Policies)

	return nil
}