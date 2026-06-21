package executor

import (
	"context"

	"gotox/internal/schedule"
)

// StubExecutor is a stub implementation of Executor for testing and scaffolding
type StubExecutor struct{}

// NewStubExecutor creates a new stub executor
func NewStubExecutor() *StubExecutor {
	return &StubExecutor{}
}

// ExecuteTask is a stub implementation that returns success
func (s *StubExecutor) ExecuteTask(task schedule.Task) (*ExecutionResult, error) {
	return &ExecutionResult{
		TaskID:   task.ID,
		Success:  true,
		Output:   "Stub execution",
		Error:    nil,
		Duration: 0,
	}, nil
}

// ExecuteSchedule is a stub implementation that returns empty results
func (s *StubExecutor) ExecuteSchedule(ctx context.Context, schedule *schedule.Schedule) ([]ExecutionResult, error) {
	return []ExecutionResult{}, nil
}

// ValidateEnvironment is a stub implementation that always succeeds
func (s *StubExecutor) ValidateEnvironment() error {
	return nil
}

// Rollback is a stub implementation that does nothing
func (s *StubExecutor) Rollback(taskID string) error {
	return nil
}
