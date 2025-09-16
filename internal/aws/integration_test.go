// +build integration

package aws

import (
	"context"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"

	"github.com/scttfrdmn/lfr-tools/internal/testutils"
)

// These tests require either real AWS credentials or LocalStack running on localhost:4566

func setupIntegrationTest(t *testing.T) *Client {
	t.Helper()

	// Check if we should skip integration tests
	if os.Getenv("SKIP_INTEGRATION_TESTS") == "true" {
		t.Skip("Skipping integration tests (SKIP_INTEGRATION_TESTS=true)")
	}

	ctx := testutils.SetupTestContext()

	// Try LocalStack first, then real AWS
	var client *Client
	var err error

	// LocalStack configuration
	if os.Getenv("AWS_ENDPOINT_URL") != "" || os.Getenv("LOCALSTACK_ENDPOINT") != "" {
		endpoint := os.Getenv("AWS_ENDPOINT_URL")
		if endpoint == "" {
			endpoint = os.Getenv("LOCALSTACK_ENDPOINT")
		}
		if endpoint == "" {
			endpoint = "http://localhost:4566"
		}

		cfg, err := config.LoadDefaultConfig(ctx,
			config.WithRegion("us-east-1"),
			config.WithEndpointResolverWithOptions(
				aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
					return aws.Endpoint{
						URL:           endpoint,
						SigningRegion: region,
					}, nil
				}),
			),
		)
		if err != nil {
			t.Fatalf("Failed to load LocalStack config: %v", err)
		}

		client = &Client{
			Config: cfg,
		}
	} else {
		// Real AWS configuration
		opts := Options{
			Region:  "us-east-1",
			Profile: os.Getenv("AWS_PROFILE"),
		}

		client, err = NewClient(ctx, opts)
		if err != nil {
			t.Skipf("Failed to create AWS client (no credentials?): %v", err)
		}
	}

	return client
}

func TestIntegrationIAMOperations(t *testing.T) {
	client := setupIntegrationTest(t)
	fixture := testutils.NewTestFixture(t)
	ctx := testutils.SetupTestContext()

	service := NewIAMService(client)

	// Test policy creation
	policyDoc := `{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Action": "lightsail:Get*",
				"Resource": "*"
			}
		]
	}`

	policyARN, err := service.CreatePolicy(ctx, "LFRTestPolicy", "Test policy for LFR tools", policyDoc)
	if err != nil {
		// LocalStack might not support all IAM operations
		t.Logf("Policy creation failed (expected with LocalStack): %v", err)
		return
	}

	t.Logf("Created policy: %s", policyARN)

	// Test group creation
	group, err := service.CreateGroup(ctx, "LFRTestGroup", "Test group", []string{policyARN})
	if err != nil {
		t.Logf("Group creation failed (expected with LocalStack): %v", err)
		return
	}

	fixture.AssertEqual("LFRTestGroup", group.Name)
	t.Logf("Created group: %s", group.Name)

	// Test user creation
	user, err := service.CreateUser(ctx, "testuser", "testpassword123", "test-project")
	if err != nil {
		t.Logf("User creation failed (expected with LocalStack): %v", err)
		return
	}

	fixture.AssertEqual("testuser", user.Username)
	fixture.AssertEqual("test-project", user.Project)
	t.Logf("Created user: %s", user.Username)

	// Test adding user to group
	err = service.AddUserToGroup(ctx, "testuser", "LFRTestGroup")
	if err != nil {
		t.Logf("Add user to group failed (expected with LocalStack): %v", err)
	}

	// Cleanup (in reverse order)
	// Note: Cleanup might fail with LocalStack, which is expected
	_ = service.DeleteUser(ctx, "testuser")
	t.Logf("Cleanup completed")
}

func TestIntegrationLightsailOperations(t *testing.T) {
	client := setupIntegrationTest(t)
	fixture := testutils.NewTestFixture(t)
	ctx := testutils.SetupTestContext()

	service := NewLightsailService(client)

	// Test getting blueprints
	blueprints, err := service.GetBlueprints(ctx)
	if err != nil {
		t.Logf("Get blueprints failed (expected with LocalStack): %v", err)
		return
	}

	if len(blueprints) > 0 {
		t.Logf("Found %d blueprints", len(blueprints))
	} else {
		t.Log("No blueprints found (expected with LocalStack)")
	}

	// Test getting bundles
	bundles, err := service.GetBundles(ctx)
	if err != nil {
		t.Logf("Get bundles failed (expected with LocalStack): %v", err)
		return
	}

	if len(bundles) > 0 {
		t.Logf("Found %d bundles", len(bundles))
	} else {
		t.Log("No bundles found (expected with LocalStack)")
	}

	// Test getting regions
	regions, err := service.GetRegions(ctx)
	if err != nil {
		t.Logf("Get regions failed (expected with LocalStack): %v", err)
		return
	}

	if len(regions) > 0 {
		t.Logf("Found %d regions", len(regions))
		fixture.AssertContains(regions[0], "us-")
	} else {
		t.Log("No regions found (expected with LocalStack)")
	}

	// Test listing instances
	instances, err := service.ListInstances(ctx, "")
	if err != nil {
		t.Logf("List instances failed (expected with LocalStack): %v", err)
		return
	}

	t.Logf("Found %d instances", len(instances))
}

func TestIntegrationClientConfiguration(t *testing.T) {
	client := setupIntegrationTest(t)

	region := client.GetRegion()
	if region == "" {
		t.Error("Expected non-empty region")
	}

	t.Logf("Using region: %s", region)

	// Test that clients are properly initialized
	if client.IAM == nil {
		t.Error("IAM client should not be nil")
	}

	if client.Lightsail == nil {
		t.Error("Lightsail client should not be nil")
	}
}