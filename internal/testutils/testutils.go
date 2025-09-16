// Package testutils provides common testing utilities and helpers.
package testutils

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/lightsail"

	"github.com/scttfrdmn/lfr-tools/internal/types"
)

// MockIAMClient provides a mock implementation of the IAM client for testing.
type MockIAMClient struct {
	ListPoliciesFunc          func(ctx context.Context, params *iam.ListPoliciesInput, optFns ...func(*iam.Options)) (*iam.ListPoliciesOutput, error)
	CreatePolicyFunc          func(ctx context.Context, params *iam.CreatePolicyInput, optFns ...func(*iam.Options)) (*iam.CreatePolicyOutput, error)
	GetGroupFunc              func(ctx context.Context, params *iam.GetGroupInput, optFns ...func(*iam.Options)) (*iam.GetGroupOutput, error)
	CreateGroupFunc           func(ctx context.Context, params *iam.CreateGroupInput, optFns ...func(*iam.Options)) (*iam.CreateGroupOutput, error)
	AttachGroupPolicyFunc     func(ctx context.Context, params *iam.AttachGroupPolicyInput, optFns ...func(*iam.Options)) (*iam.AttachGroupPolicyOutput, error)
	CreateUserFunc            func(ctx context.Context, params *iam.CreateUserInput, optFns ...func(*iam.Options)) (*iam.CreateUserOutput, error)
	CreateLoginProfileFunc    func(ctx context.Context, params *iam.CreateLoginProfileInput, optFns ...func(*iam.Options)) (*iam.CreateLoginProfileOutput, error)
	AddUserToGroupFunc        func(ctx context.Context, params *iam.AddUserToGroupInput, optFns ...func(*iam.Options)) (*iam.AddUserToGroupOutput, error)
	PutUserPolicyFunc         func(ctx context.Context, params *iam.PutUserPolicyInput, optFns ...func(*iam.Options)) (*iam.PutUserPolicyOutput, error)
	ListGroupsForUserFunc     func(ctx context.Context, params *iam.ListGroupsForUserInput, optFns ...func(*iam.Options)) (*iam.ListGroupsForUserOutput, error)
	RemoveUserFromGroupFunc   func(ctx context.Context, params *iam.RemoveUserFromGroupInput, optFns ...func(*iam.Options)) (*iam.RemoveUserFromGroupOutput, error)
	DeleteLoginProfileFunc    func(ctx context.Context, params *iam.DeleteLoginProfileInput, optFns ...func(*iam.Options)) (*iam.DeleteLoginProfileOutput, error)
	ListUserPoliciesFunc      func(ctx context.Context, params *iam.ListUserPoliciesInput, optFns ...func(*iam.Options)) (*iam.ListUserPoliciesOutput, error)
	DeleteUserPolicyFunc      func(ctx context.Context, params *iam.DeleteUserPolicyInput, optFns ...func(*iam.Options)) (*iam.DeleteUserPolicyOutput, error)
	DeleteUserFunc            func(ctx context.Context, params *iam.DeleteUserInput, optFns ...func(*iam.Options)) (*iam.DeleteUserOutput, error)
	GetUserFunc               func(ctx context.Context, params *iam.GetUserInput, optFns ...func(*iam.Options)) (*iam.GetUserOutput, error)
	ListUserTagsFunc          func(ctx context.Context, params *iam.ListUserTagsInput, optFns ...func(*iam.Options)) (*iam.ListUserTagsOutput, error)
	ListAttachedGroupPoliciesFunc func(ctx context.Context, params *iam.ListAttachedGroupPoliciesInput, optFns ...func(*iam.Options)) (*iam.ListAttachedGroupPoliciesOutput, error)
}

