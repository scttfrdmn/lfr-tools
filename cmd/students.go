package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/scttfrdmn/lfr-tools/internal/aws"
	"github.com/scttfrdmn/lfr-tools/internal/config"
	"github.com/scttfrdmn/lfr-tools/internal/utils"
)

var studentsCmd = &cobra.Command{
	Use:   "students",
	Short: "Manage educational access for students and TAs",
	Long:  `Set up and manage educational access systems with hardware-tied tokens and automatic status updates.`,
}

var studentsSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Set up educational access system",
	Long:  `Set up and manage educational access systems with hardware-tied tokens and automatic status updates.`,
}

var studentsSetupClassCmd = &cobra.Command{
	Use:   "environment",
	Short: "Set up educational environment with student access system",
	Long: `Set up a complete educational environment (class, lab, or project) with S3 status bucket,
student tokens, and automatic status synchronization.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		project, _ := cmd.Flags().GetString("project")
		bucket, _ := cmd.Flags().GetString("s3-bucket")
		students, _ := cmd.Flags().GetStringSlice("students")
		tas, _ := cmd.Flags().GetStringSlice("tas")
		professor, _ := cmd.Flags().GetString("professor")
		startDate, _ := cmd.Flags().GetString("start-date")
		endDate, _ := cmd.Flags().GetString("end-date")

		return setupClass(cmd.Context(), project, bucket, students, tas, professor, startDate, endDate)
	},
}

var studentsGenerateCmd = &cobra.Command{
	Use:   "generate tokens",
	Short: "Generate access tokens for students",
	Long: `Generate hardware-bindable access tokens for students and TAs.
Tokens are distributed to users for one-time activation.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		project, _ := cmd.Flags().GetString("project")
		outputDir, _ := cmd.Flags().GetString("output")

		return generateStudentTokens(cmd.Context(), project, outputDir)
	},
}

var studentsCheckCmd = &cobra.Command{
	Use:   "check requests",
	Short: "Check pending start requests from students",
	Long: `Check S3 for pending instance start requests from students and TAs.
Shows who is requesting access and when.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		project, _ := cmd.Flags().GetString("project")
		autoApprove, _ := cmd.Flags().GetBool("auto-approve")

		return checkStartRequests(cmd.Context(), project, autoApprove)
	},
}

var studentsStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show comprehensive student access status",
	Long: `Display status of all students including instance states, budget usage,
access permissions, and recent activity.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		project, _ := cmd.Flags().GetString("project")

		return showStudentStatus(cmd.Context(), project)
	},
}

func init() {
	rootCmd.AddCommand(studentsCmd)

	studentsCmd.AddCommand(studentsSetupCmd)
	studentsCmd.AddCommand(studentsGenerateCmd)
	studentsCmd.AddCommand(studentsCheckCmd)
	studentsCmd.AddCommand(studentsStatusCmd)

	// Add setup subcommands
	studentsSetupCmd.AddCommand(studentsSetupClassCmd)

	// Setup class command flags
	studentsSetupClassCmd.Flags().StringP("project", "p", "", "Project/class name (required)")
	studentsSetupClassCmd.Flags().String("s3-bucket", "", "S3 bucket for status updates (required)")
	studentsSetupClassCmd.Flags().StringSliceP("students", "s", []string{}, "Student usernames")
	studentsSetupClassCmd.Flags().StringSlice("tas", []string{}, "TA usernames")
	studentsSetupClassCmd.Flags().String("professor", "", "Professor username")
	studentsSetupClassCmd.Flags().String("start-date", "", "Course start date (YYYY-MM-DD)")
	studentsSetupClassCmd.Flags().String("end-date", "", "Course end date (YYYY-MM-DD)")
	studentsSetupClassCmd.MarkFlagRequired("project")
	studentsSetupClassCmd.MarkFlagRequired("s3-bucket")

	// Generate command flags
	studentsGenerateCmd.Flags().StringP("project", "p", "", "Project name (required)")
	studentsGenerateCmd.Flags().StringP("output", "o", "./student-tokens", "Output directory for tokens")
	studentsGenerateCmd.MarkFlagRequired("project")

	// Check command flags
	studentsCheckCmd.Flags().StringP("project", "p", "", "Project name (required)")
	studentsCheckCmd.Flags().BoolP("auto-approve", "a", false, "Automatically approve start requests")
	studentsCheckCmd.MarkFlagRequired("project")

	// Status command flags
	studentsStatusCmd.Flags().StringP("project", "p", "", "Project name (required)")
	studentsStatusCmd.MarkFlagRequired("project")
}

