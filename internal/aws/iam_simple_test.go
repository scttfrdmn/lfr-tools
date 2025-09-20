package aws

import (
	"testing"
)

func TestNewIAMService(t *testing.T) {
	client := &Client{}
	service := NewIAMService(client)

	if service == nil {
		t.Error("expected non-nil service")
		return
	}

	if service.client != client {
		t.Error("expected service.client to be the same as input client")
	}
}

func TestLightsailService(t *testing.T) {
	client := &Client{}
	service := NewLightsailService(client)

	if service == nil {
		t.Error("expected non-nil service")
		return
	}

	if service.client != client {
		t.Error("expected service.client to be the same as input client")
	}
}