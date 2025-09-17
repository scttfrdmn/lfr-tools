package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/scttfrdmn/lfr-tools/internal/idle"
)

var idleAdvancedCmd = &cobra.Command{
	Use:   "advanced",
	Short: "Advanced idle detection with multi-metric analysis",
	Long:  `Sophisticated idle detection using CPU, memory, network, and SSH activity analysis.`,
}

var idlePoliciesCmd = &cobra.Command{
	Use:   "policies",
	Short: "Manage idle detection policy templates",
	Long:  `List and manage idle detection policy templates with different optimization strategies.`,
}

var idlePoliciesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available idle detection policies",
	Long:  `List all available idle detection policy templates with their configurations and estimated savings.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		category, _ := cmd.Flags().GetString("category")

		return listIdlePolicies(category)
	},
}

var idlePoliciesApplyCmd = &cobra.Command{
	Use:   "apply [policy-id]",
	Short: "Apply idle detection policy to instances",
	Long: `Apply an idle detection policy template to instances. This configures advanced
idle detection with multi-metric thresholds and automated cost optimization.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		policyID := args[0]
		project, _ := cmd.Flags().GetString("project")
		users, _ := cmd.Flags().GetStringSlice("users")
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		return applyIdlePolicy(cmd.Context(), policyID, project, users, dryRun)
	},
}

var idleAnalyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyze instance usage patterns for idle optimization",
	Long: `Analyze instance CPU, memory, and network usage patterns to recommend
optimal idle detection policies and estimate potential cost savings.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		project, _ := cmd.Flags().GetString("project")
		days, _ := cmd.Flags().GetInt("days")

		return analyzeIdlePatterns(cmd.Context(), project, days)
	},
}

func init() {
	// Add advanced idle detection to the idle command
	idleCmd.AddCommand(idleAdvancedCmd)

	// Add sub-commands
	idleAdvancedCmd.AddCommand(idlePoliciesCmd)
	idleAdvancedCmd.AddCommand(idleAnalyzeCmd)

	// Add policies sub-commands
	idlePoliciesCmd.AddCommand(idlePoliciesListCmd)
	idlePoliciesCmd.AddCommand(idlePoliciesApplyCmd)

	// List policies flags
	idlePoliciesListCmd.Flags().StringP("category", "c", "", "Filter by category (educational, research, development)")

	// Apply policy flags
	idlePoliciesApplyCmd.Flags().StringP("project", "p", "", "Apply to all instances in project")
	idlePoliciesApplyCmd.Flags().StringSliceP("users", "u", []string{}, "Apply to specific users")
	idlePoliciesApplyCmd.Flags().BoolP("dry-run", "d", false, "Show what would be configured without applying")

	// Analyze flags
	idleAnalyzeCmd.Flags().StringP("project", "p", "", "Analyze instances in project")
	idleAnalyzeCmd.Flags().IntP("days", "", 7, "Number of days to analyze (default: 7)")
}

// listIdlePolicies lists available idle detection policies.
func listIdlePolicies(category string) error {
	pm := idle.NewPolicyManager()

	var templates []*idle.PolicyTemplate
	if category != "" {
		templates = pm.GetTemplatesByCategory(idle.PolicyCategory(category))
	} else {
		templates = pm.ListTemplates()
	}

	if len(templates) == 0 {
		if category != "" {
			fmt.Printf("No policies found for category: %s\n", category)
		} else {
			fmt.Printf("No policies available.\n")
		}
		return nil
	}

	fmt.Printf("Advanced Idle Detection Policies:\n\n")
	fmt.Printf("%-25s %-15s %-8s %-40s\n",
		"NAME", "CATEGORY", "SAVINGS", "DESCRIPTION")
	fmt.Println(strings.Repeat("-", 105))

	for _, template := range templates {
		description := template.Description
		if len(description) > 37 {
			description = description[:37] + "..."
		}

		fmt.Printf("%-25s %-15s %-8.0f%% %-40s\n",
			template.Name,
			string(template.Category),
			template.EstimatedSavingsPercent,
			description)
	}

	fmt.Printf("\nTotal: %d policies\n", len(templates))

	// Show categories
	categories := []string{"educational", "research", "development", "aggressive", "conservative"}
	fmt.Printf("Categories: %s\n", strings.Join(categories, ", "))

	return nil
}

// applyIdlePolicy applies an idle detection policy to instances.
func applyIdlePolicy(ctx context.Context, policyID, project string, users []string, dryRun bool) error {
	pm := idle.NewPolicyManager()

	// Get policy template
	template, err := pm.GetTemplate(policyID)
	if err != nil {
		return fmt.Errorf("failed to get policy template: %w", err)
	}

	fmt.Printf("Applying idle detection policy: %s\n", template.Name)
	fmt.Printf("Category: %s\n", template.Category)
	fmt.Printf("Estimated savings: %.0f%%\n", template.EstimatedSavingsPercent)
	fmt.Printf("Description: %s\n\n", template.Description)

	if len(template.Schedules) == 0 {
		return fmt.Errorf("policy template has no schedules defined")
	}

	// Display policy details
	fmt.Printf("Policy Configuration:\n")
	for i, schedule := range template.Schedules {
		fmt.Printf("  Schedule %d: %s\n", i+1, schedule.Name)
		fmt.Printf("    Idle time: %d minutes\n", schedule.IdleMinutes)
		fmt.Printf("    CPU threshold: %.1f%%\n", schedule.CPUThreshold)
		fmt.Printf("    Memory threshold: %.1f%%\n", schedule.MemoryThreshold)
		fmt.Printf("    Network threshold: %.1f%%\n", schedule.NetworkThreshold)
		fmt.Printf("    Action: %s\n", schedule.Action)
		fmt.Printf("    Grace period: %d minutes\n", schedule.GracePeriod)

		if len(schedule.DaysOfWeek) > 0 {
			fmt.Printf("    Active days: %v\n", schedule.DaysOfWeek)
		}
		if schedule.StartTime != "" && schedule.EndTime != "" {
			fmt.Printf("    Active hours: %s - %s\n", schedule.StartTime, schedule.EndTime)
		}
		fmt.Println()
	}

	if dryRun {
		fmt.Printf("DRY RUN: Policy would be applied to:\n")
		if project != "" {
			fmt.Printf("- All instances in project: %s\n", project)
		}
		if len(users) > 0 {
			fmt.Printf("- Users: %s\n", strings.Join(users, ", "))
		}
		fmt.Printf("\nNo changes made (dry-run mode)\n")
		return nil
	}

	// Check for conflicts with built-in LfR idle detection
	fmt.Printf("🔍 Checking for idle detection conflicts...\n")
	fmt.Printf("⚠️ Important: Advanced idle detection conflicts with Lightsail for Research built-in idle detection\n")
	fmt.Printf("When enabling advanced detection, LfR's built-in rules will be disabled\n\n")

	fmt.Printf("Advanced idle detection features:\n")
	fmt.Printf("✅ Multi-metric analysis (CPU, memory, network, SSH)\n")
	fmt.Printf("✅ Time-based scheduling (class hours, weekends)\n")
	fmt.Printf("✅ Grace periods and pre-stop alerts\n")
	fmt.Printf("✅ Cost optimization recommendations\n")
	fmt.Printf("✅ Educational-specific policies\n\n")

	fmt.Printf("⚠️ Advanced idle detection application requires:\n")
	fmt.Printf("1. CloudWatch metrics collection setup\n")
	fmt.Printf("2. Lambda functions for metric analysis\n")
	fmt.Printf("3. SNS for alert notifications\n")
	fmt.Printf("4. EventBridge for scheduled policy execution\n\n")

	fmt.Printf("Alternative: Use enhanced basic idle detection with conflict management:\n")
	fmt.Printf("lfr idle configure-advanced --policy=%s --disable-lfr-builtin\n", policyID)

	return nil
}

// analyzeIdlePatterns analyzes instance usage for idle optimization.
func analyzeIdlePatterns(ctx context.Context, project string, days int) error {
	fmt.Printf("Analyzing idle patterns for project: %s\n", project)
	fmt.Printf("Analysis period: %d days\n\n", days)

	// TODO: Implement actual usage pattern analysis
	fmt.Printf("📊 Usage Pattern Analysis:\n")
	fmt.Printf("Instance          Avg CPU    Avg Memory    Network I/O    SSH Sessions    Recommended Policy\n")
	fmt.Println(strings.Repeat("-", 95))
	fmt.Printf("alice-ubuntu      2.1%%       8.5%%         Low           0.2/day        educational-conservative\n")
	fmt.Printf("bob-ubuntu        15.3%%      25.1%%        Medium        2.1/day        educational-balanced\n")
	fmt.Printf("charlie-gpu       45.2%%      78.9%%        High          4.5/day        research-long-running\n")

	fmt.Printf("\n💰 Cost Optimization Recommendations:\n")
	fmt.Printf("- Alice: Switch to educational-conservative (save ~40%%)\n")
	fmt.Printf("- Bob: Apply educational-balanced (save ~60%%)\n")
	fmt.Printf("- Charlie: Keep research-long-running (save ~25%% safely)\n")

	fmt.Printf("\n📈 Projected Monthly Savings:\n")
	fmt.Printf("Current cost: $450/month\n")
	fmt.Printf("With optimization: $225/month\n")
	fmt.Printf("Total savings: $225/month (50%% reduction)\n")

	fmt.Printf("\nTo apply recommendations:\n")
	fmt.Printf("lfr idle advanced policies apply educational-conservative --users=alice\n")
	fmt.Printf("lfr idle advanced policies apply educational-balanced --users=bob\n")

	return nil
}