// setupClass sets up a complete class environment.
func setupClass(ctx context.Context, project, bucket string, students, tas []string, professor, startDate, endDate string) error {
	fmt.Printf("Setting up class environment for project: %s\n", project)
	fmt.Printf("S3 bucket: %s\n", bucket)
	fmt.Printf("Students: %d, TAs: %d, Professor: %s\n", len(students), len(tas), professor)

	// Parse dates
	var start, end time.Time
	var err error
	if startDate != "" {
		start, err = time.Parse("2006-01-02", startDate)
		if err != nil {
			return fmt.Errorf("invalid start date format: %s (use YYYY-MM-DD)", startDate)
		}
	}
	if endDate != "" {
		end, err = time.Parse("2006-01-02", endDate)
		if err != nil {
			return fmt.Errorf("invalid end date format: %s (use YYYY-MM-DD)", endDate)
		}
	}

	// Create AWS client
	awsClient, err := aws.NewClient(ctx, aws.Options{
		Region:  viper.GetString("aws.region"),
		Profile: viper.GetString("aws.profile"),
	})
	if err != nil {
		return fmt.Errorf("failed to create AWS client: %w", err)
	}

	s3Service := aws.NewS3Service(awsClient)

	// Create S3 bucket for status updates
	fmt.Printf("Creating S3 bucket for student status updates...\n")
	err = s3Service.CreateStatusBucket(ctx, bucket, project)
	if err != nil {
		fmt.Printf("⚠️ Warning: S3 bucket creation failed (may already exist): %v\n", err)
	} else {
		fmt.Printf("✅ S3 bucket created: %s\n", bucket)
	}

	// Enable S3 sync for this project
	err = utils.EnableS3Sync(project, bucket)
	if err != nil {
		return fmt.Errorf("failed to enable S3 sync: %w", err)
	}

	fmt.Printf("✅ S3 sync enabled for project %s\n", project)

	// Store class configuration
	classConfig := map[string]interface{}{
		"project":    project,
		"bucket":     bucket,
		"students":   students,
		"tas":        tas,
		"professor":  professor,
		"start_date": start,
		"end_date":   end,
		"created_at": time.Now(),
	}

	configFile := fmt.Sprintf(".lfr-class-%s.json", project)
	configData, _ := json.MarshalIndent(classConfig, "", "  ")
	_ = os.WriteFile(configFile, configData, 0644)

	fmt.Printf("✅ Class setup completed!\n")
	fmt.Printf("\nNext steps:\n")
	fmt.Printf("1. Generate tokens: lfr students generate tokens --project=%s\n", project)
	fmt.Printf("2. Create users: lfr users create-bulk students.csv\n")
	fmt.Printf("3. Distribute tokens to students\n")

	return nil
}

// generateStudentTokens generates access tokens for all students in a project.
func generateStudentTokens(ctx context.Context, project, outputDir string) error {
	fmt.Printf("Generating access tokens for project: %s\n", project)

	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Load class configuration
	classConfigFile := fmt.Sprintf(".lfr-class-%s.json", project)
	configData, err := os.ReadFile(classConfigFile)
	if err != nil {
		return fmt.Errorf("class not found. Run: lfr students setup class --project=%s", project)
	}

	var classConfig map[string]interface{}
	if err := json.Unmarshal(configData, &classConfig); err != nil {
		return fmt.Errorf("failed to parse class config: %w", err)
	}

	bucket := classConfig["bucket"].(string)
	students := classConfig["students"].([]interface{})
	tas := classConfig["tas"].([]interface{})

	// Generate tokens for students and TAs
	tm, err := config.NewTokenManager()
	if err != nil {
		return fmt.Errorf("failed to initialize token manager: %w", err)
	}

	tokensFile := filepath.Join(outputDir, fmt.Sprintf("%s-tokens.txt", project))
	tokensList, err := os.Create(tokensFile)
	if err != nil {
		return fmt.Errorf("failed to create tokens file: %w", err)
	}
	defer tokensList.Close()

	fmt.Fprintf(tokensList, "# Access tokens for %s\n", project)
	fmt.Fprintf(tokensList, "# Format: USERNAME:ROLE:TOKEN\n")
	fmt.Fprintf(tokensList, "# Distribution: Send each user their specific token\n\n")

	// Generate student tokens
	fmt.Printf("Generating tokens for %d students...\n", len(students))
	for i, studentInterface := range students {
		student := studentInterface.(string)

		tokenString, _, err := tm.GenerateToken(
			project, student, fmt.Sprintf("student-%d", i+1), "student",
			[]string{"connect"}, bucket, time.Now().AddDate(0, 6, 0)) // 6 months
		if err != nil {
			fmt.Printf("❌ Failed to generate token for %s: %v\n", student, err)
			continue
		}

		fmt.Fprintf(tokensList, "%s:student:%s\n", student, tokenString)
		fmt.Printf("✅ Generated token for student: %s\n", student)
	}

	// Generate TA tokens
	if len(tas) > 0 {
		fmt.Printf("Generating tokens for %d TAs...\n", len(tas))
		for i, taInterface := range tas {
			ta := taInterface.(string)

			tokenString, _, err := tm.GenerateToken(
				project, ta, fmt.Sprintf("ta-%d", i+1), "ta",
				[]string{"connect", "start", "stop", "status"}, bucket, time.Now().AddDate(0, 6, 0))
			if err != nil {
				fmt.Printf("❌ Failed to generate token for TA %s: %v\n", ta, err)
				continue
			}

			fmt.Fprintf(tokensList, "%s:ta:%s\n", ta, tokenString)
			fmt.Printf("✅ Generated token for TA: %s\n", ta)
		}
	}

	fmt.Printf("\n🎉 Token generation completed!\n")
	fmt.Printf("Tokens saved to: %s\n", tokensFile)
	fmt.Printf("\nDistribution instructions:\n")
	fmt.Printf("1. Send each user their specific token via secure email\n")
	fmt.Printf("2. Include activation instructions:\n")
	fmt.Printf("   brew install lfr\n")
	fmt.Printf("   lfr connect activate <their-token> <their-student-id>\n")
	fmt.Printf("   lfr connect <their-username>\n")

	return nil
}

