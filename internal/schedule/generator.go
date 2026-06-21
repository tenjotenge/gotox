package schedule

import (
	"context"

	"gotox/internal/analysis"
	"gotox/internal/config"
)

// Task represents a unit of work in the schedule
type Task struct {
	ID          string
	Title       string
	Description string
	EstimatedHours float64
	Dependencies []string
	Assignee    string
	Priority    int
}

// Schedule represents a generated development schedule
type Schedule struct {
	Tasks         []Task
	StartDate     string
	EndDate       string
	TotalDuration int // in days
	Metadata      map[string]interface{}
}

// Generator defines the interface for schedule generation
type Generator interface {
	// GenerateSchedule creates a development schedule based on analysis results
	GenerateSchedule(ctx context.Context, analysis *analysis.AnalysisResult, scheduling config.SchedulingConfig, seed int64) (*Schedule, error)
	
	// ValidateSchedule checks if a schedule is valid and feasible
	ValidateSchedule(schedule *Schedule) error
	
	// ExportSchedule converts schedule to a specific format
	ExportSchedule(schedule *Schedule, format string) (string, error)
}
