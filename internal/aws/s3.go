package aws

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3Types "github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// S3Service provides S3 operations for student status updates.
type S3Service struct {
	client *Client
	s3     *s3.Client
}

// NewS3Service creates a new S3 service.
func NewS3Service(client *Client) *S3Service {
	return &S3Service{
		client: client,
		s3:     s3.NewFromConfig(client.Config),
	}
}

// StudentStatus represents the status information for a student's instance.
type StudentStatus struct {
	State           string    `json:"state"`
	PublicIP        string    `json:"public_ip,omitempty"`
	LastUpdated     time.Time `json:"last_updated"`
	StartRequested  bool      `json:"start_requested"`
	RequestedAt     time.Time `json:"requested_at,omitempty"`
	RequestedBy     string    `json:"requested_by,omitempty"`
	BudgetRemaining float64   `json:"budget_remaining,omitempty"`
	AccessExpires   time.Time `json:"access_expires,omitempty"`
}

// UpdateStudentStatus updates a student's status in S3.
func (s *S3Service) UpdateStudentStatus(ctx context.Context, bucket, project, username string, status *StudentStatus) error {
	key := fmt.Sprintf("%s/%s/status.json", project, username)

	// Add timestamp
	status.LastUpdated = time.Now()

	// Marshal to JSON
	data, err := json.MarshalIndent(status, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal status: %w", err)
	}

	// Upload to S3
	_, err = s.s3.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String("application/json"),
		ACL:         s3Types.ObjectCannedACLPublicRead, // Public read for student access
	})
	if err != nil {
		return fmt.Errorf("failed to update student status in S3: %w", err)
	}

	return nil
}

// GetStudentStatus retrieves a student's status from S3.
func (s *S3Service) GetStudentStatus(ctx context.Context, bucket, project, username string) (*StudentStatus, error) {
	key := fmt.Sprintf("%s/%s/status.json", project, username)

	output, err := s.s3.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get student status from S3: %w", err)
	}
	defer output.Body.Close()

	var status StudentStatus
	if err := json.NewDecoder(output.Body).Decode(&status); err != nil {
		return nil, fmt.Errorf("failed to decode student status: %w", err)
	}

	return &status, nil
}

// CheckStartRequests checks for pending start requests in S3.
func (s *S3Service) CheckStartRequests(ctx context.Context, bucket, project string) (map[string]*StudentStartRequest, error) {
	// List all start request files for the project
	prefix := fmt.Sprintf("%s/", project)

	output, err := s.s3.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list S3 objects: %w", err)
	}

	requests := make(map[string]*StudentStartRequest)

	for _, obj := range output.Contents {
		key := aws.ToString(obj.Key)
		if !bytes.HasSuffix([]byte(key), []byte("/start-request.json")) {
			continue
		}

		// Extract username from key
		parts := bytes.Split([]byte(key), []byte("/"))
		if len(parts) < 2 {
			continue
		}
		username := string(parts[1])

		// Get the start request
		reqOutput, err := s.s3.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
		if err != nil {
			continue // Skip failed requests
		}

		var request StudentStartRequest
		if err := json.NewDecoder(reqOutput.Body).Decode(&request); err != nil {
			reqOutput.Body.Close()
			continue
		}
		reqOutput.Body.Close()

		requests[username] = &request
	}

	return requests, nil
}

// StudentStartRequest represents a student's request to start their instance.
type StudentStartRequest struct {
	Username    string    `json:"username"`
	StudentID   string    `json:"student_id"`
	Token       string    `json:"token"`
	RequestedAt time.Time `json:"requested_at"`
	MachineHash string    `json:"machine_hash"`
	RequestIP   string    `json:"request_ip,omitempty"`
}

// SubmitStartRequest submits a start request for a student.
func (s *S3Service) SubmitStartRequest(ctx context.Context, bucket, project string, request *StudentStartRequest) error {
	key := fmt.Sprintf("%s/%s/start-request.json", project, request.Username)

	// Add timestamp
	request.RequestedAt = time.Now()

	// Marshal to JSON
	data, err := json.MarshalIndent(request, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal start request: %w", err)
	}

	// Upload to S3 with public write permission
	_, err = s.s3.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String("application/json"),
	})
	if err != nil {
		return fmt.Errorf("failed to submit start request to S3: %w", err)
	}

	return nil
}

// DeleteStartRequest removes a processed start request.
func (s *S3Service) DeleteStartRequest(ctx context.Context, bucket, project, username string) error {
	key := fmt.Sprintf("%s/%s/start-request.json", project, username)

	_, err := s.s3.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete start request: %w", err)
	}

	return nil
}

// CreateStatusBucket creates an S3 bucket for student status updates.
func (s *S3Service) CreateStatusBucket(ctx context.Context, bucketName, project string) error {
	// Create bucket
	_, err := s.s3.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
		CreateBucketConfiguration: &s3Types.CreateBucketConfiguration{
			LocationConstraint: s3Types.BucketLocationConstraint(s.client.GetRegion()),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create S3 bucket: %w", err)
	}

	// Configure bucket for public read
	policy := fmt.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Principal": "*",
				"Action": "s3:GetObject",
				"Resource": "arn:aws:s3:::%s/%s/*/status.json"
			},
			{
				"Effect": "Allow",
				"Principal": "*",
				"Action": "s3:PutObject",
				"Resource": "arn:aws:s3:::%s/%s/*/start-request.json"
			}
		]
	}`, bucketName, project, bucketName, project)

	_, err = s.s3.PutBucketPolicy(ctx, &s3.PutBucketPolicyInput{
		Bucket: aws.String(bucketName),
		Policy: aws.String(policy),
	})
	if err != nil {
		return fmt.Errorf("failed to set bucket policy: %w", err)
	}

	return nil
}