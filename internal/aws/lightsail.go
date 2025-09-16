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
		regions = append(regions, string(region.Name))
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
		State:     aws.ToString(instance.State.Name),
		Blueprint: aws.ToString(instance.BlueprintName),
		Bundle:    aws.ToString(instance.BundleId),
		Region:    string(instance.Location.RegionName),
		Tags:      tags,
		CreatedAt: *instance.CreatedAt,
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
			State:     aws.ToString(instance.State.Name),
			Blueprint: aws.ToString(instance.BlueprintName),
			Bundle:    aws.ToString(instance.BundleId),
			Region:    string(instance.Location.RegionName),
			Tags:      tags,
			CreatedAt: *instance.CreatedAt,
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

// DownloadSSHKey downloads the private key for an instance.
func (s *LightsailService) DownloadSSHKey(ctx context.Context, keyPairName string) (string, error) {
	output, err := s.client.Lightsail.DownloadDefaultKeyPair(ctx, &lightsail.DownloadDefaultKeyPairInput{})
	if err != nil {
		return "", fmt.Errorf("failed to download default key pair: %w", err)
	}

	return aws.ToString(output.PrivateKeyBase64), nil
}

// GetInstanceAccessDetails retrieves SSH access information for an instance.
func (s *LightsailService) GetInstanceAccessDetails(ctx context.Context, instanceName string) (*lightsail.GetInstanceAccessDetailsOutput, error) {
	output, err := s.client.Lightsail.GetInstanceAccessDetails(ctx, &lightsail.GetInstanceAccessDetailsInput{
		InstanceName: aws.String(instanceName),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get instance access details for %s: %w", instanceName, err)
	}

	return output, nil
}

// CreateDisk creates a new block storage disk.
func (s *LightsailService) CreateDisk(ctx context.Context, diskName string, sizeGB int32, availabilityZone, project string) (*types.Disk, error) {
	_, err := s.client.Lightsail.CreateDisk(ctx, &lightsail.CreateDiskInput{
		DiskName:         aws.String(diskName),
		AvailabilityZone: aws.String(availabilityZone),
		SizeInGb:         aws.Int32(sizeGB),
		Tags: []lightsailTypes.Tag{
			{
				Key:   aws.String("Project"),
				Value: aws.String(project),
			},
			{
				Key:   aws.String("CreatedBy"),
				Value: aws.String("lfr-tools"),
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create disk %s: %w", diskName, err)
	}

	return s.GetDisk(ctx, diskName)
}

// GetDisk retrieves disk details.
func (s *LightsailService) GetDisk(ctx context.Context, diskName string) (*types.Disk, error) {
	output, err := s.client.Lightsail.GetDisk(ctx, &lightsail.GetDiskInput{
		DiskName: aws.String(diskName),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get disk %s: %w", diskName, err)
	}

	disk := output.Disk
	tags := make(map[string]string)
	for _, tag := range disk.Tags {
		tags[aws.ToString(tag.Key)] = aws.ToString(tag.Value)
	}

	return &types.Disk{
		Name:             aws.ToString(disk.Name),
		ARN:              aws.ToString(disk.Arn),
		State:            string(disk.State),
		SizeGB:           aws.ToInt32(disk.SizeInGb),
		IOPS:             aws.ToInt32(disk.Iops),
		Path:             aws.ToString(disk.Path),
		AttachedTo:       aws.ToString(disk.AttachedTo),
		AvailabilityZone: aws.ToString(disk.Location.AvailabilityZone),
		Region:           string(disk.Location.RegionName),
		Tags:             tags,
		CreatedAt:        *disk.CreatedAt,
	}, nil
}

// ListDisks lists all block storage disks.
func (s *LightsailService) ListDisks(ctx context.Context, project string) ([]*types.Disk, error) {
	output, err := s.client.Lightsail.GetDisks(ctx, &lightsail.GetDisksInput{})
	if err != nil {
		return nil, fmt.Errorf("failed to list disks: %w", err)
	}

	var disks []*types.Disk
	for _, disk := range output.Disks {
		tags := make(map[string]string)
		for _, tag := range disk.Tags {
			tags[aws.ToString(tag.Key)] = aws.ToString(tag.Value)
		}

		// Filter by project if specified
		if project != "" {
			if projectTag, exists := tags["Project"]; !exists || projectTag != project {
				continue
			}
		}

		result := &types.Disk{
			Name:             aws.ToString(disk.Name),
			ARN:              aws.ToString(disk.Arn),
			State:            string(disk.State),
			SizeGB:           aws.ToInt32(disk.SizeInGb),
			IOPS:             aws.ToInt32(disk.Iops),
			Path:             aws.ToString(disk.Path),
			AttachedTo:       aws.ToString(disk.AttachedTo),
			AvailabilityZone: aws.ToString(disk.Location.AvailabilityZone),
			Region:           string(disk.Location.RegionName),
			Tags:             tags,
			CreatedAt:        *disk.CreatedAt,
		}

		disks = append(disks, result)
	}

	return disks, nil
}

// AttachDisk attaches a disk to an instance.
func (s *LightsailService) AttachDisk(ctx context.Context, diskName, instanceName, diskPath string) error {
	_, err := s.client.Lightsail.AttachDisk(ctx, &lightsail.AttachDiskInput{
		DiskName:     aws.String(diskName),
		InstanceName: aws.String(instanceName),
		DiskPath:     aws.String(diskPath),
	})
	if err != nil {
		return fmt.Errorf("failed to attach disk %s to instance %s: %w", diskName, instanceName, err)
	}

	return nil
}

// DetachDisk detaches a disk from an instance.
func (s *LightsailService) DetachDisk(ctx context.Context, diskName string) error {
	_, err := s.client.Lightsail.DetachDisk(ctx, &lightsail.DetachDiskInput{
		DiskName: aws.String(diskName),
	})
	if err != nil {
		return fmt.Errorf("failed to detach disk %s: %w", diskName, err)
	}

	return nil
}

// DeleteDisk deletes a block storage disk.
func (s *LightsailService) DeleteDisk(ctx context.Context, diskName string) error {
	_, err := s.client.Lightsail.DeleteDisk(ctx, &lightsail.DeleteDiskInput{
		DiskName: aws.String(diskName),
	})
	if err != nil {
		return fmt.Errorf("failed to delete disk %s: %w", diskName, err)
	}

	return nil
}

// CreateInstanceSnapshot creates a snapshot of an instance.
func (s *LightsailService) CreateInstanceSnapshot(ctx context.Context, instanceName, snapshotName string) error {
	_, err := s.client.Lightsail.CreateInstanceSnapshot(ctx, &lightsail.CreateInstanceSnapshotInput{
		InstanceName:         aws.String(instanceName),
		InstanceSnapshotName: aws.String(snapshotName),
	})
	if err != nil {
		return fmt.Errorf("failed to create instance snapshot %s: %w", snapshotName, err)
	}

	return nil
}

// CreateInstanceFromSnapshot creates a new instance from a snapshot with specified bundle.
func (s *LightsailService) CreateInstanceFromSnapshot(ctx context.Context, newInstanceName, snapshotName, bundleID, availabilityZone string, tags map[string]string) (*types.Instance, error) {
	var lightsailTags []lightsailTypes.Tag
	for key, value := range tags {
		lightsailTags = append(lightsailTags, lightsailTypes.Tag{
			Key:   aws.String(key),
			Value: aws.String(value),
		})
	}

	_, err := s.client.Lightsail.CreateInstancesFromSnapshot(ctx, &lightsail.CreateInstancesFromSnapshotInput{
		InstanceNames:        []string{newInstanceName},
		InstanceSnapshotName: aws.String(snapshotName),
		BundleId:             aws.String(bundleID),
		AvailabilityZone:     aws.String(availabilityZone),
		Tags:                 lightsailTags,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create instance %s from snapshot %s: %w", newInstanceName, snapshotName, err)
	}

	return s.GetInstance(ctx, newInstanceName)
}

// DeleteInstanceSnapshot deletes an instance snapshot.
func (s *LightsailService) DeleteInstanceSnapshot(ctx context.Context, snapshotName string) error {
	_, err := s.client.Lightsail.DeleteInstanceSnapshot(ctx, &lightsail.DeleteInstanceSnapshotInput{
		InstanceSnapshotName: aws.String(snapshotName),
	})
	if err != nil {
		return fmt.Errorf("failed to delete instance snapshot %s: %w", snapshotName, err)
	}

	return nil
}

// GetInstanceSnapshot gets snapshot details.
func (s *LightsailService) GetInstanceSnapshot(ctx context.Context, snapshotName string) (*lightsailTypes.InstanceSnapshot, error) {
	output, err := s.client.Lightsail.GetInstanceSnapshot(ctx, &lightsail.GetInstanceSnapshotInput{
		InstanceSnapshotName: aws.String(snapshotName),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get instance snapshot %s: %w", snapshotName, err)
	}

	return output.InstanceSnapshot, nil
}