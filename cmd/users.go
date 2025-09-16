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
		idleThreshold, _ := cmd.Flags().GetInt("idle-threshold")
		idleDuration, _ := cmd.Flags().GetInt("idle-duration")

		return createUsersWithIdle(cmd.Context(), project, blueprint, bundle, region, users, idleThreshold, idleDuration)
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

		return removeUsers(cmd.Context(), project, users, all)
	},
}

var usersListCmd = &cobra.Command{
	Use:   "list",
	Short: "List IAM users and their instances",
	Long:  `List all IAM users and their associated Lightsail instances, optionally filtered by project.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		project, _ := cmd.Flags().GetString("project")

		return listUsers(cmd.Context(), project)
	},
}

var usersCreateBulkCmd = &cobra.Command{
	Use:   "create-bulk [csv-file]",
	Short: "Create multiple users from CSV file",
	Long: `Create multiple IAM users and instances from a CSV file. The CSV should contain
columns: username, project, blueprint, bundle, groups (optional).`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		csvFile := args[0]
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		continueOnError, _ := cmd.Flags().GetBool("continue-on-error")
		startStopped, _ := cmd.Flags().GetBool("start-stopped")

		return createBulkUsers(cmd.Context(), csvFile, dryRun, continueOnError, startStopped)
	},
}

var usersRemoveBulkCmd = &cobra.Command{
	Use:   "remove-bulk",
	Short: "Remove multiple users with progress tracking",
	Long: `Remove multiple users and their instances with detailed progress tracking
and optional rollback capabilities.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		users, _ := cmd.Flags().GetStringSlice("users")
		project, _ := cmd.Flags().GetString("project")
		csvFile, _ := cmd.Flags().GetString("csv")
		confirm, _ := cmd.Flags().GetBool("confirm")

		return removeBulkUsers(cmd.Context(), users, project, csvFile, confirm)
	},
}

