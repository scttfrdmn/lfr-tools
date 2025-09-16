// Package utils provides common utility functions including AWS operation waiters.
package utils

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// WaitConfig configures waiting behavior for AWS operations.
type WaitConfig struct {
	MaxDuration    time.Duration
	CheckInterval  time.Duration
	Operation      string
	ResourceName   string
	TargetState    string
	CurrentStateFn func() (string, error)
}

// DefaultWaitConfig returns sensible defaults for waiting operations.
func DefaultWaitConfig(operation, resourceName string) *WaitConfig {
	return &WaitConfig{
		MaxDuration:   5 * time.Minute,
		CheckInterval: 10 * time.Second,
		Operation:     operation,
		ResourceName:  resourceName,
	}
}

// WaitForState waits for a resource to reach a target state with progress updates and visual indicator.
func WaitForState(ctx context.Context, config *WaitConfig) error {
	fmt.Printf("⏳ Waiting for %s %s to reach state '%s'...\n",
		config.Operation, config.ResourceName, config.TargetState)

	start := time.Now()
	ticker := time.NewTicker(config.CheckInterval)
	defer ticker.Stop()

	// Spinner for visual feedback
	spinnerChars := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	spinnerTicker := time.NewTicker(100 * time.Millisecond)
	defer spinnerTicker.Stop()

	spinnerIndex := 0
	timeout := time.After(config.MaxDuration)

	// Start with initial state check
	currentState, err := config.CurrentStateFn()
	if err != nil {
		return fmt.Errorf("error checking initial state: %w", err)
	}

	if currentState == config.TargetState {
		fmt.Printf("✅ %s %s is already in state '%s'\n",
			config.Operation, config.ResourceName, config.TargetState)
		return nil
	}

	lastState := currentState
	fmt.Printf("  %s %s: %s → %s\n", spinnerChars[0], config.ResourceName, currentState, config.TargetState)

	for {
		select {
		case <-ctx.Done():
			fmt.Print("\r") // Clear spinner
			return ctx.Err()

		case <-timeout:
			fmt.Print("\r") // Clear spinner
			elapsed := time.Since(start)
			return fmt.Errorf("timeout waiting for %s %s after %v",
				config.Operation, config.ResourceName, elapsed.Round(time.Second))

		case <-spinnerTicker.C:
			// Update spinner
			elapsed := time.Since(start)
			fmt.Printf("\r  %s %s: %s → %s (elapsed: %v)",
				spinnerChars[spinnerIndex], config.ResourceName, lastState, config.TargetState, elapsed.Round(time.Second))
			spinnerIndex = (spinnerIndex + 1) % len(spinnerChars)

		case <-ticker.C:
			currentState, err := config.CurrentStateFn()
			if err != nil {
				fmt.Print("\r") // Clear spinner
				return fmt.Errorf("error checking state: %w", err)
			}

			if currentState != lastState {
				fmt.Print("\r") // Clear spinner
				elapsed := time.Since(start)
				fmt.Printf("  🔄 %s: %s → %s (elapsed: %v)\n",
					config.ResourceName, lastState, currentState, elapsed.Round(time.Second))
				lastState = currentState
			}

			if currentState == config.TargetState {
				fmt.Print("\r") // Clear spinner
				elapsed := time.Since(start)
				fmt.Printf("✅ %s %s reached state '%s' after %v\n",
					config.Operation, config.ResourceName, config.TargetState, elapsed.Round(time.Second))
				return nil
			}

			// Check for error states
			if strings.Contains(strings.ToLower(currentState), "error") ||
			   strings.Contains(strings.ToLower(currentState), "failed") {
				fmt.Print("\r") // Clear spinner
				return fmt.Errorf("%s %s entered error state: %s",
					config.Operation, config.ResourceName, currentState)
			}
		}
	}
}

// WaitForInstanceState waits for a Lightsail instance to reach target state.
func WaitForInstanceState(ctx context.Context, instanceName, targetState string, checkStateFn func() (string, error)) error {
	config := DefaultWaitConfig("instance", instanceName)
	config.TargetState = targetState
	config.CurrentStateFn = checkStateFn

	// Instance operations can take longer
	config.MaxDuration = 10 * time.Minute

	return WaitForState(ctx, config)
}

// WaitForEFSState waits for an EFS file system to reach target state.
func WaitForEFSState(ctx context.Context, filesystemID, targetState string, checkStateFn func() (string, error)) error {
	config := DefaultWaitConfig("EFS", filesystemID)
	config.TargetState = targetState
	config.CurrentStateFn = checkStateFn

	// EFS operations are usually faster
	config.MaxDuration = 3 * time.Minute
	config.CheckInterval = 5 * time.Second

	return WaitForState(ctx, config)
}

// WaitForMountTargetState waits for EFS mount targets to become available.
func WaitForMountTargetState(ctx context.Context, mountTargetID, targetState string, checkStateFn func() (string, error)) error {
	config := DefaultWaitConfig("mount target", mountTargetID)
	config.TargetState = targetState
	config.CurrentStateFn = checkStateFn
	config.MaxDuration = 3 * time.Minute
	config.CheckInterval = 5 * time.Second

	return WaitForState(ctx, config)
}

// WaitForDiskState waits for a Lightsail disk to reach target state.
func WaitForDiskState(ctx context.Context, diskName, targetState string, checkStateFn func() (string, error)) error {
	config := DefaultWaitConfig("volume", diskName)
	config.TargetState = targetState
	config.CurrentStateFn = checkStateFn

	// Disk operations are usually fast
	config.MaxDuration = 5 * time.Minute
	config.CheckInterval = 5 * time.Second

	return WaitForState(ctx, config)
}

// WaitForSnapshotState waits for a Lightsail snapshot to reach target state.
func WaitForSnapshotState(ctx context.Context, snapshotName, targetState string, checkStateFn func() (string, error)) error {
	config := DefaultWaitConfig("snapshot", snapshotName)
	config.TargetState = targetState
	config.CurrentStateFn = checkStateFn

	// Snapshots can take longer
	config.MaxDuration = 15 * time.Minute
	config.CheckInterval = 10 * time.Second

	return WaitForState(ctx, config)
}