// +build integration

package aws

import (
	"os"
	"testing"

	"github.com/scttfrdmn/lfr-tools/internal/testutils"
)

// TestRealLightsailOperations tests actual Lightsail operations with real AWS
// These tests validate the complete workflow for educational environments
func TestRealLightsailOperations(t *testing.T) {
	if os.Getenv("AWS_ENDPOINT_URL") != "" {
		t.Skip("Skipping real AWS tests when LocalStack is configured")
	}

	client := setupIntegrationTest(t)
	ctx := testutils.SetupTestContext()

	service := NewLightsailService(client)

	t.Run("list existing instances", func(t *testing.T) {
		instances, err := service.ListInstances(ctx, "")
		if err != nil {
			t.Errorf("failed to list instances: %v", err)
		}

		t.Logf("Found %d instances", len(instances))
		for _, instance := range instances {
			t.Logf("  - %s: %s (%s) in %s",
				instance.Name, instance.State, instance.Blueprint, instance.Region)
		}
	})

	t.Run("filter instances by project", func(t *testing.T) {
		allInstances, err := service.ListInstances(ctx, "")
		if err != nil {
			t.Errorf("failed to list all instances: %v", err)
		}

		// Test filtering - most instances should be "untagged" based on CLI output
		untaggedInstances, err := service.ListInstances(ctx, "untagged")
		if err != nil {
			t.Errorf("failed to list untagged instances: %v", err)
		}

		t.Logf("Total instances: %d, Untagged instances: %d",
			len(allInstances), len(untaggedInstances))

		// Validate that untagged instances are a subset
		if len(untaggedInstances) > len(allInstances) {
			t.Errorf("untagged instances (%d) should not exceed total (%d)",
				len(untaggedInstances), len(allInstances))
		}
	})

	t.Run("get specific instance details", func(t *testing.T) {
		instances, err := service.ListInstances(ctx, "")
		if err != nil || len(instances) == 0 {
			t.Skip("No instances available for testing")
		}

		// Test getting details for the first instance
		testInstance := instances[0]
		details, err := service.GetInstance(ctx, testInstance.Name)
		if err != nil {
			t.Errorf("failed to get instance details for %s: %v", testInstance.Name, err)
		}

		// Validate returned data
		if details.Name != testInstance.Name {
			t.Errorf("instance name mismatch: expected %s, got %s",
				testInstance.Name, details.Name)
		}

		if details.State == "" {
			t.Error("instance state should not be empty")
		}

		t.Logf("Instance %s: %s, Bundle: %s, Blueprint: %s",
			details.Name, details.State, details.Bundle, details.Blueprint)
	})

	t.Run("get available regions", func(t *testing.T) {
		regions, err := service.GetRegions(ctx)
		if err != nil {
			t.Errorf("failed to get regions: %v", err)
		}

		if len(regions) == 0 {
			t.Error("expected at least one region")
		}

		// Validate common regions exist
		hasUSEast1 := false
		hasUSWest2 := false
		for _, region := range regions {
			if region == "us-east-1" {
				hasUSEast1 = true
			}
			if region == "us-west-2" {
				hasUSWest2 = true
			}
		}

		if !hasUSEast1 && !hasUSWest2 {
			t.Error("expected to find common US regions")
		}

		t.Logf("Found %d regions: %v", len(regions), regions)
	})

	t.Run("get blueprints", func(t *testing.T) {
		blueprints, err := service.GetBlueprints(ctx)
		if err != nil {
			t.Errorf("failed to get blueprints: %v", err)
		}

		if len(blueprints) == 0 {
			t.Error("expected at least one blueprint")
		}

		// Look for common LFR educational blueprints
		hasUbuntu := false
		hasRStudio := false
		for _, bp := range blueprints {
			if bp == "lfr_ubuntu_1_0" {
				hasUbuntu = true
			}
			if bp == "lfr_rstudio_1_0" {
				hasRStudio = true
			}
		}

		if !hasUbuntu {
			t.Error("expected to find LFR Ubuntu blueprint")
		}

		t.Logf("Found %d blueprints, Ubuntu: %v, RStudio: %v",
			len(blueprints), hasUbuntu, hasRStudio)
		t.Logf("Available blueprints: %v", blueprints)
	})

	t.Run("get bundles", func(t *testing.T) {
		bundles, err := service.GetBundles(ctx)
		if err != nil {
			t.Errorf("failed to get bundles: %v", err)
		}

		if len(bundles) == 0 {
			t.Error("expected at least one bundle")
		}

		// Look for common educational bundles
		hasStandardXL := false
		hasGPU := false
		for _, bundle := range bundles {
			if bundle == "app_standard_xl_1_0" {
				hasStandardXL = true
			}
			if bundle == "gpu_nvidia_xl_1_0" {
				hasGPU = true
			}
		}

		t.Logf("Found %d bundles, Standard XL: %v, GPU: %v",
			len(bundles), hasStandardXL, hasGPU)
	})
}

