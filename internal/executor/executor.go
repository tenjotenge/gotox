package executor

import (
	"context"

	"gotox/internal/schedule"
)

// ExecutionResult represents the result of executing a schedule
type ExecutionResult struct {
	TaskID      string
	Success     bool
	Output      string
	Error       error
	Duration    int64 // in milliseconds
}

// Executor defines the interface for executing generated schedules
type Executor interface {
	// ExecuteTask performs a single task from the schedule
	ExecuteTask(task schedule.Task) (*ExecutionResult, error)
	
	// ExecuteSchedule runs all tasks in a schedule
	ExecuteSchedule(ctx context.Context, schedule *schedule.Schedule) ([]ExecutionResult, error)
	
	// ValidateEnvironment checks if the execution environment is ready
	ValidateEnvironment() error
	
	// Rollback reverts changes made during execution
	Rollback(taskID string) error
}
