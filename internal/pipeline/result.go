package pipeline

import (
	"gotox/internal/analysis"
	"gotox/internal/executor"
	"gotox/internal/github"
	"gotox/internal/schedule"
)

// Result represents the outcome of a pipeline run.
// Artifact fields are the source of truth; summary fields are derived views.
type Result struct {
	RepositoryData    []github.RepositoryData
	AnalysisResult    *analysis.AnalysisResult
	GeneratedSchedule *schedule.Schedule
	ExecutionResults  []executor.ExecutionResult

	Ingestion     *IngestionSummary
	Analysis      *AnalysisSummary
	Schedule      *ScheduleSummary
	Execution     *ExecutionSummary
	Verification  *VerificationSummary
	Error         error
}

// IngestionSummary summarizes the GitHub ingestion phase
type IngestionSummary struct {
	RepositoriesFetched int
	TotalCommits        int
	DurationMs          int64
	Success             bool
}

// AnalysisSummary summarizes the analysis phase
type AnalysisSummary struct {
	DevelopersAnalyzed   int
	RepositoriesAnalyzed int
	PatternsExtracted    int
	DurationMs           int64
	Success              bool
}

// ScheduleSummary summarizes the schedule generation phase
type ScheduleSummary struct {
	TasksGenerated    int
	TotalDurationDays int
	DurationMs        int64
	Success           bool
}

// ExecutionSummary summarizes the execution phase
type ExecutionSummary struct {
	TasksExecuted  int
	TasksSucceeded int
	TasksFailed    int
	DurationMs     int64
	Success        bool
}

// VerificationSummary summarizes the verification phase
type VerificationSummary struct {
	ItemsVerified int
	ItemsPassed   int
	ItemsFailed   int
	DurationMs    int64
	Success       bool
}

// NewResult creates a new empty Result
func NewResult() *Result {
	return &Result{
		Ingestion:    &IngestionSummary{},
		Analysis:     &AnalysisSummary{},
		Schedule:     &ScheduleSummary{},
		Execution:    &ExecutionSummary{},
		Verification: &VerificationSummary{},
	}
}

// FromGitHubData creates an IngestionSummary from GitHub data
func FromGitHubData(data []github.RepositoryData, durationMs int64, success bool) *IngestionSummary {
	totalCommits := 0
	for _, repo := range data {
		totalCommits += len(repo.Commits)
	}

	return &IngestionSummary{
		RepositoriesFetched: len(data),
		TotalCommits:        totalCommits,
		DurationMs:          durationMs,
		Success:             success,
	}
}

// FromAnalysisResult creates an AnalysisSummary from analysis result
func FromAnalysisResult(result *analysis.AnalysisResult, durationMs int64, success bool) *AnalysisSummary {
	if result == nil {
		return &AnalysisSummary{
			DurationMs: durationMs,
			Success:    success,
		}
	}

	return &AnalysisSummary{
		DevelopersAnalyzed:   len(result.DeveloperProfiles),
		RepositoriesAnalyzed: len(result.RepositoryProfiles),
		PatternsExtracted:    len(result.GlobalPatterns),
		DurationMs:           durationMs,
		Success:              success,
	}
}

// FromSchedule creates a ScheduleSummary from schedule
func FromSchedule(sched *schedule.Schedule, durationMs int64, success bool) *ScheduleSummary {
	if sched == nil {
		return &ScheduleSummary{
			DurationMs: durationMs,
			Success:    success,
		}
	}

	return &ScheduleSummary{
		TasksGenerated:    len(sched.Tasks),
		TotalDurationDays: sched.TotalDuration,
		DurationMs:        durationMs,
		Success:           success,
	}
}

// FromExecutionResults creates an ExecutionSummary from execution results
func FromExecutionResults(results []executor.ExecutionResult, durationMs int64, success bool) *ExecutionSummary {
	successCount := 0
	failedCount := 0
	for _, r := range results {
		if r.Success {
			successCount++
		} else {
			failedCount++
		}
	}

	return &ExecutionSummary{
		TasksExecuted:  len(results),
		TasksSucceeded: successCount,
		TasksFailed:    failedCount,
		DurationMs:     durationMs,
		Success:        success && failedCount == 0,
	}
}
