package utils

import (
	"testing"
)

func TestGetBundleInfo(t *testing.T) {
	tests := []struct {
		bundleID string
		expected bool
		isGPU    bool
	}{
		{"app_standard_xl_1_0", true, false},
		{"app_standard_2xl_1_0", true, false},
		{"gpu_nvidia_xl_1_0", true, true},
		{"gpu_nvidia_2xl_1_0", true, true},
		{"nonexistent_bundle", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.bundleID, func(t *testing.T) {
			bundle, err := GetBundleInfo(tt.bundleID)

			if tt.expected {
				if err != nil {
					t.Errorf("expected to find bundle %s, got error: %v", tt.bundleID, err)
					return
				}

				if bundle == nil {
					t.Errorf("expected non-nil bundle for %s", tt.bundleID)
					return
				}

				if bundle.ID != tt.bundleID {
					t.Errorf("expected bundle ID %s, got %s", tt.bundleID, bundle.ID)
				}

				if bundle.IsGPU != tt.isGPU {
					t.Errorf("expected IsGPU=%v for %s, got %v", tt.isGPU, tt.bundleID, bundle.IsGPU)
				}
			} else {
				if err == nil {
					t.Errorf("expected error for nonexistent bundle %s", tt.bundleID)
				}
			}
		})
	}
}

func TestGetNextSizeBundle(t *testing.T) {
	tests := []struct {
		current  string
		expected string
		hasNext  bool
	}{
		{"app_standard_xl_1_0", "app_standard_2xl_1_0", true},
		{"app_standard_2xl_1_0", "app_standard_4xl_1_0", true},
		{"app_standard_4xl_1_0", "", false}, // No larger bundle
		{"gpu_nvidia_xl_1_0", "gpu_nvidia_2xl_1_0", true},
		{"gpu_nvidia_4xl_1_0", "", false}, // No larger bundle
	}

	for _, tt := range tests {
		t.Run(tt.current+"_to_"+tt.expected, func(t *testing.T) {
			next, err := GetNextSizeBundle(tt.current)

			if tt.hasNext {
				if err != nil {
					t.Errorf("expected to find next bundle for %s, got error: %v", tt.current, err)
					return
				}

				if next == nil {
					t.Errorf("expected non-nil next bundle for %s", tt.current)
					return
				}

				if next.ID != tt.expected {
					t.Errorf("expected next bundle %s, got %s", tt.expected, next.ID)
				}
			} else {
				if err == nil {
					t.Errorf("expected error for bundle %s (no larger bundle available)", tt.current)
				}
			}
		})
	}
}

func TestGetEquivalentGPUBundle(t *testing.T) {
	tests := []struct {
		standard string
		gpu      string
		hasGPU   bool
	}{
		{"app_standard_xl_1_0", "gpu_nvidia_xl_1_0", true},
		{"app_standard_2xl_1_0", "gpu_nvidia_2xl_1_0", true},
		{"app_standard_4xl_1_0", "gpu_nvidia_4xl_1_0", true},
	}

	for _, tt := range tests {
		t.Run(tt.standard+"_to_"+tt.gpu, func(t *testing.T) {
			gpu, err := GetEquivalentGPUBundle(tt.standard)

			if tt.hasGPU {
				if err != nil {
					t.Errorf("expected to find GPU equivalent for %s, got error: %v", tt.standard, err)
					return
				}

				if gpu == nil {
					t.Errorf("expected non-nil GPU bundle for %s", tt.standard)
					return
				}

				if gpu.ID != tt.gpu {
					t.Errorf("expected GPU bundle %s, got %s", tt.gpu, gpu.ID)
				}

				if !gpu.IsGPU {
					t.Errorf("expected GPU bundle to have IsGPU=true")
				}
			} else {
				if err == nil {
					t.Errorf("expected error for standard bundle %s", tt.standard)
				}
			}
		})
	}
}

func TestGetEquivalentStandardBundle(t *testing.T) {
	tests := []struct {
		gpu      string
		standard string
		hasStd   bool
	}{
		{"gpu_nvidia_xl_1_0", "app_standard_xl_1_0", true},
		{"gpu_nvidia_2xl_1_0", "app_standard_2xl_1_0", true},
		{"gpu_nvidia_4xl_1_0", "app_standard_4xl_1_0", true},
	}

	for _, tt := range tests {
		t.Run(tt.gpu+"_to_"+tt.standard, func(t *testing.T) {
			std, err := GetEquivalentStandardBundle(tt.gpu)

			if tt.hasStd {
				if err != nil {
					t.Errorf("expected to find standard equivalent for %s, got error: %v", tt.gpu, err)
					return
				}

				if std == nil {
					t.Errorf("expected non-nil standard bundle for %s", tt.gpu)
					return
				}

				if std.ID != tt.standard {
					t.Errorf("expected standard bundle %s, got %s", tt.standard, std.ID)
				}

				if std.IsGPU {
					t.Errorf("expected standard bundle to have IsGPU=false")
				}
			} else {
				if err == nil {
					t.Errorf("expected error for GPU bundle %s", tt.gpu)
				}
			}
		})
	}
}

func TestListBundlesByCategory(t *testing.T) {
	standard, gpu := ListBundlesByCategory()

	if len(standard) == 0 {
		t.Error("expected at least one standard bundle")
	}

	if len(gpu) == 0 {
		t.Error("expected at least one GPU bundle")
	}

	// Verify all standard bundles are not GPU
	for _, bundle := range standard {
		if bundle.IsGPU {
			t.Errorf("standard bundle %s has IsGPU=true", bundle.ID)
		}
	}

	// Verify all GPU bundles are GPU
	for _, bundle := range gpu {
		if !bundle.IsGPU {
			t.Errorf("GPU bundle %s has IsGPU=false", bundle.ID)
		}
	}

	// Verify size ranking is consistent
	for i := 1; i < len(standard); i++ {
		if standard[i].SizeRank <= standard[i-1].SizeRank {
			t.Errorf("standard bundles not properly sorted by size rank")
		}
	}

	for i := 1; i < len(gpu); i++ {
		if gpu[i].SizeRank <= gpu[i-1].SizeRank {
			t.Errorf("GPU bundles not properly sorted by size rank")
		}
	}
}

func TestFormatBundleComparison(t *testing.T) {
	xl, err := GetBundleInfo("app_standard_xl_1_0")
	if err != nil {
		t.Fatalf("failed to get XL bundle: %v", err)
	}

	xxl, err := GetBundleInfo("app_standard_2xl_1_0")
	if err != nil {
		t.Fatalf("failed to get 2XL bundle: %v", err)
	}

	// Test upgrade comparison
	comparison := FormatBundleComparison(xl, xxl)
	if comparison == "" {
		t.Error("expected non-empty comparison string")
	}

	if !contains(comparison, "LARGER") {
		t.Error("expected comparison to indicate size increase")
	}

	if !contains(comparison, xl.Name) {
		t.Error("expected comparison to include source bundle name")
	}

	if !contains(comparison, xxl.Name) {
		t.Error("expected comparison to include target bundle name")
	}
}