// Package aws provides AWS SDK integration and common AWS operations.
package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/lightsail"
)

// Client wraps AWS service clients.
type Client struct {
	IAM       *iam.Client
	Lightsail *lightsail.Client
	Config    aws.Config
}

// Options for creating a new AWS client.
type Options struct {
	Region  string
	Profile string
}

// NewClient creates a new AWS client with the specified options.
func NewClient(ctx context.Context, opts Options) (*Client, error) {
	var configOpts []func(*config.LoadOptions) error

	if opts.Region != "" {
		configOpts = append(configOpts, config.WithRegion(opts.Region))
	}

	if opts.Profile != "" {
		configOpts = append(configOpts, config.WithSharedConfigProfile(opts.Profile))
	}

	cfg, err := config.LoadDefaultConfig(ctx, configOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	return &Client{
		IAM:       iam.NewFromConfig(cfg),
		Lightsail: lightsail.NewFromConfig(cfg),
		Config:    cfg,
	}, nil
}

// GetRegion returns the configured AWS region.
func (c *Client) GetRegion() string {
	return c.Config.Region
}