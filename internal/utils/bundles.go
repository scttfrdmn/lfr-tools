// Package utils provides bundle management utilities for Lightsail for Research.
package utils

import (
	"fmt"
	"sort"
	"strings"
)

// BundleInfo represents bundle specifications.
type BundleInfo struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	RAM      float64 `json:"ram_gb"`
	VCPU     int     `json:"vcpu"`
	DiskGB   int     `json:"disk_gb"`
	IsGPU    bool    `json:"is_gpu"`
	SizeRank int     `json:"size_rank"`
}

// LightsailBundles contains the available Lightsail for Research bundles.
var LightsailBundles = []BundleInfo{
	{ID: "app_standard_xl_1_0", Name: "Standard XL", RAM: 8.0, VCPU: 4, DiskGB: 50, IsGPU: false, SizeRank: 1},
	{ID: "app_standard_2xl_1_0", Name: "Standard 2XL", RAM: 16.0, VCPU: 8, DiskGB: 50, IsGPU: false, SizeRank: 2},
	{ID: "app_standard_4xl_1_0", Name: "Standard 4XL", RAM: 32.0, VCPU: 16, DiskGB: 50, IsGPU: false, SizeRank: 3},
	{ID: "gpu_nvidia_xl_1_0", Name: "GPU XL", RAM: 16.0, VCPU: 4, DiskGB: 50, IsGPU: true, SizeRank: 1},
	{ID: "gpu_nvidia_2xl_1_0", Name: "GPU 2XL", RAM: 32.0, VCPU: 8, DiskGB: 50, IsGPU: true, SizeRank: 2},
	{ID: "gpu_nvidia_4xl_1_0", Name: "GPU 4XL", RAM: 64.0, VCPU: 16, DiskGB: 50, IsGPU: true, SizeRank: 3},
}

// GetBundleInfo returns bundle information by ID.
func GetBundleInfo(bundleID string) (*BundleInfo, error) {
	for _, bundle := range LightsailBundles {
		if bundle.ID == bundleID {
			return &bundle, nil
		}
	}
	return nil, fmt.Errorf("bundle not found: %s", bundleID)
}

// GetNextSizeBundle returns the next larger bundle in the same category (standard/GPU).
func GetNextSizeBundle(currentBundleID string) (*BundleInfo, error) {
	current, err := GetBundleInfo(currentBundleID)
	if err != nil {
		return nil, err
	}

	// Find bundles in the same category
	var candidates []BundleInfo
	for _, bundle := range LightsailBundles {
		if bundle.IsGPU == current.IsGPU && bundle.SizeRank > current.SizeRank {
			candidates = append(candidates, bundle)
		}
	}

	if len(candidates) == 0 {
		return nil, fmt.Errorf("no larger bundle available for %s", currentBundleID)
	}

	// Sort by size rank and return the smallest larger one
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].SizeRank < candidates[j].SizeRank
	})

	return &candidates[0], nil
}

// GetPreviousSizeBundle returns the next smaller bundle in the same category.
func GetPreviousSizeBundle(currentBundleID string) (*BundleInfo, error) {
	current, err := GetBundleInfo(currentBundleID)
	if err != nil {
		return nil, err
	}

	// Find bundles in the same category
	var candidates []BundleInfo
	for _, bundle := range LightsailBundles {
		if bundle.IsGPU == current.IsGPU && bundle.SizeRank < current.SizeRank {
			candidates = append(candidates, bundle)
		}
	}

	if len(candidates) == 0 {
		return nil, fmt.Errorf("no smaller bundle available for %s", currentBundleID)
	}

	// Sort by size rank and return the largest smaller one
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].SizeRank > candidates[j].SizeRank
	})

	return &candidates[0], nil
}

// GetEquivalentGPUBundle returns the equivalent GPU bundle for a standard bundle.
func GetEquivalentGPUBundle(standardBundleID string) (*BundleInfo, error) {
	current, err := GetBundleInfo(standardBundleID)
	if err != nil {
		return nil, err
	}

	if current.IsGPU {
		return nil, fmt.Errorf("bundle %s is already a GPU bundle", standardBundleID)
	}

	// Find GPU bundle with same or similar size rank
	for _, bundle := range LightsailBundles {
		if bundle.IsGPU && bundle.SizeRank == current.SizeRank {
			return &bundle, nil
		}
	}

	return nil, fmt.Errorf("no equivalent GPU bundle found for %s", standardBundleID)
}

// GetEquivalentStandardBundle returns the equivalent standard bundle for a GPU bundle.
func GetEquivalentStandardBundle(gpuBundleID string) (*BundleInfo, error) {
	current, err := GetBundleInfo(gpuBundleID)
	if err != nil {
		return nil, err
	}

	if !current.IsGPU {
		return nil, fmt.Errorf("bundle %s is not a GPU bundle", gpuBundleID)
	}

	// Find standard bundle with same or similar size rank
	for _, bundle := range LightsailBundles {
		if !bundle.IsGPU && bundle.SizeRank == current.SizeRank {
			return &bundle, nil
		}
	}

	return nil, fmt.Errorf("no equivalent standard bundle found for %s", gpuBundleID)
}

// ListBundlesByCategory returns bundles grouped by type.
func ListBundlesByCategory() (standard, gpu []BundleInfo) {
	for _, bundle := range LightsailBundles {
		if bundle.IsGPU {
			gpu = append(gpu, bundle)
		} else {
			standard = append(standard, bundle)
		}
	}

	// Sort by size rank
	sort.Slice(standard, func(i, j int) bool { return standard[i].SizeRank < standard[j].SizeRank })
	sort.Slice(gpu, func(i, j int) bool { return gpu[i].SizeRank < gpu[j].SizeRank })

	return standard, gpu
}

// FormatBundleComparison formats a bundle comparison for display.
func FormatBundleComparison(from, to *BundleInfo) string {
	var change string
	if to.RAM > from.RAM {
		change = "⬆️ LARGER"
	} else if to.RAM < from.RAM {
		change = "⬇️ SMALLER"
	} else {
		change = "↔️ SAME SIZE"
	}

	return fmt.Sprintf("%s: %s (%.1fGB RAM, %dv CPU) → %s (%.1fGB RAM, %dv CPU) %s",
		change, from.Name, from.RAM, from.VCPU, to.Name, to.RAM, to.VCPU,
		strings.Repeat(" ", 10))
}