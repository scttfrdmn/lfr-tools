// Package config handles configuration management for lfr-tools.
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config represents the application configuration.
type Config struct {
	AWS      AWSConfig      `mapstructure:"aws" json:"aws" yaml:"aws"`
	Defaults DefaultsConfig `mapstructure:"defaults" json:"defaults" yaml:"defaults"`
	SSH      SSHConfig      `mapstructure:"ssh" json:"ssh" yaml:"ssh"`
	Debug    bool           `mapstructure:"debug" json:"debug" yaml:"debug"`
}

// AWSConfig holds AWS-related configuration.
type AWSConfig struct {
	Profile string `mapstructure:"profile" json:"profile" yaml:"profile"`
	Region  string `mapstructure:"region" json:"region" yaml:"region"`
}

// DefaultsConfig holds default values for commands.
type DefaultsConfig struct {
	Blueprint     string `mapstructure:"blueprint" json:"blueprint" yaml:"blueprint"`
	Bundle        string `mapstructure:"bundle" json:"bundle" yaml:"bundle"`
	IdleThreshold int    `mapstructure:"idle_threshold" json:"idle_threshold" yaml:"idle_threshold"`
}

// SSHConfig holds SSH-related configuration.
type SSHConfig struct {
	KeyPath    string `mapstructure:"key_path" json:"key_path" yaml:"key_path"`
	ConfigPath string `mapstructure:"config_path" json:"config_path" yaml:"config_path"`
}

// Load reads and parses the configuration from file and environment variables.
func Load() (*Config, error) {
	config := &Config{
		AWS: AWSConfig{
			Profile: "default",
			Region:  "us-east-1",
		},
		Defaults: DefaultsConfig{
			Blueprint:     "ubuntu_22_04",
			Bundle:        "nano_2_0",
			IdleThreshold: 120,
		},
		SSH: SSHConfig{
			KeyPath:    filepath.Join(mustGetHomeDir(), ".ssh", "lfr-tools"),
			ConfigPath: filepath.Join(mustGetHomeDir(), ".ssh", "config.d", "lfr-tools"),
		},
		Debug: false,
	}

	// Unmarshal the configuration
	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Expand tilde in paths
	config.SSH.KeyPath = expandPath(config.SSH.KeyPath)
	config.SSH.ConfigPath = expandPath(config.SSH.ConfigPath)

	return config, nil
}

// expandPath expands ~ to the user's home directory.
func expandPath(path string) string {
	if len(path) == 0 || path[0] != '~' {
		return path
	}

	homeDir := mustGetHomeDir()
	if len(path) == 1 {
		return homeDir
	}

	return filepath.Join(homeDir, path[1:])
}

// mustGetHomeDir returns the user's home directory or panics.
func mustGetHomeDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Sprintf("failed to get user home directory: %v", err))
	}
	return homeDir
}