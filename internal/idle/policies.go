// Package idle provides advanced idle detection and policy management for Lightsail for Research.
package idle

import (
	"fmt"
	"time"
)

// PolicyTemplate represents a pre-configured idle detection policy.
type PolicyTemplate struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Category    PolicyCategory    `json:"category"`
	Schedules   []Schedule        `json:"schedules"`
	Tags        map[string]string `json:"tags"`

	// Cost analysis
	EstimatedSavingsPercent float64  `json:"estimated_savings_percent"`
	SuitableFor             []string `json:"suitable_for"`

	// Configuration
	AutoApply bool     `json:"auto_apply"`
	Priority  int      `json:"priority"`
	Conflicts []string `json:"conflicts"` // IDs of conflicting templates
}

// PolicyCategory categorizes idle detection policies.
type PolicyCategory string

const (
	CategoryAggressive   PolicyCategory = "aggressive"
	CategoryBalanced     PolicyCategory = "balanced"
	CategoryConservative PolicyCategory = "conservative"
	CategoryDevelopment  PolicyCategory = "development"
	CategoryProduction   PolicyCategory = "production"
	CategoryResearch     PolicyCategory = "research"
	CategoryEducational  PolicyCategory = "educational"
	CategoryCustom       PolicyCategory = "custom"
)

// Schedule represents an idle detection schedule with multi-metric thresholds.
type Schedule struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Type        ScheduleType `json:"type"`
	Enabled     bool         `json:"enabled"`

	// Time-based scheduling
	StartTime  string      `json:"start_time,omitempty"` // HH:MM format
	EndTime    string      `json:"end_time,omitempty"`   // HH:MM format
	DaysOfWeek []DayOfWeek `json:"days_of_week,omitempty"`
	Timezone   string      `json:"timezone,omitempty"`

	// Multi-metric idle detection thresholds
	IdleMinutes      int     `json:"idle_minutes"`      // Minutes of inactivity
	CPUThreshold     float64 `json:"cpu_threshold"`     // CPU usage % threshold
	MemoryThreshold  float64 `json:"memory_threshold"`  // Memory usage % threshold
	NetworkThreshold float64 `json:"network_threshold"` // Network I/O threshold
	SSHConnections   int     `json:"ssh_connections"`   // Active SSH sessions

	// Actions
	Action       string `json:"action"`        // stop, hibernate, alert
	GracePeriod  int    `json:"grace_period"`  // Minutes before action
	PreStopAlert bool   `json:"pre_stop_alert"` // Send alert before stopping

	// Cost tracking
	EstimatedMonthlySavings float64   `json:"estimated_monthly_savings"`
	LastExecuted            time.Time `json:"last_executed"`
	TotalSavings            float64   `json:"total_savings"`
}

// ScheduleType defines the type of idle detection schedule.
type ScheduleType string

const (
	ScheduleTypeDaily       ScheduleType = "daily"
	ScheduleTypeWeekly      ScheduleType = "weekly"
	ScheduleTypeWorkHours   ScheduleType = "work_hours"
	ScheduleTypeClassHours  ScheduleType = "class_hours"
	ScheduleTypeIdleBased   ScheduleType = "idle_based"
	ScheduleTypeCustom      ScheduleType = "custom"
)

// DayOfWeek represents a day of the week.
type DayOfWeek string

const (
	Monday    DayOfWeek = "monday"
	Tuesday   DayOfWeek = "tuesday"
	Wednesday DayOfWeek = "wednesday"
	Thursday  DayOfWeek = "thursday"
	Friday    DayOfWeek = "friday"
	Saturday  DayOfWeek = "saturday"
	Sunday    DayOfWeek = "sunday"
)

// PolicyManager manages idle detection policies for educational environments.
type PolicyManager struct {
	templates map[string]*PolicyTemplate
	applied   map[string][]string // instance -> policy IDs
}

// NewPolicyManager creates a new policy manager with educational templates.
func NewPolicyManager() *PolicyManager {
	pm := &PolicyManager{
		templates: make(map[string]*PolicyTemplate),
		applied:   make(map[string][]string),
	}

	// Load educational policy templates
	pm.loadEducationalTemplates()

	return pm
}

