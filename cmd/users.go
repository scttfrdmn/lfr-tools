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

		return createUsers(cmd.Context(), project, blueprint, bundle, region, users)
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

	fmt.Printf("%-20s %-20s %-15s %-15s %-15s %-15s\n",
		"USERNAME", "INSTANCE", "STATE", "BLUEPRINT", "BUNDLE", "PUBLIC IP")
	fmt.Println(strings.Repeat("-", 110))

	for _, instance := range instances {
		// Extract username from instance name
		username := "unknown"
		parts := strings.Split(instance.Name, "-")
		if len(parts) > 0 {
			username = parts[0]
		}

		fmt.Printf("%-20s %-20s %-15s %-15s %-15s %-15s\n",
			username,
			instance.Name,
			instance.State,
			instance.Blueprint,
			instance.Bundle,
			instance.PublicIP,
		)
	}

	fmt.Printf("\nTotal: %d instances\n", len(instances))
	return nil
}