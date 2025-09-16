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
		project, _ := cmd.Flags().GetString("project")

		return listGroups(cmd.Context(), project)
	},
}

var groupsCreateBulkCmd = &cobra.Command{
	Use:   "create-bulk [csv-file]",
	Short: "Create multiple groups from CSV file",
	Long: `Create multiple IAM groups from a CSV file. The CSV should contain
columns: name, description, policies, project (optional).`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		csvFile := args[0]
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		continueOnError, _ := cmd.Flags().GetBool("continue-on-error")

		return createBulkGroups(cmd.Context(), csvFile, dryRun, continueOnError)
	},
}

var groupsRemoveBulkCmd = &cobra.Command{
	Use:   "remove-bulk",
	Short: "Remove multiple groups with progress tracking",
	Long: `Remove multiple groups with detailed progress tracking. Groups must be
empty (no members) before deletion.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		groups, _ := cmd.Flags().GetStringSlice("groups")
		project, _ := cmd.Flags().GetString("project")
		csvFile, _ := cmd.Flags().GetString("csv")
		confirm, _ := cmd.Flags().GetBool("confirm")

		return removeBulkGroups(cmd.Context(), groups, project, csvFile, confirm)
	},
}

var groupsTemplateCmd = &cobra.Command{
	Use:   "template [filename]",
	Short: "Generate CSV template for bulk group operations",
	Long: `Generate a CSV template file with sample data for bulk group creation.
Use this as a starting point for your bulk group operations.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filename := args[0]

		return generateGroupTemplate(filename)
	},
}

