package main

import (
	"log"
	"os"

	"gotox/internal/analysis"
	"gotox/internal/executor"
	"gotox/internal/github"
	"gotox/internal/pipeline"
	"gotox/internal/schedule"
	"gotox/internal/ui"
	"gotox/internal/verify"
)

func main() {
	// Use real GitHub client if PAT is provided, otherwise use stub
	pat := os.Getenv("GITHUB_PAT")
	var ingestor github.Ingestor

	if pat != "" {
		ingestor = github.NewClient(pat)
		log.Println("[gotox] Running in REAL mode with GitHub API integration")
	} else {
		ingestor = github.NewStubIngestor()
		log.Println("[gotox] Running in STUB mode - no GitHub PAT provided")
	}

	// Create stub implementations for other pipeline dependencies
	analyzer := analysis.NewStubAnalyzer()
	generator := schedule.NewStubGenerator()
	executor := executor.NewStubExecutor()
	verifier := verify.NewStubVerifier()

	log.Println("[gotox] Pipeline components initialized:")
	log.Printf("[gotox]   - Ingestor: %T", ingestor)
	log.Printf("[gotox]   - Analyzer: %T", analyzer)
	log.Printf("[gotox]   - Generator: %T", generator)
	log.Printf("[gotox]   - Executor: %T", executor)
	log.Printf("[gotox]   - Verifier: %T", verifier)

	// Create pipeline with injected dependencies
	pipeline := pipeline.NewPipeline(ingestor, analyzer, generator, executor, verifier)

	// Create and run UI with pipeline
	app := ui.NewApp()
	app.SetPipeline(pipeline)
	log.Println("[gotox] Starting UI...")
	app.Run()
}
