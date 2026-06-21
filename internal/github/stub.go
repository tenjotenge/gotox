package github

import (
	"context"

	"gotox/internal/config"
)

// StubIngestor is a stub implementation of Ingestor for testing and scaffolding
type StubIngestor struct{}

// NewStubIngestor creates a new stub ingestor
func NewStubIngestor() *StubIngestor {
	return &StubIngestor{}
}

// FetchCommits is a stub implementation that returns empty data
func (s *StubIngestor) FetchCommits(ctx context.Context, cfg *config.GitHubConfig, startYear, endYear int) ([]RepositoryData, error) {
	return []RepositoryData{}, nil
}

// ValidateCredentials is a stub implementation that always succeeds
func (s *StubIngestor) ValidateCredentials(username, pat string) error {
	return nil
}

// FetchRepositories is a stub implementation that returns empty list
func (s *StubIngestor) FetchRepositories(username, pat string) ([]config.Repository, error) {
	return []config.Repository{}, nil
}