var usersTemplateCmd = &cobra.Command{
	Use:   "template [filename]",
	Short: "Generate CSV template for bulk operations",
	Long: `Generate a CSV template file with sample data for bulk user creation.
Use this as a starting point for your bulk user operations.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filename := args[0]

		return generateUserTemplate(filename)
	},
}

func init() {
	rootCmd.AddCommand(usersCmd)

	usersCmd.AddCommand(usersCreateCmd)
	usersCmd.AddCommand(usersRemoveCmd)
	usersCmd.AddCommand(usersListCmd)
	usersCmd.AddCommand(usersCreateBulkCmd)
	usersCmd.AddCommand(usersRemoveBulkCmd)
	usersCmd.AddCommand(usersTemplateCmd)

	// Create command flags
	usersCreateCmd.Flags().StringP("project", "p", "", "Project name (required)")
	usersCreateCmd.Flags().StringP("blueprint", "b", "", "Lightsail blueprint ID (required)")
	usersCreateCmd.Flags().String("bundle", "", "Lightsail bundle ID (required)")
	usersCreateCmd.Flags().StringP("region", "r", "", "AWS region (required)")
	usersCreateCmd.Flags().StringSliceP("users", "u", []string{}, "Comma-separated list of usernames (required)")
	usersCreateCmd.Flags().IntP("idle-threshold", "", 120, "Idle threshold in minutes (default: 120)")
	usersCreateCmd.Flags().IntP("idle-duration", "", 30, "Duration in minutes before stopping (default: 30)")

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

	// Bulk create command flags
	usersCreateBulkCmd.Flags().BoolP("dry-run", "d", false, "Show what would be created without executing")
	usersCreateBulkCmd.Flags().BoolP("continue-on-error", "c", false, "Continue creating users even if some fail")
	usersCreateBulkCmd.Flags().BoolP("start-stopped", "s", false, "Create instances but immediately stop them to save costs")
	usersCreateBulkCmd.Flags().String("from-snapshot", "", "Create instances from existing snapshot instead of blueprint")

	// Bulk remove command flags
	usersRemoveBulkCmd.Flags().StringSliceP("users", "u", []string{}, "Specific usernames to remove")
	usersRemoveBulkCmd.Flags().StringP("project", "p", "", "Remove all users from project")
	usersRemoveBulkCmd.Flags().StringP("csv", "", "", "CSV file containing users to remove")
	usersRemoveBulkCmd.Flags().BoolP("confirm", "y", false, "Skip confirmation prompts")
}

// createUsers implements the core user creation logic from the original script.
func createUsers(ctx context.Context, project, blueprint, bundle, region string, usernames []string) error {
	fmt.Printf("Creating %d users for project: %s\n", len(usernames), project)
	fmt.Printf("Blueprint: %s, Bundle: %s, Region: %s\n", blueprint, bundle, region)

	// Load configuration
	_, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Create AWS client
	awsClient, err := aws.NewClient(ctx, aws.Options{
		Region:  region,
		Profile: viper.GetString("aws.profile"),
	})
	if err != nil {
		return fmt.Errorf("failed to create AWS client: %w", err)
	}

	iamService := aws.NewIAMService(awsClient)
	lightsailService := aws.NewLightsailService(awsClient)

	// Step 1: Ensure LightsailReadOnly policy exists
	lightsailPolicyDoc := `{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Action": [
					"lightsail:Get*"
				],
				"Resource": "*"
			}
		]
	}`

	policyARN, err := iamService.CreatePolicy(ctx, "LightsailReadOnly", "Read-only access to the Lightsail service", lightsailPolicyDoc)
	if err != nil {
		return fmt.Errorf("failed to create or get LightsailReadOnly policy: %w", err)
	}

	// Step 2: Ensure Lightsail-Users group exists
	changePasswordPolicyARN := "arn:aws:iam::aws:policy/IAMUserChangePassword"
	_, err = iamService.CreateGroup(ctx, "Lightsail-Users", "Group for Lightsail for Research users", []string{policyARN, changePasswordPolicyARN})
	if err != nil {
		return fmt.Errorf("failed to create or get Lightsail-Users group: %w", err)
	}

	// Step 3: Create users and instances
	availabilityZone := region + "a"

	fmt.Printf("\nCreating %d users and instances...\n", len(usernames))

	for i, username := range usernames {
		fmt.Printf("\n[%d/%d] Creating user: %s\n", i+1, len(usernames), username)

		// Generate secure password
		password, err := utils.GeneratePassword()
		if err != nil {
			fmt.Printf("❌ Error generating password for %s: %v\n", username, err)
			continue
		}

		// Create IAM user
		_, err = iamService.CreateUser(ctx, username, password, project)
		if err != nil {
			fmt.Printf("❌ Error creating user %s: %v\n", username, err)
			continue
		}

		// Add user to Lightsail-Users group
		err = iamService.AddUserToGroup(ctx, username, "Lightsail-Users")
		if err != nil {
			fmt.Printf("❌ Error adding user %s to group: %v\n", username, err)
			continue
		}

		// Create Lightsail instance
		instanceName := username + "-" + blueprint
		instance, err := lightsailService.CreateInstance(ctx, instanceName, blueprint, bundle, availabilityZone, project)
		if err != nil {
			fmt.Printf("❌ Error creating instance for %s: %v\n", username, err)
			continue
		}

		// Create user-specific policy for their instance
		userPolicyDoc := fmt.Sprintf(`{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Effect": "Allow",
					"Action": [
						"lightsail:*"
					],
					"Resource": "%s"
				}
			]
		}`, instance.ARN)

		policyName := fmt.Sprintf("LightsailLimitedAccess-%s", username)
		err = iamService.PutUserPolicy(ctx, username, policyName, userPolicyDoc)
		if err != nil {
			fmt.Printf("❌ Error creating user policy for %s: %v\n", username, err)
			continue
		}

		fmt.Printf("✅ %s : %s : %s\n", username, password, instance.ARN)
	}

	fmt.Printf("\n🎉 User creation completed!\n")
	return nil
}

// removeUsers removes IAM users and their Lightsail instances.
func removeUsers(ctx context.Context, project string, usernames []string, all bool) error {
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
	lightsailService := aws.NewLightsailService(awsClient)

	// Get list of users to remove
	var usersToRemove []string
	if all {
		// Get all instances for the project and extract usernames
		instances, err := lightsailService.ListInstances(ctx, project)
		if err != nil {
			return fmt.Errorf("failed to list instances for project %s: %w", project, err)
		}

		for _, instance := range instances {
			// Extract username from instance name (format: username-blueprint)
			parts := strings.Split(instance.Name, "-")
			if len(parts) > 0 {
				usersToRemove = append(usersToRemove, parts[0])
			}
		}
	} else {
		usersToRemove = usernames
	}

	if len(usersToRemove) == 0 {
		fmt.Println("No users to remove.")
		return nil
	}

	fmt.Printf("Removing %d users from project: %s\n", len(usersToRemove), project)
	fmt.Printf("⚠️  This will delete all user data and instances!\n\n")

	for i, username := range usersToRemove {
		fmt.Printf("[%d/%d] Removing user: %s\n", i+1, len(usersToRemove), username)

		// Get user's instances to delete
		instances, err := lightsailService.ListInstances(ctx, project)
		if err != nil {
			fmt.Printf("❌ Error listing instances for %s: %v\n", username, err)
			continue
		}

		// Delete instances belonging to this user
		for _, instance := range instances {
			if strings.HasPrefix(instance.Name, username+"-") {
				fmt.Printf("  Deleting instance: %s\n", instance.Name)
				err = lightsailService.DeleteInstance(ctx, instance.Name)
				if err != nil {
					fmt.Printf("❌ Error deleting instance %s: %v\n", instance.Name, err)
				}
			}
		}

		// Delete IAM user (this handles group removal, policies, login profile)
		err = iamService.DeleteUser(ctx, username)
		if err != nil {
			fmt.Printf("❌ Error deleting user %s: %v\n", username, err)
			continue
		}

		fmt.Printf("✅ Removed user: %s\n", username)
	}

	fmt.Printf("\n🎉 User removal completed!\n")
	return nil
}

// listUsers lists IAM users and their instances.
func listUsers(ctx context.Context, project string) error {
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

	// Display results
	if project != "" {
		fmt.Printf("Users and instances for project: %s\n\n", project)
	} else {
		fmt.Printf("All users and instances:\n\n")
	}

	fmt.Printf("%-15s %-20s %-12s %-15s %-20s %-18s\n",
		"USERNAME", "INSTANCE", "STATE", "BLUEPRINT", "BUNDLE", "PUBLIC IP")
	fmt.Println(strings.Repeat("-", 120))

	for _, instance := range instances {
		// Extract username from instance name
		username := "unknown"
		parts := strings.Split(instance.Name, "-")
		if len(parts) > 0 {
			username = parts[0]
		}

		publicIP := instance.PublicIP
		if publicIP == "" {
			publicIP = "-"
		}

		fmt.Printf("%-15s %-20s %-12s %-15s %-20s %-18s\n",
			username,
			instance.Name,
			instance.State,
			instance.Blueprint,
			instance.Bundle,
			publicIP,
		)
	}

	fmt.Printf("\nTotal: %d instances\n", len(instances))
	return nil
}

// createBulkUsers creates multiple users from a CSV file.
func createBulkUsers(ctx context.Context, csvFile string, dryRun, continueOnError, startStopped bool) error {
	// Parse CSV file
	users, err := utils.ParseUsersCSV(csvFile)
	if err != nil {
		return fmt.Errorf("failed to parse CSV file: %w", err)
	}

	if len(users) == 0 {
		return fmt.Errorf("no users found in CSV file")
	}

	fmt.Printf("Bulk user creation from: %s\n", csvFile)
	fmt.Printf("Total users: %d\n\n", len(users))

	if dryRun {
		fmt.Printf("DRY RUN - Users that would be created:\n\n")
		fmt.Printf("%-15s %-15s %-15s %-20s %-15s\n",
			"USERNAME", "PROJECT", "BLUEPRINT", "BUNDLE", "GROUPS")
		fmt.Println(strings.Repeat("-", 95))

		for _, user := range users {
			groupsStr := strings.Join(user.Groups, ",")
			if groupsStr == "" {
				groupsStr = "-"
			}
			fmt.Printf("%-15s %-15s %-15s %-20s %-15s\n",
				user.Username, user.Project, user.Blueprint, user.Bundle, groupsStr)
		}
		fmt.Printf("\nTotal: %d users would be created\n", len(users))
		return nil
	}

	// Group users by project and configuration for efficient creation
	type ProjectConfig struct {
		Project   string
		Blueprint string
		Bundle    string
		Region    string
		Users     []string
	}

	projectGroups := make(map[string]*ProjectConfig)
	for _, user := range users {
		key := fmt.Sprintf("%s-%s-%s", user.Project, user.Blueprint, user.Bundle)
		if _, exists := projectGroups[key]; !exists {
			projectGroups[key] = &ProjectConfig{
				Project:   user.Project,
				Blueprint: user.Blueprint,
				Bundle:    user.Bundle,
				Region:    viper.GetString("aws.region"),
				Users:     []string{},
			}
		}
		projectGroups[key].Users = append(projectGroups[key].Users, user.Username)
	}

	fmt.Printf("Creating users in %d batch(es):\n\n", len(projectGroups))

	successCount := 0
	failCount := 0
	batchNum := 1

	for _, config := range projectGroups {
		fmt.Printf("Batch %d: %s/%s/%s (%d users)\n",
			batchNum, config.Project, config.Blueprint, config.Bundle, len(config.Users))
		batchNum++

		err := createUsers(ctx, config.Project, config.Blueprint, config.Bundle, config.Region, config.Users)
		if err != nil {
			fmt.Printf("❌ Batch failed: %v\n", err)
			failCount += len(config.Users)
			if !continueOnError {
				return fmt.Errorf("bulk creation failed, use --continue-on-error to proceed despite failures")
			}
		} else {
			successCount += len(config.Users)

			// Stop instances immediately if requested
			if startStopped {
				fmt.Printf("🛑 Stopping instances to save costs...\n")
				stopErr := stopInstances(ctx, config.Users, config.Project, false)
				if stopErr != nil {
					fmt.Printf("⚠️  Warning: Failed to stop some instances: %v\n", stopErr)
				} else {
					fmt.Printf("✅ Instances stopped for cost savings\n")
				}
			}
		}
		fmt.Println()
	}

	fmt.Printf("🎉 Bulk user creation completed!\n")
	fmt.Printf("✅ Success: %d users\n", successCount)
	if failCount > 0 {
		fmt.Printf("❌ Failed: %d users\n", failCount)
	}

	return nil
}

// removeBulkUsers removes multiple users with progress tracking.
func removeBulkUsers(ctx context.Context, users []string, project, csvFile string, confirm bool) error {
	var usersToRemove []string

	// Determine users to remove
	if csvFile != "" {
		bulkUsers, err := utils.ParseUsersCSV(csvFile)
		if err != nil {
			return fmt.Errorf("failed to parse CSV file: %w", err)
		}
		for _, user := range bulkUsers {
			usersToRemove = append(usersToRemove, user.Username)
		}
	} else if project != "" {
		// Get all users from project (via instances)
		return removeUsers(ctx, project, []string{}, true)
	} else if len(users) > 0 {
		usersToRemove = users
	} else {
		return fmt.Errorf("must specify --users, --project, or --csv")
	}

	if len(usersToRemove) == 0 {
		fmt.Println("No users to remove.")
		return nil
	}

	fmt.Printf("Bulk user removal:\n")
	fmt.Printf("Users to remove: %v\n", usersToRemove)
	fmt.Printf("Total: %d users\n\n", len(usersToRemove))

	if !confirm {
		fmt.Printf("⚠️  This will permanently delete:\n")
		fmt.Printf("   - %d IAM users and their login profiles\n", len(usersToRemove))
		fmt.Printf("   - All associated Lightsail instances\n")
		fmt.Printf("   - All user data and configurations\n\n")
		fmt.Printf("Run with --confirm to proceed.\n")
		return nil
	}

	fmt.Printf("Removing %d users...\n\n", len(usersToRemove))

	successCount := 0
	failCount := 0

	for i, username := range usersToRemove {
		fmt.Printf("[%d/%d] Removing user: %s\n", i+1, len(usersToRemove), username)

		err := removeUsers(ctx, "", []string{username}, false)
		if err != nil {
			fmt.Printf("❌ Failed to remove %s: %v\n", username, err)
			failCount++
		} else {
			successCount++
		}
	}

	fmt.Printf("\n🎉 Bulk user removal completed!\n")
	fmt.Printf("✅ Success: %d users\n", successCount)
	if failCount > 0 {
		fmt.Printf("❌ Failed: %d users\n", failCount)
	}

	return nil
}

// generateUserTemplate generates a CSV template for bulk user creation.
func generateUserTemplate(filename string) error {
	err := utils.GenerateUsersCSVTemplate(filename)
	if err != nil {
		return fmt.Errorf("failed to generate template: %w", err)
	}

	fmt.Printf("✅ User CSV template created: %s\n", filename)
	fmt.Printf("\nTemplate contains sample data for:\n")
	fmt.Printf("- Research team users (alice, bob)\n")
	fmt.Printf("- ML team with GPU instance (charlie)\n")
	fmt.Printf("- Custom groups and bundle configurations\n\n")
	fmt.Printf("Edit the file and run: lfr users create-bulk %s\n", filename)

	return nil
}

// createUsersWithIdle implements user creation with custom idle detection settings.
func createUsersWithIdle(ctx context.Context, project, blueprint, bundle, region string, usernames []string, idleThreshold, idleDuration int) error {
	fmt.Printf("Creating %d users for project: %s\n", len(usernames), project)
	fmt.Printf("Blueprint: %s, Bundle: %s, Region: %s\n", blueprint, bundle, region)
	fmt.Printf("Idle detection: %d minutes threshold, %d minutes duration\n", idleThreshold, idleDuration)

	// For now, use the standard createUsers function
	// TODO: Implement CreateInstanceWithIdleDetection in lightsail service
	fmt.Printf("Note: Custom idle detection requires API enhancement\n")
	fmt.Printf("Using standard idle detection (120min/30min) for now\n\n")

	return createUsers(ctx, project, blueprint, bundle, region, usernames)
}