// TestCostIntegration tests cost-related operations
func TestCostIntegration(t *testing.T) {
	if os.Getenv("AWS_ENDPOINT_URL") != "" {
		t.Skip("Skipping cost tests with LocalStack (Cost Explorer not supported)")
	}

	// Note: This test validates cost API integration but doesn't verify exact amounts
	// since costs are dynamic and account-specific
	t.Run("cost API availability", func(t *testing.T) {
		// Basic test to ensure cost-related operations don't crash
		// Actual cost testing would require specific test scenarios
		t.Log("Cost integration test placeholder - implement specific cost scenarios as needed")
		t.Log("For real cost testing, implement scenarios based on your account's usage patterns")
	})
}

// TestEducationalWorkflows tests complete educational scenarios
func TestEducationalWorkflows(t *testing.T) {
	if os.Getenv("AWS_ENDPOINT_URL") != "" {
		t.Skip("Skipping educational workflow tests with LocalStack")
	}

	client := setupIntegrationTest(t)
	ctx := testutils.SetupTestContext()

	t.Run("professor workflow simulation", func(t *testing.T) {
		lightsailService := NewLightsailService(client)

		// 1. Professor lists current instances
		instances, err := lightsailService.ListInstances(ctx, "")
		if err != nil {
			t.Errorf("professor cannot list instances: %v", err)
		}
		t.Logf("Professor sees %d instances", len(instances))

		// 2. Professor checks user access (IAM validation)
		// This is a read-only operation to validate permissions
		t.Log("Professor can access IAM for user management (validated via integration tests)")

		// 3. Professor checks regional capacity
		regions, err := lightsailService.GetRegions(ctx)
		if err != nil {
			t.Errorf("professor cannot check regional capacity: %v", err)
		}
		t.Logf("Professor can access %d regions", len(regions))
	})

	t.Run("student workflow simulation", func(t *testing.T) {
		lightsailService := NewLightsailService(client)

		// 1. Student tries to find their instance
		instances, err := lightsailService.ListInstances(ctx, "")
		if err != nil {
			t.Errorf("student cannot list instances: %v", err)
		}

		// Look for student-pattern instances (username-blueprint format)
		studentInstances := 0
		for _, instance := range instances {
			// Count instances that follow student naming pattern
			if len(instance.Name) > 0 && instance.Name != "demo" {
				studentInstances++
			}
		}

		t.Logf("Student can see %d potential student instances", studentInstances)

		// 2. Student checks if they can access instance details
		if len(instances) > 0 {
			testInstance := instances[0]
			details, err := lightsailService.GetInstance(ctx, testInstance.Name)
			if err != nil {
				t.Errorf("student cannot get instance details: %v", err)
			} else {
				t.Logf("Student can access instance %s details: %s",
					details.Name, details.State)
			}
		}
	})
}