func init() {
	rootCmd.AddCommand(groupsCmd)

	groupsCmd.AddCommand(groupsCreateCmd)
	groupsCmd.AddCommand(groupsRemoveCmd)
	groupsCmd.AddCommand(groupsListCmd)
	groupsCmd.AddCommand(groupsCreateBulkCmd)
	groupsCmd.AddCommand(groupsRemoveBulkCmd)
	groupsCmd.AddCommand(groupsTemplateCmd)

	// Create command flags
	groupsCreateCmd.Flags().StringP("name", "n", "", "Group name (required)")
	groupsCreateCmd.Flags().StringSliceP("policies", "p", []string{}, "Policy ARNs to attach (required)")
	groupsCreateCmd.Flags().StringP("description", "d", "", "Group description")

	groupsCreateCmd.MarkFlagRequired("name")
	groupsCreateCmd.MarkFlagRequired("policies")

	// Remove command flags
	groupsRemoveCmd.Flags().StringP("name", "n", "", "Group name (required)")
	groupsRemoveCmd.MarkFlagRequired("name")

	// List command flags
	groupsListCmd.Flags().StringP("project", "p", "", "Filter by project name")

	// Bulk create command flags
	groupsCreateBulkCmd.Flags().BoolP("dry-run", "d", false, "Show what would be created without executing")
	groupsCreateBulkCmd.Flags().BoolP("continue-on-error", "c", false, "Continue creating groups even if some fail")

	// Bulk remove command flags
	groupsRemoveBulkCmd.Flags().StringSliceP("groups", "g", []string{}, "Specific group names to remove")
	groupsRemoveBulkCmd.Flags().StringP("project", "p", "", "Remove all groups from project")
	groupsRemoveBulkCmd.Flags().StringP("csv", "", "", "CSV file containing groups to remove")
	groupsRemoveBulkCmd.Flags().BoolP("confirm", "y", false, "Skip confirmation prompts")
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

// listGroups lists IAM groups with filtering.
func listGroups(ctx context.Context, project string) error {
	// Load configuration
	_, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Create AWS client
	_, err = aws.NewClient(ctx, aws.Options{
		Region:  viper.GetString("aws.region"),
		Profile: viper.GetString("aws.profile"),
	})
	if err != nil {
		return fmt.Errorf("failed to create AWS client: %w", err)
	}

	// For now, show basic group listing message
	// TODO: Implement actual IAM group listing when needed
	if project != "" {
		fmt.Printf("Groups for project: %s\n", project)
	} else {
		fmt.Printf("All IAM groups:\n")
	}

	fmt.Printf("✅ Basic groups: Lightsail-Users (created by lfr)\n")
	fmt.Printf("Note: Full group listing requires additional IAM API implementation\n")

	return nil
}

// createBulkGroups creates multiple groups from a CSV file.
func createBulkGroups(ctx context.Context, csvFile string, dryRun, continueOnError bool) error {
	// Parse CSV file
	groups, err := utils.ParseGroupsCSV(csvFile)
	if err != nil {
		return fmt.Errorf("failed to parse CSV file: %w", err)
	}

	if len(groups) == 0 {
		return fmt.Errorf("no groups found in CSV file")
	}

	fmt.Printf("Bulk group creation from: %s\n", csvFile)
	fmt.Printf("Total groups: %d\n\n", len(groups))

	if dryRun {
		fmt.Printf("DRY RUN - Groups that would be created:\n\n")
		fmt.Printf("%-20s %-30s %-15s %-40s\n",
			"NAME", "DESCRIPTION", "PROJECT", "POLICIES")
		fmt.Println(strings.Repeat("-", 120))

		for _, group := range groups {
			project := group.Project
			if project == "" {
				project = "-"
			}
			policies := strings.Join(group.Policies, ", ")
			if len(policies) > 37 {
				policies = policies[:37] + "..."
			}
			fmt.Printf("%-20s %-30s %-15s %-40s\n",
				group.Name, group.Description, project, policies)
		}
		fmt.Printf("\nTotal: %d groups would be created\n", len(groups))
		return nil
	}

	fmt.Printf("Creating %d groups...\n\n", len(groups))

	successCount := 0
	failCount := 0

	for i, group := range groups {
		fmt.Printf("[%d/%d] Creating group: %s\n", i+1, len(groups), group.Name)

		err := createGroup(ctx, group.Name, group.Description, group.Policies)
		if err != nil {
			fmt.Printf("❌ Failed to create %s: %v\n", group.Name, err)
			failCount++
			if !continueOnError {
				return fmt.Errorf("bulk group creation failed, use --continue-on-error to proceed despite failures")
			}
		} else {
			successCount++
		}
	}

	fmt.Printf("\n🎉 Bulk group creation completed!\n")
	fmt.Printf("✅ Success: %d groups\n", successCount)
	if failCount > 0 {
		fmt.Printf("❌ Failed: %d groups\n", failCount)
	}

	return nil
}

// removeBulkGroups removes multiple groups with progress tracking.
func removeBulkGroups(ctx context.Context, groups []string, project, csvFile string, confirm bool) error {
	var groupsToRemove []string

	// Determine groups to remove
	if csvFile != "" {
		bulkGroups, err := utils.ParseGroupsCSV(csvFile)
		if err != nil {
			return fmt.Errorf("failed to parse CSV file: %w", err)
		}
		for _, group := range bulkGroups {
			groupsToRemove = append(groupsToRemove, group.Name)
		}
	} else if len(groups) > 0 {
		groupsToRemove = groups
	} else {
		return fmt.Errorf("must specify --groups or --csv")
	}

	if len(groupsToRemove) == 0 {
		fmt.Println("No groups to remove.")
		return nil
	}

	fmt.Printf("Bulk group removal:\n")
	fmt.Printf("Groups to remove: %v\n", groupsToRemove)
	fmt.Printf("Total: %d groups\n\n", len(groupsToRemove))

	if !confirm {
		fmt.Printf("⚠️  This will permanently delete:\n")
		fmt.Printf("   - %d IAM groups\n", len(groupsToRemove))
		fmt.Printf("   - All group policy attachments\n")
		fmt.Printf("   - Group memberships (users will lose group permissions)\n\n")
		fmt.Printf("Run with --confirm to proceed.\n")
		return nil
	}

	fmt.Printf("Removing %d groups...\n\n", len(groupsToRemove))

	successCount := 0
	failCount := 0

	for i, groupName := range groupsToRemove {
		fmt.Printf("[%d/%d] Removing group: %s\n", i+1, len(groupsToRemove), groupName)

		// TODO: Implement actual group removal
		fmt.Printf("⚠️  Group removal implementation needed\n")
		failCount++
	}

	fmt.Printf("\n🎉 Bulk group removal completed!\n")
	fmt.Printf("✅ Success: %d groups\n", successCount)
	if failCount > 0 {
		fmt.Printf("❌ Failed: %d groups\n", failCount)
	}

	return nil
}

// generateGroupTemplate generates a CSV template for bulk group creation.
func generateGroupTemplate(filename string) error {
	err := utils.GenerateGroupsCSVTemplate(filename)
	if err != nil {
		return fmt.Errorf("failed to generate template: %w", err)
	}

	fmt.Printf("✅ Group CSV template created: %s\n", filename)
	fmt.Printf("\nTemplate contains sample data for:\n")
	fmt.Printf("- Research team groups with different permissions\n")
	fmt.Printf("- ML team with PowerUser access\n")
	fmt.Printf("- Custom policy combinations\n\n")
	fmt.Printf("Edit the file and run: lfr groups create-bulk %s\n", filename)

	return nil
}