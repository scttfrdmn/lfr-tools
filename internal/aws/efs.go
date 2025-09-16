package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2Types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/efs"
	efsTypes "github.com/aws/aws-sdk-go-v2/service/efs/types"
	"github.com/aws/aws-sdk-go-v2/service/lightsail"
)

// EFSService provides EFS operations and Lightsail integration.
type EFSService struct {
	client       *Client
	efsClient    *efs.Client
	ec2Client    *ec2.Client
}

// NewEFSService creates a new EFS service.
func NewEFSService(client *Client) *EFSService {
	return &EFSService{
		client:    client,
		efsClient: efs.NewFromConfig(client.Config),
		ec2Client: ec2.NewFromConfig(client.Config),
	}
}

// EFSFileSystem represents an EFS file system configuration.
type EFSFileSystem struct {
	ID               string            `json:"id"`
	Name             string            `json:"name"`
	State            string            `json:"state"`
	Region           string            `json:"region"`
	MountTargets     []EFSMountTarget  `json:"mount_targets"`
	Tags             map[string]string `json:"tags"`
	CreationTime     string            `json:"creation_time"`
	PerformanceMode  string            `json:"performance_mode"`
	ThroughputMode   string            `json:"throughput_mode"`
}

// EFSMountTarget represents an EFS mount target.
type EFSMountTarget struct {
	ID               string `json:"id"`
	IPAddress        string `json:"ip_address"`
	SubnetID         string `json:"subnet_id"`
	AvailabilityZone string `json:"availability_zone"`
	State            string `json:"state"`
}

// EnableVPCPeering enables VPC peering for Lightsail in the current region.
func (s *EFSService) EnableVPCPeering(ctx context.Context) error {
	_, err := s.client.Lightsail.PeerVpc(ctx, &lightsail.PeerVpcInput{})
	if err != nil {
		return fmt.Errorf("failed to enable VPC peering: %w", err)
	}

	return nil
}

// IsVPCPeered checks if VPC peering is enabled.
func (s *EFSService) IsVPCPeered(ctx context.Context) (bool, error) {
	output, err := s.client.Lightsail.IsVpcPeered(ctx, &lightsail.IsVpcPeeredInput{})
	if err != nil {
		return false, fmt.Errorf("failed to check VPC peering status: %w", err)
	}

	return aws.ToBool(output.IsPeered), nil
}

