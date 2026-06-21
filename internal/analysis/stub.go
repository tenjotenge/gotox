package analysis

import (
	"context"

	"gotox/internal/github"
)

// StubAnalyzer is a stub implementation of Analyzer for testing and scaffolding
type StubAnalyzer struct{}

// NewStubAnalyzer creates a new stub analyzer
func NewStubAnalyzer() *StubAnalyzer {
	return &StubAnalyzer{}
}

// AnalyzeCommits is a stub implementation that returns empty analysis result
func (s *StubAnalyzer) AnalyzeCommits(ctx context.Context, data []github.RepositoryData) (*AnalysisResult, error) {
	return &AnalysisResult{
		DeveloperProfiles:  make(map[string]DeveloperProfile),
		RepositoryProfiles: make(map[string]RepositoryProfile),
		GlobalPatterns:     make(map[string]interface{}),
	}, nil
}

// ExtractDeveloperProfiles is a stub implementation that returns empty profiles
func (s *StubAnalyzer) ExtractDeveloperProfiles(data []github.RepositoryData) (map[string]DeveloperProfile, error) {
	return make(map[string]DeveloperProfile), nil
}

// ExtractRepositoryPatterns is a stub implementation that returns empty patterns
func (s *StubAnalyzer) ExtractRepositoryPatterns(data []github.RepositoryData) (map[string]RepositoryProfile, error) {
	return make(map[string]RepositoryProfile), nil
}
