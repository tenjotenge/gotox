package github

import (
	"context"

	"gotox/internal/config"
)

// Commit represents a GitHub commit with metadata
type Commit struct {
	SHA          string
	Author       string
	Message      string
	Timestamp    int64
	Additions    int
	Deletions    int
	ChangedFiles int
}

// RepositoryData holds all commits for a repository
type RepositoryData struct {
	Owner   string
	Name    string
	Commits []Commit
}

// Ingestor defines the interface for GitHub data ingestion
type Ingestor interface {
	// FetchCommits retrieves commit history for configured repositories
	FetchCommits(ctx context.Context, cfg *config.GitHubConfig, startYear, endYear int) ([]RepositoryData, error)

	// ValidateCredentials checks if GitHub credentials are valid
	ValidateCredentials(username, pat string) error

	// FetchRepositories lists repositories for a user
	FetchRepositories(ctx context.Context, username, pat string) ([]config.Repository, error)
}