// MockLightsailClient provides a mock implementation of the Lightsail client for testing.
type MockLightsailClient struct {
	GetBlueprintsFunc     func(ctx context.Context, params *lightsail.GetBlueprintsInput, optFns ...func(*lightsail.Options)) (*lightsail.GetBlueprintsOutput, error)
	GetBundlesFunc        func(ctx context.Context, params *lightsail.GetBundlesInput, optFns ...func(*lightsail.Options)) (*lightsail.GetBundlesOutput, error)
	GetRegionsFunc        func(ctx context.Context, params *lightsail.GetRegionsInput, optFns ...func(*lightsail.Options)) (*lightsail.GetRegionsOutput, error)
	CreateInstancesFunc   func(ctx context.Context, params *lightsail.CreateInstancesInput, optFns ...func(*lightsail.Options)) (*lightsail.CreateInstancesOutput, error)
	GetInstanceFunc       func(ctx context.Context, params *lightsail.GetInstanceInput, optFns ...func(*lightsail.Options)) (*lightsail.GetInstanceOutput, error)
	GetInstancesFunc      func(ctx context.Context, params *lightsail.GetInstancesInput, optFns ...func(*lightsail.Options)) (*lightsail.GetInstancesOutput, error)
	DeleteInstanceFunc    func(ctx context.Context, params *lightsail.DeleteInstanceInput, optFns ...func(*lightsail.Options)) (*lightsail.DeleteInstanceOutput, error)
	StartInstanceFunc     func(ctx context.Context, params *lightsail.StartInstanceInput, optFns ...func(*lightsail.Options)) (*lightsail.StartInstanceOutput, error)
	StopInstanceFunc      func(ctx context.Context, params *lightsail.StopInstanceInput, optFns ...func(*lightsail.Options)) (*lightsail.StopInstanceOutput, error)
	GetKeyPairFunc        func(ctx context.Context, params *lightsail.GetKeyPairInput, optFns ...func(*lightsail.Options)) (*lightsail.GetKeyPairOutput, error)
}

// TestFixture provides common test data and utilities.
type TestFixture struct {
	T *testing.T
}

// NewTestFixture creates a new test fixture.
func NewTestFixture(t *testing.T) *TestFixture {
	return &TestFixture{T: t}
}

// TestProject returns a sample project for testing.
func (f *TestFixture) TestProject() types.Project {
	return types.Project{
		Name:      "test-project",
		Blueprint: "ubuntu_22_04",
		Bundle:    "nano_2_0",
		Region:    "us-east-1",
		CreatedAt: time.Now(),
	}
}

// TestUser returns a sample user for testing.
func (f *TestFixture) TestUser() types.User {
	return types.User{
		Username:     "test-user",
		Project:      "test-project",
		InstanceARN:  "arn:aws:lightsail:us-east-1:123456789012:Instance/test-instance",
		InstanceName: "test-user-ubuntu_22_04",
		Password:     "test-password-123",
		CreatedAt:    time.Now(),
	}
}

// TestGroup returns a sample group for testing.
func (f *TestFixture) TestGroup() types.Group {
	return types.Group{
		Name:        "test-group",
		Policies:    []string{"arn:aws:iam::aws:policy/ReadOnlyAccess"},
		Description: "Test group for unit tests",
		CreatedAt:   time.Now(),
	}
}

// TestInstance returns a sample instance for testing.
func (f *TestFixture) TestInstance() types.Instance {
	return types.Instance{
		Name:      "test-instance",
		ARN:       "arn:aws:lightsail:us-east-1:123456789012:Instance/test-instance",
		State:     "running",
		Blueprint: "ubuntu_22_04",
		Bundle:    "nano_2_0",
		Region:    "us-east-1",
		Tags: map[string]string{
			"Project": "test-project",
		},
		CreatedAt: time.Now(),
		PublicIP:  "203.0.113.1",
		PrivateIP: "10.0.1.100",
	}
}

// AssertNoError fails the test if err is not nil.
func (f *TestFixture) AssertNoError(err error) {
	f.T.Helper()
	if err != nil {
		f.T.Fatalf("unexpected error: %v", err)
	}
}

// AssertError fails the test if err is nil.
func (f *TestFixture) AssertError(err error) {
	f.T.Helper()
	if err == nil {
		f.T.Fatal("expected error but got nil")
	}
}

// AssertEqual fails the test if expected != actual.
func (f *TestFixture) AssertEqual(expected, actual interface{}) {
	f.T.Helper()
	if expected != actual {
		f.T.Fatalf("expected %v, got %v", expected, actual)
	}
}

// AssertContains fails the test if haystack doesn't contain needle.
func (f *TestFixture) AssertContains(haystack, needle string) {
	f.T.Helper()
	if needle == "" || haystack == "" {
		f.T.Fatal("empty strings not allowed in AssertContains")
	}
	if !contains(haystack, needle) {
		f.T.Fatalf("expected %q to contain %q", haystack, needle)
	}
}

// contains checks if a string contains a substring.
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(substr) <= len(s) && s[0:len(substr)] == substr) ||
		(len(substr) < len(s) && containsHelper(s, substr)))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// SetupTestContext creates a context with timeout for testing.
func SetupTestContext() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	return ctx
}

// MockAWSConfig returns a basic AWS config for testing.
func MockAWSConfig() aws.Config {
	return aws.Config{
		Region: "us-east-1",
	}
}