// GetDefaultVPC gets the default VPC for the current region.
func (s *EFSService) GetDefaultVPC(ctx context.Context) (string, error) {
	output, err := s.ec2Client.DescribeVpcs(ctx, &ec2.DescribeVpcsInput{
		Filters: []ec2Types.Filter{
			{
				Name:   aws.String("is-default"),
				Values: []string{"true"},
			},
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to describe default VPC: %w", err)
	}

	if len(output.Vpcs) == 0 {
		return "", fmt.Errorf("no default VPC found in region")
	}

	return aws.ToString(output.Vpcs[0].VpcId), nil
}

// CreateEFSFileSystem creates an EFS file system with proper configuration for Lightsail.
func (s *EFSService) CreateEFSFileSystem(ctx context.Context, name, project string) (*EFSFileSystem, error) {
	// Create the EFS file system
	createOutput, err := s.efsClient.CreateFileSystem(ctx, &efs.CreateFileSystemInput{
		CreationToken:   aws.String(fmt.Sprintf("lfr-%s-%s", project, name)),
		PerformanceMode: efsTypes.PerformanceModeGeneralPurpose,
		ThroughputMode:  efsTypes.ThroughputModeBursting,
		Tags: []efsTypes.Tag{
			{
				Key:   aws.String("Name"),
				Value: aws.String(name),
			},
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
		return nil, fmt.Errorf("failed to create EFS file system: %w", err)
	}

	return &EFSFileSystem{
		ID:              aws.ToString(createOutput.FileSystemId),
		Name:            name,
		State:           string(createOutput.LifeCycleState),
		Region:          s.client.GetRegion(),
		PerformanceMode: string(createOutput.PerformanceMode),
		ThroughputMode:  string(createOutput.ThroughputMode),
		CreationTime:    createOutput.CreationTime.String(),
	}, nil
}

// CreateMountTargets creates mount targets for EFS in all available subnets.
func (s *EFSService) CreateMountTargets(ctx context.Context, fileSystemID string) ([]EFSMountTarget, error) {
	// Get default VPC
	vpcID, err := s.GetDefaultVPC(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get default VPC: %w", err)
	}

	// Get subnets in the default VPC
	subnetsOutput, err := s.ec2Client.DescribeSubnets(ctx, &ec2.DescribeSubnetsInput{
		Filters: []ec2Types.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []string{vpcID},
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to describe subnets: %w", err)
	}

	// Create security group for EFS access
	sgID, err := s.createEFSSecurityGroup(ctx, vpcID)
	if err != nil {
		return nil, fmt.Errorf("failed to create security group: %w", err)
	}

	var mountTargets []EFSMountTarget

	// Create mount targets in each subnet
	for _, subnet := range subnetsOutput.Subnets {
		mtOutput, err := s.efsClient.CreateMountTarget(ctx, &efs.CreateMountTargetInput{
			FileSystemId:   aws.String(fileSystemID),
			SubnetId:       subnet.SubnetId,
			SecurityGroups: []string{sgID},
		})
		if err != nil {
			fmt.Printf("Warning: Failed to create mount target in subnet %s: %v\n", aws.ToString(subnet.SubnetId), err)
			continue
		}

		mountTarget := EFSMountTarget{
			ID:               aws.ToString(mtOutput.MountTargetId),
			IPAddress:        aws.ToString(mtOutput.IpAddress),
			SubnetID:         aws.ToString(mtOutput.SubnetId),
			AvailabilityZone: aws.ToString(subnet.AvailabilityZone),
			State:            string(mtOutput.LifeCycleState),
		}

		mountTargets = append(mountTargets, mountTarget)
	}

	return mountTargets, nil
}

// createEFSSecurityGroup creates a security group allowing NFS traffic from Lightsail.
func (s *EFSService) createEFSSecurityGroup(ctx context.Context, vpcID string) (string, error) {
	// Get Lightsail VPC CIDR (typically 172.26.0.0/16)
	lightsailCIDR := "172.26.0.0/16" // Lightsail's VPC CIDR range

	sgOutput, err := s.ec2Client.CreateSecurityGroup(ctx, &ec2.CreateSecurityGroupInput{
		GroupName:   aws.String("lfr-efs-access"),
		Description: aws.String("Allow NFS access from Lightsail instances"),
		VpcId:       aws.String(vpcID),
	})
	if err != nil {
		return "", fmt.Errorf("failed to create security group: %w", err)
	}

	sgID := aws.ToString(sgOutput.GroupId)

	// Add inbound rule for NFS (port 2049) from Lightsail VPC
	_, err = s.ec2Client.AuthorizeSecurityGroupIngress(ctx, &ec2.AuthorizeSecurityGroupIngressInput{
		GroupId: aws.String(sgID),
		IpPermissions: []ec2Types.IpPermission{
			{
				IpProtocol: aws.String("tcp"),
				FromPort:   aws.Int32(2049),
				ToPort:     aws.Int32(2049),
				IpRanges: []ec2Types.IpRange{
					{
						CidrIp:      aws.String(lightsailCIDR),
						Description: aws.String("NFS access from Lightsail instances"),
					},
				},
			},
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to authorize security group ingress: %w", err)
	}

	return sgID, nil
}

// ListEFSFileSystems lists EFS file systems.
func (s *EFSService) ListEFSFileSystems(ctx context.Context, project string) ([]*EFSFileSystem, error) {
	output, err := s.efsClient.DescribeFileSystems(ctx, &efs.DescribeFileSystemsInput{})
	if err != nil {
		return nil, fmt.Errorf("failed to describe EFS file systems: %w", err)
	}

	var fileSystems []*EFSFileSystem

	for _, fs := range output.FileSystems {
		// Get tags
		tagsOutput, err := s.efsClient.DescribeTags(ctx, &efs.DescribeTagsInput{
			FileSystemId: fs.FileSystemId,
		})
		if err != nil {
			continue // Skip if can't get tags
		}

		tags := make(map[string]string)
		var name string
		var fsProject string
		for _, tag := range tagsOutput.Tags {
			tagKey := aws.ToString(tag.Key)
			tagValue := aws.ToString(tag.Value)
			tags[tagKey] = tagValue

			if tagKey == "Name" {
				name = tagValue
			}
			if tagKey == "Project" {
				fsProject = tagValue
			}
		}

		// Filter by project if specified
		if project != "" && fsProject != project {
			continue
		}

		fileSystem := &EFSFileSystem{
			ID:              aws.ToString(fs.FileSystemId),
			Name:            name,
			State:           string(fs.LifeCycleState),
			Region:          s.client.GetRegion(),
			Tags:            tags,
			CreationTime:    fs.CreationTime.String(),
			PerformanceMode: string(fs.PerformanceMode),
			ThroughputMode:  string(fs.ThroughputMode),
		}

		// Get mount targets
		mtOutput, err := s.efsClient.DescribeMountTargets(ctx, &efs.DescribeMountTargetsInput{
			FileSystemId: fs.FileSystemId,
		})
		if err == nil {
			for _, mt := range mtOutput.MountTargets {
				mountTarget := EFSMountTarget{
					ID:        aws.ToString(mt.MountTargetId),
					IPAddress: aws.ToString(mt.IpAddress),
					SubnetID:  aws.ToString(mt.SubnetId),
					State:     string(mt.LifeCycleState),
				}
				fileSystem.MountTargets = append(fileSystem.MountTargets, mountTarget)
			}
		}

		fileSystems = append(fileSystems, fileSystem)
	}

	return fileSystems, nil
}