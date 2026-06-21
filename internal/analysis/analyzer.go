package analysis

import (
	"context"

	"gotox/internal/github"
)

// DeveloperProfile represents a developer's coding patterns
type DeveloperProfile struct {
	Username          string
	TotalCommits      int
	AverageCommitsPerDay float64
	MostActiveHour    int
	MostActiveDay     string
	FileTypeDistribution map[string]int
}

// RepositoryProfile represents patterns for a repository
type RepositoryProfile struct {
	Owner            string
	Name             string
	TotalCommits     int
	ActiveContributors []string
	CommitFrequency  map[string]int // day -> commit count
}

// AnalysisResult contains the output of commit history analysis
type AnalysisResult struct {
	DeveloperProfiles  map[string]DeveloperProfile
	RepositoryProfiles map[string]RepositoryProfile
	GlobalPatterns     map[string]interface{}
}

// Analyzer defines the interface for commit history analysis
type Analyzer interface {
	// AnalyzeCommits processes commit data and extracts patterns
	AnalyzeCommits(ctx context.Context, data []github.RepositoryData) (*AnalysisResult, error)
	
	// ExtractDeveloperProfiles creates profiles for individual developers
	ExtractDeveloperProfiles(data []github.RepositoryData) (map[string]DeveloperProfile, error)
	
	// ExtractRepositoryPatterns identifies repository-level patterns
	ExtractRepositoryPatterns(data []github.RepositoryData) (map[string]RepositoryProfile, error)
}
