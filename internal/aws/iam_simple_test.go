package aws

import (
	"testing"

	"github.com/scttfrdmn/lfr-tools/internal/testutils"
)

func TestNewIAMService(t *testing.T) {
	fixture := testutils.NewTestFixture(t)

	client := &Client{}
	service := NewIAMService(client)

	if service == nil {
		fixture.AssertError(nil) // This will fail the test
	}

	if service.client != client {
		t.Error("expected service.client to be the same as input client")
	}
}

func TestLightsailService(t *testing.T) {
	fixture := testutils.NewTestFixture(t)

	client := &Client{}
	service := NewLightsailService(client)

	if service == nil {
		fixture.AssertError(nil) // This will fail the test
	}

	if service.client != client {
		t.Error("expected service.client to be the same as input client")
	}
}