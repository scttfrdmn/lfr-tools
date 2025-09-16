package aws

import (
	"testing"

	"github.com/scttfrdmn/lfr-tools/internal/testutils"
)

func TestNewClientOptions(t *testing.T) {
	opts := Options{
		Region:  "us-east-1",
		Profile: "default",
	}

	if opts.Region != "us-east-1" {
		t.Errorf("expected region 'us-east-1', got %s", opts.Region)
	}

	if opts.Profile != "default" {
		t.Errorf("expected profile 'default', got %s", opts.Profile)
	}
}

func TestGetRegion(t *testing.T) {
	client := &Client{
		Config: testutils.MockAWSConfig(),
	}

	region := client.GetRegion()
	if region != "us-east-1" {
		t.Errorf("expected region 'us-east-1', got %s", region)
	}
}