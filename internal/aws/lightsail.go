package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lightsail"
	lightsailTypes "github.com/aws/aws-sdk-go-v2/service/lightsail/types"

	"github.com/scttfrdmn/lfr-tools/internal/types"
)

// LightsailService provides Lightsail operations.
type LightsailService struct {
	client *Client
}

// NewLightsailService creates a new Lightsail service.
func NewLightsailService(client *Client) *LightsailService {
	return &LightsailService{
		client: client,
	}
}

// GetBlueprints retrieves available Lightsail blueprints for LfR.
func (s *LightsailService) GetBlueprints(ctx context.Context) ([]string, error) {
	output, err := s.client.Lightsail.GetBlueprints(ctx, &lightsail.GetBlueprintsInput{
		AppCategory: lightsailTypes.AppCategoryLfR,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get blueprints: %w", err)
	}

	var blueprints []string
	for _, blueprint := range output.Blueprints {
		blueprints = append(blueprints, aws.ToString(blueprint.BlueprintId))
	}

	return blueprints, nil
}

// GetBundles retrieves available Lightsail bundles for LfR.
func (s *LightsailService) GetBundles(ctx context.Context) ([]string, error) {
	output, err := s.client.Lightsail.GetBundles(ctx, &lightsail.GetBundlesInput{
		AppCategory: lightsailTypes.AppCategoryLfR,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get bundles: %w", err)
	}

	var bundles []string
	for _, bundle := range output.Bundles {
		bundles = append(bundles, aws.ToString(bundle.BundleId))
	}

	return bundles, nil
}

// GetRegions retrieves available Lightsail regions.
func (s *LightsailService) GetRegions(ctx context.Context) ([]string, error) {
	output, err := s.client.Lightsail.GetRegions(ctx, &lightsail.GetRegionsInput{})
	if err != nil {
		return nil, fmt.Errorf("failed to get regions: %w", err)
	}

	var regions []string
	for _, region := range output.Regions {
		regions = append(regions, aws.ToString(region.Name))
	}

	return regions, nil
}

// CreateInstance creates a Lightsail instance.
func (s *LightsailService) CreateInstance(ctx context.Context, name, blueprintID, bundleID, availabilityZone, project string) (*types.Instance, error) {
	_, err := s.client.Lightsail.CreateInstances(ctx, &lightsail.CreateInstancesInput{
		InstanceNames:    []string{name},
		BlueprintId:      aws.String(blueprintID),
		BundleId:         aws.String(bundleID),
		AvailabilityZone: aws.String(availabilityZone),
		Tags: []lightsailTypes.Tag{
			{
				Key:   aws.String("Project"),
				Value: aws.String(project),
			},
		},
		AddOns: []lightsailTypes.AddOnRequest{
			{
				AddOnType: lightsailTypes.AddOnTypeStopInstanceOnIdle,
				StopInstanceOnIdleRequest: &lightsailTypes.StopInstanceOnIdleRequest{
					Threshold: aws.String("2"),   // 2 hours
					Duration:  aws.String("30"),  // 30 minutes
				},
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create instance %s: %w", name, err)
	}

	// Get the created instance details
	return s.GetInstance(ctx, name)
}

// GetInstance retrieves instance details.
func (s *LightsailService) GetInstance(ctx context.Context, name string) (*types.Instance, error) {
	output, err := s.client.Lightsail.GetInstance(ctx, &lightsail.GetInstanceInput{
		InstanceName: aws.String(name),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get instance %s: %w", name, err)
	}

	instance := output.Instance
	tags := make(map[string]string)
	for _, tag := range instance.Tags {
		tags[aws.ToString(tag.Key)] = aws.ToString(tag.Value)
	}

	result := &types.Instance{
		Name:      aws.ToString(instance.Name),
		ARN:       aws.ToString(instance.Arn),
		State:     string(instance.State.Name),
		Blueprint: aws.ToString(instance.BlueprintName),
		Bundle:    aws.ToString(instance.BundleName),
		Region:    aws.ToString(instance.Location.RegionName),
		Tags:      tags,
		CreatedAt: aws.ToTime(instance.CreatedAt),
	}

	if instance.PublicIpAddress != nil {
		result.PublicIP = aws.ToString(instance.PublicIpAddress)
	}
	if instance.PrivateIpAddress != nil {
		result.PrivateIP = aws.ToString(instance.PrivateIpAddress)
	}

	return result, nil
}

// ListInstances lists all Lightsail instances, optionally filtered by project.
func (s *LightsailService) ListInstances(ctx context.Context, project string) ([]*types.Instance, error) {
	output, err := s.client.Lightsail.GetInstances(ctx, &lightsail.GetInstancesInput{})
	if err != nil {
		return nil, fmt.Errorf("failed to list instances: %w", err)
	}

	var instances []*types.Instance
	for _, instance := range output.Instances {
		tags := make(map[string]string)
		for _, tag := range instance.Tags {
			tags[aws.ToString(tag.Key)] = aws.ToString(tag.Value)
		}

		// Filter by project if specified
		if project != "" {
			if projectTag, exists := tags["Project"]; !exists || projectTag != project {
				continue
			}
		}

		result := &types.Instance{
			Name:      aws.ToString(instance.Name),
			ARN:       aws.ToString(instance.Arn),
			State:     string(instance.State.Name),
			Blueprint: aws.ToString(instance.BlueprintName),
			Bundle:    aws.ToString(instance.BundleName),
			Region:    aws.ToString(instance.Location.RegionName),
			Tags:      tags,
			CreatedAt: aws.ToTime(instance.CreatedAt),
		}

		if instance.PublicIpAddress != nil {
			result.PublicIP = aws.ToString(instance.PublicIpAddress)
		}
		if instance.PrivateIpAddress != nil {
			result.PrivateIP = aws.ToString(instance.PrivateIpAddress)
		}

		instances = append(instances, result)
	}

	return instances, nil
}

// DeleteInstance deletes a Lightsail instance.
func (s *LightsailService) DeleteInstance(ctx context.Context, name string) error {
	_, err := s.client.Lightsail.DeleteInstance(ctx, &lightsail.DeleteInstanceInput{
		InstanceName: aws.String(name),
	})
	if err != nil {
		return fmt.Errorf("failed to delete instance %s: %w", name, err)
	}

	return nil
}

// StartInstance starts a stopped Lightsail instance.
func (s *LightsailService) StartInstance(ctx context.Context, name string) error {
	_, err := s.client.Lightsail.StartInstance(ctx, &lightsail.StartInstanceInput{
		InstanceName: aws.String(name),
	})
	if err != nil {
		return fmt.Errorf("failed to start instance %s: %w", name, err)
	}

	return nil
}

// StopInstance stops a running Lightsail instance.
func (s *LightsailService) StopInstance(ctx context.Context, name string) error {
	_, err := s.client.Lightsail.StopInstance(ctx, &lightsail.StopInstanceInput{
		InstanceName: aws.String(name),
	})
	if err != nil {
		return fmt.Errorf("failed to stop instance %s: %w", name, err)
	}

	return nil
}

// GetInstanceSSHKeyPair retrieves the SSH key pair for an instance.
func (s *LightsailService) GetInstanceSSHKeyPair(ctx context.Context, keyPairName string) (*lightsailTypes.KeyPair, error) {
	output, err := s.client.Lightsail.GetKeyPair(ctx, &lightsail.GetKeyPairInput{
		KeyPairName: aws.String(keyPairName),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get key pair %s: %w", keyPairName, err)
	}

	return output.KeyPair, nil
}