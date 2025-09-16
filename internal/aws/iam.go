package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	iamTypes "github.com/aws/aws-sdk-go-v2/service/iam/types"

	"github.com/scttfrdmn/lfr-tools/internal/types"
)

// IAMService provides IAM operations.
type IAMService struct {
	client *Client
}

// NewIAMService creates a new IAM service.
func NewIAMService(client *Client) *IAMService {
	return &IAMService{
		client: client,
	}
}

// CreatePolicy creates an IAM policy if it doesn't exist.
func (s *IAMService) CreatePolicy(ctx context.Context, name, description, document string) (string, error) {
	// Check if policy already exists
	listOutput, err := s.client.IAM.ListPolicies(ctx, &iam.ListPoliciesInput{
		Scope: iamTypes.PolicyScopeTypeLocal,
	})
	if err != nil {
		return "", fmt.Errorf("failed to list policies: %w", err)
	}

	for _, policy := range listOutput.Policies {
		if aws.ToString(policy.PolicyName) == name {
			return aws.ToString(policy.Arn), nil
		}
	}

	// Create policy
	output, err := s.client.IAM.CreatePolicy(ctx, &iam.CreatePolicyInput{
		PolicyName:     aws.String(name),
		Description:    aws.String(description),
		PolicyDocument: aws.String(document),
	})
	if err != nil {
		return "", fmt.Errorf("failed to create policy %s: %w", name, err)
	}

	return aws.ToString(output.Policy.Arn), nil
}