// loadEducationalTemplates loads pre-configured templates for educational use.
func (pm *PolicyManager) loadEducationalTemplates() {
	// Educational - Conservative (for classes)
	pm.templates["educational-conservative"] = &PolicyTemplate{
		ID:          "educational-conservative",
		Name:        "Educational Conservative",
		Description: "Safe idle detection for student environments. Minimal risk of interrupting work.",
		Category:    CategoryEducational,
		Schedules: []Schedule{
			{
				ID:               "class-hours-conservative",
				Name:             "Class Hours - Conservative",
				Type:             ScheduleTypeClassHours,
				Enabled:          true,
				StartTime:        "08:00",
				EndTime:          "18:00",
				DaysOfWeek:       []DayOfWeek{Monday, Tuesday, Wednesday, Thursday, Friday},
				IdleMinutes:      180, // 3 hours
				CPUThreshold:     2.0, // Very low CPU
				MemoryThreshold:  5.0, // Very low memory changes
				NetworkThreshold: 1.0, // Very low network
				SSHConnections:   0,   // No active SSH
				Action:           "stop",
				GracePeriod:      15, // 15 minute warning
				PreStopAlert:     true,
			},
		},
		EstimatedSavingsPercent: 40.0,
		SuitableFor:             []string{"students", "beginners", "classes"},
		AutoApply:               false,
		Priority:                1,
	}

	// Educational - Balanced (for labs)
	pm.templates["educational-balanced"] = &PolicyTemplate{
		ID:          "educational-balanced",
		Name:        "Educational Balanced",
		Description: "Balanced idle detection for lab environments. Good cost savings with reasonable safety.",
		Category:    CategoryEducational,
		Schedules: []Schedule{
			{
				ID:               "lab-hours-balanced",
				Name:             "Lab Hours - Balanced",
				Type:             ScheduleTypeWorkHours,
				Enabled:          true,
				StartTime:        "07:00",
				EndTime:          "20:00",
				DaysOfWeek:       []DayOfWeek{Monday, Tuesday, Wednesday, Thursday, Friday},
				IdleMinutes:      120, // 2 hours (LfR default)
				CPUThreshold:     5.0,
				MemoryThreshold:  10.0,
				NetworkThreshold: 5.0,
				SSHConnections:   0,
				Action:           "stop",
				GracePeriod:      10,
				PreStopAlert:     true,
			},
		},
		EstimatedSavingsPercent: 60.0,
		SuitableFor:             []string{"lab-work", "research", "development"},
		AutoApply:               false,
		Priority:                2,
	}

	// Research - Long Running
	pm.templates["research-long-running"] = &PolicyTemplate{
		ID:          "research-long-running",
		Name:        "Research Long Running",
		Description: "Extended idle detection for long-running research tasks. Minimal interruptions.",
		Category:    CategoryResearch,
		Schedules: []Schedule{
			{
				ID:               "research-extended",
				Name:             "Research Extended Hours",
				Type:             ScheduleTypeCustom,
				Enabled:          true,
				IdleMinutes:      480, // 8 hours
				CPUThreshold:     1.0, // Very sensitive to any activity
				MemoryThreshold:  2.0,
				NetworkThreshold: 1.0,
				SSHConnections:   0,
				Action:           "stop",
				GracePeriod:      30, // Long warning period
				PreStopAlert:     true,
			},
		},
		EstimatedSavingsPercent: 25.0,
		SuitableFor:             []string{"research", "long-tasks", "data-processing"},
		AutoApply:               false,
		Priority:                3,
	}

	// Development - Aggressive
	pm.templates["development-aggressive"] = &PolicyTemplate{
		ID:          "development-aggressive",
		Name:        "Development Aggressive",
		Description: "Aggressive cost optimization for development environments. Quick shutdown when idle.",
		Category:    CategoryDevelopment,
		Schedules: []Schedule{
			{
				ID:               "dev-aggressive",
				Name:             "Development Quick Stop",
				Type:             ScheduleTypeIdleBased,
				Enabled:          true,
				IdleMinutes:      30, // 30 minutes
				CPUThreshold:     10.0,
				MemoryThreshold:  15.0,
				NetworkThreshold: 10.0,
				SSHConnections:   0,
				Action:           "stop",
				GracePeriod:      5,
				PreStopAlert:     true,
			},
		},
		EstimatedSavingsPercent: 80.0,
		SuitableFor:             []string{"development", "testing", "short-tasks"},
		AutoApply:               false,
		Priority:                4,
	}
}

// GetTemplate returns a policy template by ID.
func (pm *PolicyManager) GetTemplate(id string) (*PolicyTemplate, error) {
	template, exists := pm.templates[id]
	if !exists {
		return nil, fmt.Errorf("policy template not found: %s", id)
	}
	return template, nil
}

// ListTemplates returns all available policy templates.
func (pm *PolicyManager) ListTemplates() []*PolicyTemplate {
	var templates []*PolicyTemplate
	for _, template := range pm.templates {
		templates = append(templates, template)
	}
	return templates
}

// GetTemplatesByCategory returns templates filtered by category.
func (pm *PolicyManager) GetTemplatesByCategory(category PolicyCategory) []*PolicyTemplate {
	var templates []*PolicyTemplate
	for _, template := range pm.templates {
		if template.Category == category {
			templates = append(templates, template)
		}
	}
	return templates
}