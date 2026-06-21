package schedule

import (
	"context"

	"gotox/internal/analysis"
	"gotox/internal/config"
)

// StubGenerator is a stub implementation of Generator for testing and scaffolding
type StubGenerator struct{}

// NewStubGenerator creates a new stub generator
func NewStubGenerator() *StubGenerator {
	return &StubGenerator{}
}

// GenerateSchedule is a stub implementation that returns empty schedule
func (s *StubGenerator) GenerateSchedule(ctx context.Context, analysis *analysis.AnalysisResult, scheduling config.SchedulingConfig, seed int64) (*Schedule, error) {
	return &Schedule{
		Tasks:         []Task{},
		StartDate:     "",
		EndDate:       "",
		TotalDuration: 0,
		Metadata:      make(map[string]interface{}),
	}, nil
}

// ValidateSchedule is a stub implementation that always succeeds
func (s *StubGenerator) ValidateSchedule(schedule *Schedule) error {
	return nil
}

// ExportSchedule is a stub implementation that returns empty string
func (s *StubGenerator) ExportSchedule(schedule *Schedule, format string) (string, error) {
	return "", nil
}
