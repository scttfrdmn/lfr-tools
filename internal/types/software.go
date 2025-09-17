// Package types defines software pack data structures.
package types

// SoftwarePack represents a deployable software package.
type SoftwarePack struct {
	ID          string            `json:"id" yaml:"id"`
	Name        string            `json:"name" yaml:"name"`
	Description string            `json:"description" yaml:"description"`
	Category    string            `json:"category" yaml:"category"`
	Type        PackType          `json:"type" yaml:"type"`
	Version     string            `json:"version" yaml:"version"`
	Dependencies []string         `json:"dependencies" yaml:"dependencies"`
	Packages    []Package         `json:"packages" yaml:"packages"`
	Scripts     []Script          `json:"scripts" yaml:"scripts"`
	Environment map[string]string `json:"environment" yaml:"environment"`
	Tags        []string          `json:"tags" yaml:"tags"`
	Supported   []string          `json:"supported_platforms" yaml:"supported_platforms"`
}

// PackType defines the type of software pack.
type PackType string

const (
	PackTypeAPT       PackType = "apt"
	PackTypeContainer PackType = "container"
	PackTypeScript    PackType = "script"
	PackTypeMixed     PackType = "mixed"
)

// Package represents an individual software package.
type Package struct {
	Name     string   `json:"name" yaml:"name"`
	Version  string   `json:"version,omitempty" yaml:"version,omitempty"`
	Source   string   `json:"source" yaml:"source"` // apt, snap, pip, npm, etc.
	Options  []string `json:"options,omitempty" yaml:"options,omitempty"`
	PostInstall []string `json:"post_install,omitempty" yaml:"post_install,omitempty"`
}

// Script represents a setup or configuration script.
type Script struct {
	Name        string            `json:"name" yaml:"name"`
	Description string            `json:"description" yaml:"description"`
	Content     string            `json:"content" yaml:"content"`
	Type        string            `json:"type" yaml:"type"` // bash, python, etc.
	RunAs       string            `json:"run_as,omitempty" yaml:"run_as,omitempty"`
	Environment map[string]string `json:"environment,omitempty" yaml:"environment,omitempty"`
}

// InstallResult represents the result of a software pack installation.
type InstallResult struct {
	PackID      string    `json:"pack_id"`
	Success     bool      `json:"success"`
	Message     string    `json:"message"`
	Duration    string    `json:"duration"`
	InstalledAt string    `json:"installed_at"`
	Packages    []string  `json:"packages_installed"`
	Errors      []string  `json:"errors,omitempty"`
}