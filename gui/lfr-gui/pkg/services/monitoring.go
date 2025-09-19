// Package services provides real-time monitoring services for the GUI
package services

import (
	"log"
	"sync"
	"time"

	"github.com/scttfrdmn/lfr-tools/pkg/api"
)

// MonitoringService provides real-time status monitoring
type MonitoringService struct {
	instanceAPI *api.InstanceAPI
	subscribers map[string][]chan *StatusUpdate
	mu          sync.RWMutex
	stopChan    chan bool
	running     bool
}

// StatusUpdate represents a real-time status update
type StatusUpdate struct {
	Type      string      `json:"type"`      // instance, project, budget, alert
	Timestamp string      `json:"timestamp"`
	Project   string      `json:"project"`
	Data      interface{} `json:"data"`
}

// InstanceStatusUpdate represents instance status changes
type InstanceStatusUpdate struct {
	Name     string `json:"name"`
	State    string `json:"state"`
	PublicIP string `json:"public_ip"`
	Username string `json:"username"`
	Previous string `json:"previous_state"`
}

// ProjectStatusUpdate represents project-level changes
type ProjectStatusUpdate struct {
	Name           string  `json:"name"`
	StudentsOnline int     `json:"students_online"`
	BudgetUsed     float64 `json:"budget_used"`
	AlertCount     int     `json:"alert_count"`
}

// NewMonitoringService creates a new monitoring service
func NewMonitoringService() *MonitoringService {
	return &MonitoringService{
		instanceAPI: api.NewInstanceAPI(),
		subscribers: make(map[string][]chan *StatusUpdate),
		stopChan:    make(chan bool),
		running:     false,
	}
}

// StartMonitoring begins real-time status monitoring
func (m *MonitoringService) StartMonitoring(projects []string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.running {
		return nil // Already running
	}

	m.running = true

	// Start monitoring goroutine
	go m.monitoringLoop(projects)

	log.Printf("Started real-time monitoring for projects: %v", projects)
	return nil
}

// StopMonitoring stops real-time monitoring
func (m *MonitoringService) StopMonitoring() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.running {
		return
	}

	m.running = false
	m.stopChan <- true

	log.Println("Stopped real-time monitoring")
}

// SubscribeToProject subscribes to updates for a specific project
func (m *MonitoringService) SubscribeToProject(project string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Create subscriber channel
	subscriber := make(chan *StatusUpdate, 100)
	m.subscribers[project] = append(m.subscribers[project], subscriber)

	// Start goroutine to forward updates to frontend
	go m.forwardUpdatesToFrontend(project, subscriber)

	return nil
}

// GetProjectStatus returns current project status
func (m *MonitoringService) GetProjectStatus(project string) (*ProjectStatusUpdate, error) {
	projectInfo, err := m.instanceAPI.GetProjectInfo(project)
	if err != nil {
		return nil, err
	}

	return &ProjectStatusUpdate{
		Name:           projectInfo.Name,
		StudentsOnline: projectInfo.RunningCount,
		BudgetUsed:     projectInfo.BudgetUsed,
		AlertCount:     0, // TODO: Implement alert counting
	}, nil
}

// RefreshInstanceStatus manually refreshes instance status for a project
func (m *MonitoringService) RefreshInstanceStatus(project string) error {
	instances, err := m.instanceAPI.ListInstances(project)
	if err != nil {
		return err
	}

	// Send update to subscribers
	update := &StatusUpdate{
		Type:      "instances_refresh",
		Timestamp: time.Now().Format(time.RFC3339),
		Project:   project,
		Data:      instances,
	}

	m.sendUpdate(project, update)
	return nil
}

// monitoringLoop runs the main monitoring loop
func (m *MonitoringService) monitoringLoop(projects []string) {
	ticker := time.NewTicker(30 * time.Second) // Check every 30 seconds
	defer ticker.Stop()

	// Store previous state for comparison
	previousState := make(map[string]map[string]string) // project -> instance -> state

	for {
		select {
		case <-m.stopChan:
			return
		case <-ticker.C:
			for _, project := range projects {
				m.checkProjectChanges(project, previousState)
			}
		}
	}
}

// checkProjectChanges checks for changes in a project's instances
func (m *MonitoringService) checkProjectChanges(project string, previousState map[string]map[string]string) {
	instances, err := m.instanceAPI.ListInstances(project)
	if err != nil {
		log.Printf("Failed to check project %s: %v", project, err)
		return
	}

	if previousState[project] == nil {
		previousState[project] = make(map[string]string)
	}

	var changes []InstanceStatusUpdate
	runningCount := 0

	for _, instance := range instances {
		if instance.State == "running" {
			runningCount++
		}

		// Check for state changes
		prevState := previousState[project][instance.Name]
		if prevState != "" && prevState != instance.State {
			changes = append(changes, InstanceStatusUpdate{
				Name:     instance.Name,
				State:    instance.State,
				PublicIP: instance.PublicIP,
				Username: instance.Username,
				Previous: prevState,
			})
		}

		// Update previous state
		previousState[project][instance.Name] = instance.State
	}

	// Send updates if there are changes
	if len(changes) > 0 {
		update := &StatusUpdate{
			Type:      "instance_changes",
			Timestamp: time.Now().Format(time.RFC3339),
			Project:   project,
			Data:      changes,
		}
		m.sendUpdate(project, update)
	}

	// Send project status update
	projectUpdate := &StatusUpdate{
		Type:      "project_status",
		Timestamp: time.Now().Format(time.RFC3339),
		Project:   project,
		Data: ProjectStatusUpdate{
			Name:           project,
			StudentsOnline: runningCount,
			BudgetUsed:     340.50, // TODO: Calculate real budget usage
			AlertCount:     0,      // TODO: Calculate real alert count
		},
	}
	m.sendUpdate(project, projectUpdate)
}

// sendUpdate sends an update to all subscribers for a project
func (m *MonitoringService) sendUpdate(project string, update *StatusUpdate) {
	m.mu.RLock()
	subscribers := m.subscribers[project]
	m.mu.RUnlock()

	for _, subscriber := range subscribers {
		select {
		case subscriber <- update:
		default:
			// Channel full, skip this update
		}
	}
}

// forwardUpdatesToFrontend forwards updates to the frontend via Wails events
func (m *MonitoringService) forwardUpdatesToFrontend(project string, subscriber chan *StatusUpdate) {
	for update := range subscriber {
		// Log update for now - in production this would emit via Wails event system
		log.Printf("Status update for %s: %s", project, update.Type)

		// TODO: Emit event to frontend when app reference is available
		// app.Event.Emit(fmt.Sprintf("status_update_%s", project), update)
	}
}