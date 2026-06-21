package pipeline

import (
	"context"
	"time"

	"gotox/internal/analysis"
	"gotox/internal/config"
	"gotox/internal/executor"
	"gotox/internal/github"
	"gotox/internal/schedule"
	"gotox/internal/verify"
)

// Pipeline orchestrates the execution flow through all subsystems
type Pipeline struct {
	ingestor  github.Ingestor
	analyzer  analysis.Analyzer
	generator schedule.Generator
	executor  executor.Executor
	verifier  verify.Verifier
	progress  *ProgressReporter
}

// NewPipeline creates a new pipeline with injected dependencies
func NewPipeline(
	ingestor github.Ingestor,
	analyzer analysis.Analyzer,
	generator schedule.Generator,
	executor executor.Executor,
	verifier verify.Verifier,
) *Pipeline {
	return &Pipeline{
		ingestor:  ingestor,
		analyzer:  analyzer,
		generator: generator,
		executor:  executor,
		verifier:  verifier,
		progress:  NewProgressReporter(),
	}
}

// AddProgressHandler adds a progress event handler to the pipeline
func (p *Pipeline) AddProgressHandler(handler ProgressHandler) {
	p.progress.AddHandler(handler)
}

// Run executes the pipeline with the given configuration
func (p *Pipeline) Run(ctx context.Context, cfg config.Config) (*Result, error) {
	result := NewResult()

	// Phase 1: Ingestion
	p.progress.PhaseStarted(PhaseIngestion, "Starting GitHub data ingestion")
	ingestionData, err := p.runIngestion(ctx, cfg, result)
	if err != nil {
		p.progress.PhaseFailed(PhaseIngestion, "GitHub ingestion failed", err)
		result.Error = err
		return result, err
	}
	p.progress.PhaseCompleted(PhaseIngestion, "GitHub ingestion completed")

	// Phase 2: Analysis
	p.progress.PhaseStarted(PhaseAnalysis, "Starting commit history analysis")
	analysisResult, err := p.runAnalysis(ctx, ingestionData, result)
	if err != nil {
		p.progress.PhaseFailed(PhaseAnalysis, "Analysis failed", err)
		result.Error = err
		return result, err
	}
	p.progress.PhaseCompleted(PhaseAnalysis, "Analysis completed")

	// Phase 3: Schedule Generation
	p.progress.PhaseStarted(PhaseSchedule, "Starting schedule generation")
	scheduleResult, err := p.runSchedule(ctx, analysisResult, cfg.Scheduling, cfg.RandomSeed, result)
	if err != nil {
		p.progress.PhaseFailed(PhaseSchedule, "Schedule generation failed", err)
		result.Error = err
		return result, err
	}
	p.progress.PhaseCompleted(PhaseSchedule, "Schedule generation completed")

	// Phase 4: Execution (optional - can be skipped)
	if p.executor != nil {
		p.progress.PhaseStarted(PhaseExecution, "Starting schedule execution")
		if err := p.runExecution(ctx, scheduleResult, result); err != nil {
			p.progress.PhaseFailed(PhaseExecution, "Execution failed", err)
			result.Error = err
			return result, err
		}
		p.progress.PhaseCompleted(PhaseExecution, "Execution completed")
	}

	// Phase 5: Verification (optional - can be skipped)
	if p.verifier != nil {
		p.progress.PhaseStarted(PhaseVerification, "Starting verification")
		verificationResult, err := p.runVerification(ctx, scheduleResult)
		if err != nil {
			p.progress.PhaseFailed(PhaseVerification, "Verification failed", err)
			result.Error = err
			return result, err
		}
		p.progress.PhaseCompleted(PhaseVerification, "Verification completed")
		result.Verification = verificationResult
	}

	return result, nil
}

// runIngestion executes the GitHub ingestion phase
func (p *Pipeline) runIngestion(ctx context.Context, cfg config.Config, pipelineResult *Result) ([]github.RepositoryData, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	start := time.Now()

	data, err := p.ingestor.FetchCommits(ctx, &cfg.GitHub, cfg.Analysis.StartYear, cfg.Analysis.EndYear)
	if err != nil {
		return nil, err
	}

	duration := time.Since(start).Milliseconds()
	pipelineResult.RepositoryData = data
	pipelineResult.Ingestion = FromGitHubData(data, duration, true)

	return data, nil
}

// runAnalysis executes the analysis phase
func (p *Pipeline) runAnalysis(ctx context.Context, data []github.RepositoryData, pipelineResult *Result) (*analysis.AnalysisResult, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	start := time.Now()

	analysisResult, err := p.analyzer.AnalyzeCommits(ctx, data)
	if err != nil {
		return nil, err
	}

	duration := time.Since(start).Milliseconds()
	pipelineResult.AnalysisResult = analysisResult
	pipelineResult.Analysis = FromAnalysisResult(analysisResult, duration, true)

	return analysisResult, nil
}

// runSchedule executes the schedule generation phase
func (p *Pipeline) runSchedule(ctx context.Context, analysisResult *analysis.AnalysisResult, scheduling config.SchedulingConfig, seed int64, pipelineResult *Result) (*schedule.Schedule, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	start := time.Now()

	scheduleResult, err := p.generator.GenerateSchedule(ctx, analysisResult, scheduling, seed)
	if err != nil {
		return nil, err
	}

	duration := time.Since(start).Milliseconds()
	pipelineResult.GeneratedSchedule = scheduleResult
	pipelineResult.Schedule = FromSchedule(scheduleResult, duration, true)

	return scheduleResult, nil
}

// runExecution executes the execution phase
func (p *Pipeline) runExecution(ctx context.Context, sched *schedule.Schedule, pipelineResult *Result) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	start := time.Now()

	results, err := p.executor.ExecuteSchedule(ctx, sched)
	if err != nil {
		return err
	}

	duration := time.Since(start).Milliseconds()
	pipelineResult.ExecutionResults = results
	pipelineResult.Execution = FromExecutionResults(results, duration, true)

	return nil
}

// runVerification executes the verification phase
func (p *Pipeline) runVerification(ctx context.Context, sched *schedule.Schedule) (*VerificationSummary, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	start := time.Now()

	result, err := p.verifier.VerifySchedule(sched)
	if err != nil {
		return nil, err
	}

	duration := time.Since(start).Milliseconds()

	return &VerificationSummary{
		ItemsVerified: len(result.Errors) + len(result.Warnings),
		ItemsPassed:     len(result.Warnings),
		ItemsFailed:     len(result.Errors),
		DurationMs:      duration,
		Success:         result.Valid,
	}, nil
}