// CreateGroup creates an IAM group if it doesn't exist.
func (s *IAMService) CreateGroup(ctx context.Context, name, description string, policyARNs []string) (*types.Group, error) {
	// Check if group already exists
	_, err := s.client.IAM.GetGroup(ctx, &iam.GetGroupInput{
		GroupName: aws.String(name),
	})
	if err == nil {
		// Group exists, return info
		return s.getGroupInfo(ctx, name)
	}

	// Create group
	_, err = s.client.IAM.CreateGroup(ctx, &iam.CreateGroupInput{
		GroupName: aws.String(name),
		Path:      aws.String("/"),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create group %s: %w", name, err)
	}

	// Attach policies
	for _, policyARN := range policyARNs {
		_, err = s.client.IAM.AttachGroupPolicy(ctx, &iam.AttachGroupPolicyInput{
			GroupName: aws.String(name),
			PolicyArn: aws.String(policyARN),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to attach policy %s to group %s: %w", policyARN, name, err)
		}
	}

	return s.getGroupInfo(ctx, name)
}

// CreateUser creates an IAM user with login profile.
func (s *IAMService) CreateUser(ctx context.Context, username, password, project string) (*types.User, error) {
	// Create user
	_, err := s.client.IAM.CreateUser(ctx, &iam.CreateUserInput{
		UserName: aws.String(username),
		Tags: []iamTypes.Tag{
			{
				Key:   aws.String("Project"),
				Value: aws.String(project),
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user %s: %w", username, err)
	}

	// Create login profile
	_, err = s.client.IAM.CreateLoginProfile(ctx, &iam.CreateLoginProfileInput{
		UserName:              aws.String(username),
		Password:              aws.String(password),
		PasswordResetRequired: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create login profile for user %s: %w", username, err)
	}

	return s.getUserInfo(ctx, username)
}

// AddUserToGroup adds a user to a group.
func (s *IAMService) AddUserToGroup(ctx context.Context, username, groupName string) error {
	_, err := s.client.IAM.AddUserToGroup(ctx, &iam.AddUserToGroupInput{
		UserName:  aws.String(username),
		GroupName: aws.String(groupName),
	})
	if err != nil {
		return fmt.Errorf("failed to add user %s to group %s: %w", username, groupName, err)
	}

	return nil
}

// PutUserPolicy creates or updates an inline policy for a user.
func (s *IAMService) PutUserPolicy(ctx context.Context, username, policyName, policyDocument string) error {
	_, err := s.client.IAM.PutUserPolicy(ctx, &iam.PutUserPolicyInput{
		UserName:       aws.String(username),
		PolicyName:     aws.String(policyName),
		PolicyDocument: aws.String(policyDocument),
	})
	if err != nil {
		return fmt.Errorf("failed to put user policy %s for user %s: %w", policyName, username, err)
	}

	return nil
}

// DeleteUser removes a user and all associated resources.
func (s *IAMService) DeleteUser(ctx context.Context, username string) error {
	// Remove user from all groups
	groups, err := s.client.IAM.ListGroupsForUser(ctx, &iam.ListGroupsForUserInput{
		UserName: aws.String(username),
	})
	if err != nil {
		return fmt.Errorf("failed to get groups for user %s: %w", username, err)
	}

	for _, group := range groups.Groups {
		_, err = s.client.IAM.RemoveUserFromGroup(ctx, &iam.RemoveUserFromGroupInput{
			UserName:  aws.String(username),
			GroupName: group.GroupName,
		})
		if err != nil {
			return fmt.Errorf("failed to remove user %s from group %s: %w", username, aws.ToString(group.GroupName), err)
		}
	}

	// Delete login profile
	_, err = s.client.IAM.DeleteLoginProfile(ctx, &iam.DeleteLoginProfileInput{
		UserName: aws.String(username),
	})
	if err != nil {
		// Login profile might not exist, continue
		fmt.Printf("Warning: failed to delete login profile for user %s: %v\n", username, err)
	}

	// Delete inline policies
	policies, err := s.client.IAM.ListUserPolicies(ctx, &iam.ListUserPoliciesInput{
		UserName: aws.String(username),
	})
	if err != nil {
		return fmt.Errorf("failed to list user policies for %s: %w", username, err)
	}

	for _, policyName := range policies.PolicyNames {
		_, err = s.client.IAM.DeleteUserPolicy(ctx, &iam.DeleteUserPolicyInput{
			UserName:   aws.String(username),
			PolicyName: aws.String(policyName),
		})
		if err != nil {
			return fmt.Errorf("failed to delete user policy %s for user %s: %w", policyName, username, err)
		}
	}

	// Delete user
	_, err = s.client.IAM.DeleteUser(ctx, &iam.DeleteUserInput{
		UserName: aws.String(username),
	})
	if err != nil {
		return fmt.Errorf("failed to delete user %s: %w", username, err)
	}

	return nil
}

// getGroupInfo retrieves group information.
func (s *IAMService) getGroupInfo(ctx context.Context, name string) (*types.Group, error) {
	output, err := s.client.IAM.GetGroup(ctx, &iam.GetGroupInput{
		GroupName: aws.String(name),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get group info for %s: %w", name, err)
	}

	// Get attached policies
	policies, err := s.client.IAM.ListAttachedGroupPolicies(ctx, &iam.ListAttachedGroupPoliciesInput{
		GroupName: aws.String(name),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list group policies for %s: %w", name, err)
	}

	var policyARNs []string
	for _, policy := range policies.AttachedPolicies {
		policyARNs = append(policyARNs, aws.ToString(policy.PolicyArn))
	}

	group := &types.Group{
		Name:        aws.ToString(output.Group.GroupName),
		Policies:    policyARNs,
		Description: aws.ToString(output.Group.Path), // Using path as description placeholder
		CreatedAt:   *output.Group.CreateDate,
	}

	return group, nil
}

// getUserInfo retrieves user information.
func (s *IAMService) getUserInfo(ctx context.Context, username string) (*types.User, error) {
	output, err := s.client.IAM.GetUser(ctx, &iam.GetUserInput{
		UserName: aws.String(username),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get user info for %s: %w", username, err)
	}

	// Get user tags
	tags, err := s.client.IAM.ListUserTags(ctx, &iam.ListUserTagsInput{
		UserName: aws.String(username),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get user tags for %s: %w", username, err)
	}

	var project string
	for _, tag := range tags.Tags {
		if aws.ToString(tag.Key) == "Project" {
			project = aws.ToString(tag.Value)
			break
		}
	}

	user := &types.User{
		Username:  aws.ToString(output.User.UserName),
		Project:   project,
		CreatedAt: *output.User.CreateDate,
	}

	return user, nil
}