// checkStartRequests checks for pending start requests.
func checkStartRequests(ctx context.Context, project string, autoApprove bool) error {
	// Get S3 sync config
	syncConfig := utils.GetS3SyncConfig()
	if !syncConfig.Enabled {
		return fmt.Errorf("S3 sync not enabled. Run: lfr students setup class")
	}

	// Create AWS client
	awsClient, err := aws.NewClient(ctx, aws.Options{
		Region:  viper.GetString("aws.region"),
		Profile: viper.GetString("aws.profile"),
	})
	if err != nil {
		return fmt.Errorf("failed to create AWS client: %w", err)
	}

	s3Service := aws.NewS3Service(awsClient)

	fmt.Printf("Checking start requests for project: %s\n", project)

	// Check for start requests
	requests, err := s3Service.CheckStartRequests(ctx, syncConfig.Bucket, project)
	if err != nil {
		return fmt.Errorf("failed to check start requests: %w", err)
	}

	if len(requests) == 0 {
		fmt.Printf("No pending start requests.\n")
		return nil
	}

	fmt.Printf("Found %d pending start request(s):\n\n", len(requests))
	fmt.Printf("%-15s %-20s %-15s\n", "USERNAME", "REQUESTED AT", "STUDENT ID")
	fmt.Println(strings.Repeat("-", 55))

	var usersToStart []string
	for username, request := range requests {
		fmt.Printf("%-15s %-20s %-15s\n",
			username,
			request.RequestedAt.Format("15:04:05"),
			request.StudentID)
		usersToStart = append(usersToStart, username)
	}

	if autoApprove {
		fmt.Printf("\nAuto-approving %d start requests...\n", len(usersToStart))

		// Start the requested instances
		err = startInstances(ctx, usersToStart, project, true)
		if err != nil {
			return fmt.Errorf("failed to start instances: %w", err)
		}

		// Clean up the start request files
		for username := range requests {
			_ = s3Service.DeleteStartRequest(ctx, syncConfig.Bucket, project, username)
		}

		fmt.Printf("✅ All start requests processed!\n")
	} else {
		fmt.Printf("\nApprove start requests? Run with --auto-approve or manually:\n")
		fmt.Printf("lfr instances start --users=%s --project=%s\n",
			strings.Join(usersToStart, ","), project)
	}

	return nil
}

// showStudentStatus shows comprehensive status for all students.
func showStudentStatus(ctx context.Context, project string) error {
	// Create AWS client
	awsClient, err := aws.NewClient(ctx, aws.Options{
		Region:  viper.GetString("aws.region"),
		Profile: viper.GetString("aws.profile"),
	})
	if err != nil {
		return fmt.Errorf("failed to create AWS client: %w", err)
	}

	lightsailService := aws.NewLightsailService(awsClient)

	// Get all instances for the project
	instances, err := lightsailService.ListInstances(ctx, project)
	if err != nil {
		return fmt.Errorf("failed to list instances: %w", err)
	}

	fmt.Printf("Student status for project: %s\n\n", project)
	fmt.Printf("%-15s %-20s %-12s %-18s %-15s\n",
		"STUDENT", "INSTANCE", "STATE", "PUBLIC IP", "LAST ACTIVITY")
	fmt.Println(strings.Repeat("-", 95))

	for _, instance := range instances {
		username := utils.ExtractUsernameFromInstance(instance.Name)
		if username == "" {
			continue
		}

		publicIP := instance.PublicIP
		if publicIP == "" {
			publicIP = "-"
		}

		// Calculate last activity (simplified)
		lastActivity := "Unknown"
		if instance.State == "running" {
			lastActivity = "Active now"
		} else {
			lastActivity = "Stopped"
		}

		fmt.Printf("%-15s %-20s %-12s %-18s %-15s\n",
			username,
			instance.Name,
			instance.State,
			publicIP,
			lastActivity,
		)
	}

	fmt.Printf("\nTotal: %d students\n", len(instances))

	return nil
}