package gitutil

import (
	"time"
)

// GitUtil provides git-related utility functions
// This package will contain helper functions for git operations
// such as commit parsing, diff analysis, and repository management

// GitCommit represents a parsed git commit
type GitCommit struct {
	Hash      string
	Author    string
	Email     string
	Date      string
	Message   string
	Files     []string
}

// Parser defines the interface for parsing git data
type Parser interface {
	ParseCommit(data string) (*GitCommit, error)
	ParseDiff(data string) (map[string]string, error)
}

// Helper provides git utility functions
type Helper struct{}

// NewHelper creates a new git utility helper
func NewHelper() *Helper {
	return &Helper{}
}

// ToGitHubCommit converts a GitCommit to a github.Commit-like structure
// This adapter function enables data flow between git parsing and GitHub-based analysis
func (g *GitCommit) ToGitHubCommit() map[string]interface{} {
	return map[string]interface{}{
		"sha":           g.Hash,
		"author":        g.Author,
		"message":       g.Message,
		"timestamp":     parseGitDate(g.Date),
		"changed_files": len(g.Files),
	}
}

// FromGitHubCommit creates a GitCommit from GitHub commit data
// This adapter function enables data flow from GitHub API to git utilities
func FromGitHubCommit(sha, author, message string, timestamp int64, files []string) *GitCommit {
	return &GitCommit{
		Hash:    sha,
		Author:  author,
		Email:   "", // GitHub API may not provide email in all contexts
		Date:    time.Unix(timestamp, 0).Format("2006-01-02 15:04:05"),
		Message: message,
		Files:   files,
	}
}

// parseGitDate parses a git date string to Unix timestamp
func parseGitDate(dateStr string) int64 {
	// Try common git date formats
	formats := []string{
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05Z",
		time.RFC3339,
	}
	
	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t.Unix()
		}
	}
	
	// Default to current time if parsing fails
	return time.Now().Unix()
}
