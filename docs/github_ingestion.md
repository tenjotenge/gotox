# GitHub Ingestion Module

## Overview

The GitHub ingestion module (`internal/github`) is responsible for fetching commit history and repository data from GitHub's API. This module implements the `Ingestor` interface.

## Interface Definition

```go
type Ingestor interface {
    FetchCommits(cfg *config.GitHubConfig, startYear, endYear int) ([]RepositoryData, error)
    ValidateCredentials(username, pat string) error
    FetchRepositories(username, pat string) ([]config.Repository, error)
}
```

## Data Structures

### Commit
```go
type Commit struct {
    SHA          string
    Author       string
    Message      string
    Timestamp    int64
    Additions    int
    Deletions    int
    ChangedFiles int
}
```

### RepositoryData
```go
type RepositoryData struct {
    Owner   string
    Name    string
    Commits []Commit
}
```

## Implementation Requirements

### 1. Authentication
- Use GitHub Personal Access Token (PAT) for authentication
- PAT should have `repo` scope for private repositories
- Public repositories can be accessed without PAT (with rate limits)
- Validate credentials before making API calls

### 2. Rate Limit Handling
- GitHub API has rate limits (5000 requests/hour with token, 60/hour without)
- Implement exponential backoff on rate limit errors
- Cache responses to minimize API calls
- Check rate limit status from response headers

### 3. Commit Fetching
- Use GitHub REST API: `GET /repos/{owner}/{repo}/commits`
- Fetch commits paginated (default 30 per page, max 100)
- Filter commits by date range (startYear to endYear)
- Fetch commit details including:
  - SHA
  - Author name and email
  - Commit timestamp
  - Commit message
  - Stats (additions, deletions, files changed)

### 4. Repository Listing
- Use GitHub REST API: `GET /users/{username}/repos`
- List all repositories for the authenticated user
- Include both owned and collaborated repositories
- Allow filtering by repository name pattern

### 5. Error Handling
- Handle network timeouts
- Handle authentication failures (401)
- Handle rate limit errors (403)
- Handle not found errors (404)
- Handle server errors (5xx)

## API Endpoints

### Get Repository Commits
```
GET /repos/{owner}/{repo}/commits
Parameters:
  - since: ISO 8601 timestamp (start of year range)
  - until: ISO 8601 timestamp (end of year range)
  - per_page: 1-100 (default 30)
  - page: page number for pagination
```

### Get User Repositories
```
GET /users/{username}/repos
Parameters:
  - type: all, owner, member (default: all)
  - per_page: 1-100 (default 30)
  - page: page number for pagination
```

### Get Commit Details
```
GET /repos/{owner}/{repo}/commits/{sha}
Returns detailed commit information including stats
```

## Implementation Notes

### Pagination
- Use Link header to navigate pages
- Fetch all pages until no more results
- Handle empty repository cases

### Date Filtering
- Convert year range to ISO 8601 timestamps
- Use `since` and `until` parameters
- Filter client-side as backup

### Caching Strategy
- Cache repository lists (valid for 1 hour)
- Cache commit data (valid for 24 hours)
- Use filesystem or in-memory cache
- Invalidate cache on configuration change

### Concurrent Fetching
- Fetch multiple repositories concurrently
- Limit concurrent requests to avoid rate limits
- Use worker pool pattern
- Aggregate results

## Testing Considerations

### Unit Tests
- Mock GitHub API responses
- Test pagination handling
- Test error scenarios
- Test date filtering logic

### Integration Tests
- Use GitHub test fixtures
- Test with real API (use test account)
- Test rate limit handling
- Test concurrent fetching

## Dependencies

- Go standard library
- HTTP client for API calls
- JSON parsing
- Time handling for date conversion

## Future Enhancements

- Support GraphQL API for more efficient queries
- Support GitHub Enterprise
- Support webhook-based real-time updates
- Support branch filtering
- Support commit message filtering
