// Package types defines common data structures used across the application.
package types

import "time"

// Project represents a Lightsail for Research project configuration.
type Project struct {
	Name      string `json:"name" yaml:"name"`
	Blueprint string `json:"blueprint" yaml:"blueprint"`
	Bundle    string `json:"bundle" yaml:"bundle"`
	Region    string `json:"region" yaml:"region"`
	CreatedAt time.Time `json:"created_at" yaml:"created_at"`
}

// User represents an IAM user with associated Lightsail resources.
type User struct {
	Username     string    `json:"username" yaml:"username"`
	Project      string    `json:"project" yaml:"project"`
	InstanceARN  string    `json:"instance_arn" yaml:"instance_arn"`
	InstanceName string    `json:"instance_name" yaml:"instance_name"`
	Password     string    `json:"password,omitempty" yaml:"password,omitempty"`
	CreatedAt    time.Time `json:"created_at" yaml:"created_at"`
}

// Group represents an IAM group configuration.
type Group struct {
	Name        string   `json:"name" yaml:"name"`
	Policies    []string `json:"policies" yaml:"policies"`
	Description string   `json:"description" yaml:"description"`
	CreatedAt   time.Time `json:"created_at" yaml:"created_at"`
}

// Instance represents a Lightsail instance.
type Instance struct {
	Name         string            `json:"name" yaml:"name"`
	ARN          string            `json:"arn" yaml:"arn"`
	State        string            `json:"state" yaml:"state"`
	Blueprint    string            `json:"blueprint" yaml:"blueprint"`
	Bundle       string            `json:"bundle" yaml:"bundle"`
	Region       string            `json:"region" yaml:"region"`
	Tags         map[string]string `json:"tags" yaml:"tags"`
	CreatedAt    time.Time         `json:"created_at" yaml:"created_at"`
	PublicIP     string            `json:"public_ip,omitempty" yaml:"public_ip,omitempty"`
	PrivateIP    string            `json:"private_ip,omitempty" yaml:"private_ip,omitempty"`
}