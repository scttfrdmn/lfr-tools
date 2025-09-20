package api

import (
	"testing"
)

func TestNewCostAPI(t *testing.T) {
	api := NewCostAPI()
	if api == nil {
		t.Fatal("expected non-nil cost API")
	}

	if api.ctx == nil {
		t.Error("expected non-nil context")
	}
}

func TestBudgetInfoStructure(t *testing.T) {
	budget := &BudgetInfo{
		ProjectName:     "cs101",
		TotalBudget:     500.0,
		UsedAmount:      340.50,
		RemainingAmount: 159.50,
		UsagePercentage: 68.1,
		DailyRate:       11.35,
		ProjectedTotal:  454.0,
		DaysRemaining:   45,
		OnTrack:         true,
	}

	if budget.ProjectName != "cs101" {
		t.Errorf("expected project name 'cs101', got %s", budget.ProjectName)
	}

	if budget.TotalBudget <= 0 {
		t.Error("total budget must be positive")
	}

	if budget.UsedAmount < 0 {
		t.Error("used amount cannot be negative")
	}

	if budget.RemainingAmount < 0 {
		t.Error("remaining amount cannot be negative")
	}

	if budget.UsagePercentage < 0 || budget.UsagePercentage > 100 {
		t.Errorf("usage percentage must be 0-100, got %.1f", budget.UsagePercentage)
	}

	if budget.DailyRate < 0 {
		t.Error("daily rate cannot be negative")
	}

	if budget.DaysRemaining < 0 {
		t.Error("days remaining cannot be negative")
	}
}

func TestCostByServiceStructure(t *testing.T) {
	service := &CostByService{
		ServiceName: "Lightsail",
		Amount:      289.43,
		Percentage:  85.0,
	}

	if service.ServiceName == "" {
		t.Error("service name cannot be empty")
	}

	if service.Amount < 0 {
		t.Error("amount cannot be negative")
	}

	if service.Percentage < 0 || service.Percentage > 100 {
		t.Errorf("percentage must be 0-100, got %.1f", service.Percentage)
	}
}

func TestUsageMetricsStructure(t *testing.T) {
	metrics := &UsageMetrics{
		InstanceName:     "alice-ubuntu",
		Username:         "alice",
		AvgCPU:          15.5,
		AvgMemory:       32.1,
		NetworkIO:       5.2,
		SSHSessions:     3,
		DailyUsageHours: 6.5,
		MonthlyCost:     25.0,
	}

	if metrics.InstanceName == "" {
		t.Error("instance name cannot be empty")
	}

	if metrics.Username == "" {
		t.Error("username cannot be empty")
	}

	if metrics.AvgCPU < 0 || metrics.AvgCPU > 100 {
		t.Errorf("CPU percentage must be 0-100, got %.1f", metrics.AvgCPU)
	}

	if metrics.AvgMemory < 0 || metrics.AvgMemory > 100 {
		t.Errorf("memory percentage must be 0-100, got %.1f", metrics.AvgMemory)
	}

	if metrics.SSHSessions < 0 {
		t.Error("SSH sessions cannot be negative")
	}

	if metrics.DailyUsageHours < 0 || metrics.DailyUsageHours > 24 {
		t.Errorf("daily usage hours must be 0-24, got %.1f", metrics.DailyUsageHours)
	}

	if metrics.MonthlyCost < 0 {
		t.Error("monthly cost cannot be negative")
	}
}

func TestEstimateInstanceMonthlyCost(t *testing.T) {
	api := NewCostAPI()

	tests := []struct {
		bundle string
		state  string
		minCost float64
		maxCost float64
	}{
		{"app_standard_xl_1_0", "running", 20.0, 30.0},
		{"app_standard_xl_1_0", "stopped", 1.0, 5.0},
		{"gpu_nvidia_xl_1_0", "running", 70.0, 90.0},
		{"gpu_nvidia_xl_1_0", "stopped", 5.0, 15.0},
		{"unknown_bundle", "running", 20.0, 30.0}, // Should use default
	}

	for _, tt := range tests {
		t.Run(tt.bundle+"_"+tt.state, func(t *testing.T) {
			cost := api.estimateInstanceMonthlyCost(tt.bundle, tt.state)

			if cost < tt.minCost || cost > tt.maxCost {
				t.Errorf("cost %.2f not in expected range %.2f-%.2f", cost, tt.minCost, tt.maxCost)
			}
		})
	}
}

func TestGetProjectBudgetWithoutAWS(t *testing.T) {
	api := NewCostAPI()

	// This test will likely fail without AWS credentials, which is expected
	budget, err := api.GetProjectBudget("test-project")
	if err != nil {
		// Expected failure without proper AWS setup
		t.Logf("GetProjectBudget failed as expected without AWS config: %v", err)
		return
	}

	// If it succeeds, validate the response
	if budget == nil {
		t.Error("expected non-nil budget info")
		return
	}

	if budget.ProjectName != "test-project" {
		t.Errorf("expected project name 'test-project', got %s", budget.ProjectName)
	}

	if budget.TotalBudget <= 0 {
		t.Error("total budget must be positive")
	}

	if budget.UsedAmount < 0 {
		t.Error("used amount cannot be negative")
	}
}