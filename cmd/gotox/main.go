package main

import (
	"gotox/internal/analysis"
	"gotox/internal/executor"
	"gotox/internal/github"
	"gotox/internal/pipeline"
	"gotox/internal/schedule"
	"gotox/internal/ui"
	"gotox/internal/verify"
)

func main() {
	// Create stub implementations for all pipeline dependencies
	ingestor := github.NewStubIngestor()
	analyzer := analysis.NewStubAnalyzer()
	generator := schedule.NewStubGenerator()
	executor := executor.NewStubExecutor()
	verifier := verify.NewStubVerifier()

	// Create pipeline with injected dependencies
	pipeline := pipeline.NewPipeline(ingestor, analyzer, generator, executor, verifier)

	// Create and run UI with pipeline
	app := ui.NewApp()
	app.SetPipeline(pipeline)
	app.Run()
}
