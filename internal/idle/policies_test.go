package idle

import (
	"testing"
)

func TestNewPolicyManager(t *testing.T) {
	pm := NewPolicyManager()
	if pm == nil {
		t.Fatal("expected non-nil policy manager")
	}

	if pm.templates == nil {
		t.Error("expected non-nil templates map")
	}

	if pm.applied == nil {
		t.Error("expected non-nil applied map")
	}
}

func TestGetTemplate(t *testing.T) {
	pm := NewPolicyManager()

	// Test getting existing template
	template, err := pm.GetTemplate("educational-conservative")
	if err != nil {
		t.Fatalf("failed to get educational-conservative template: %v", err)
	}

	if template == nil {
		t.Fatal("expected non-nil template")
	}

	if template.ID != "educational-conservative" {
		t.Errorf("expected ID 'educational-conservative', got %s", template.ID)
	}

	if template.Category != CategoryEducational {
		t.Errorf("expected category %s, got %s", CategoryEducational, template.Category)
	}

	// Test getting non-existent template
	_, err = pm.GetTemplate("non-existent")
	if err == nil {
		t.Error("expected error for non-existent template")
	}
}

func TestListTemplates(t *testing.T) {
	pm := NewPolicyManager()

	templates := pm.ListTemplates()
	if templates == nil {
		t.Fatal("expected non-nil templates slice")
	}

	if len(templates) == 0 {
		t.Error("expected at least one template")
	}

	// Verify all templates have required fields
	for _, template := range templates {
		if template.ID == "" {
			t.Error("template ID cannot be empty")
		}
		if template.Name == "" {
			t.Error("template name cannot be empty")
		}
		if template.Description == "" {
			t.Error("template description cannot be empty")
		}
		if template.EstimatedSavingsPercent < 0 || template.EstimatedSavingsPercent > 100 {
			t.Errorf("invalid savings percentage: %.1f", template.EstimatedSavingsPercent)
		}
	}
}

func TestGetTemplatesByCategory(t *testing.T) {
	pm := NewPolicyManager()

	// Test educational category
	eduTemplates := pm.GetTemplatesByCategory(CategoryEducational)
	if len(eduTemplates) == 0 {
		t.Error("expected at least one educational template")
	}

	for _, template := range eduTemplates {
		if template.Category != CategoryEducational {
			t.Errorf("expected educational category, got %s", template.Category)
		}
	}

	// Test research category
	researchTemplates := pm.GetTemplatesByCategory(CategoryResearch)
	if len(researchTemplates) == 0 {
		t.Error("expected at least one research template")
	}

	// Test non-existent category
	nonExistentTemplates := pm.GetTemplatesByCategory("non-existent")
	if len(nonExistentTemplates) != 0 {
		t.Error("expected no templates for non-existent category")
	}
}

func TestPolicyTemplateValidation(t *testing.T) {
	pm := NewPolicyManager()

	templates := pm.ListTemplates()
	for _, template := range templates {
		t.Run(template.ID, func(t *testing.T) {
			// Validate basic fields
			if template.ID == "" {
				t.Error("template ID cannot be empty")
			}
			if template.Name == "" {
				t.Error("template name cannot be empty")
			}
			if template.Description == "" {
				t.Error("template description cannot be empty")
			}

			// Validate schedules
			if len(template.Schedules) == 0 {
				t.Error("template must have at least one schedule")
			}

			for _, schedule := range template.Schedules {
				if schedule.ID == "" {
					t.Error("schedule ID cannot be empty")
				}
				if schedule.Name == "" {
					t.Error("schedule name cannot be empty")
				}
				if schedule.IdleMinutes <= 0 {
					t.Error("idle minutes must be positive")
				}
				if schedule.CPUThreshold < 0 || schedule.CPUThreshold > 100 {
					t.Errorf("CPU threshold must be 0-100, got %.1f", schedule.CPUThreshold)
				}
				if schedule.MemoryThreshold < 0 || schedule.MemoryThreshold > 100 {
					t.Errorf("memory threshold must be 0-100, got %.1f", schedule.MemoryThreshold)
				}
				if schedule.GracePeriod < 0 {
					t.Error("grace period cannot be negative")
				}
			}

			// Validate suitable for
			if len(template.SuitableFor) == 0 {
				t.Error("template must specify what it's suitable for")
			}

			// Validate priority
			if template.Priority <= 0 {
				t.Error("template priority must be positive")
			}
		})
	}
}

func TestScheduleTypeConstants(t *testing.T) {
	scheduleTypes := []ScheduleType{
		ScheduleTypeDaily,
		ScheduleTypeWeekly,
		ScheduleTypeWorkHours,
		ScheduleTypeClassHours,
		ScheduleTypeIdleBased,
		ScheduleTypeCustom,
	}

	for _, scheduleType := range scheduleTypes {
		if string(scheduleType) == "" {
			t.Error("schedule type cannot be empty")
		}
	}
}

func TestDayOfWeekConstants(t *testing.T) {
	days := []DayOfWeek{
		Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday,
	}

	expectedDays := []string{
		"monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday",
	}

	if len(days) != len(expectedDays) {
		t.Errorf("expected %d days, got %d", len(expectedDays), len(days))
	}

	for i, day := range days {
		if string(day) != expectedDays[i] {
			t.Errorf("expected day %s, got %s", expectedDays[i], string(day))
		}
	}
}

func TestPolicyCategoryConstants(t *testing.T) {
	categories := []PolicyCategory{
		CategoryAggressive,
		CategoryBalanced,
		CategoryConservative,
		CategoryDevelopment,
		CategoryProduction,
		CategoryResearch,
		CategoryEducational,
		CategoryCustom,
	}

	for _, category := range categories {
		if string(category) == "" {
			t.Error("category cannot be empty")
		}
	}

	// Verify specific values
	if CategoryEducational != "educational" {
		t.Errorf("expected 'educational', got %s", CategoryEducational)
	}

	if CategoryResearch != "research" {
		t.Errorf("expected 'research', got %s", CategoryResearch)
	}
}