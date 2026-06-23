package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"

	"gotox/internal/config"
)

const (
	githubAPIURL = "https://api.github.com"
	maxRetries   = 3
	baseDelay    = 1 * time.Second
)

// Client implements the Ingestor interface using GitHub API
type Client struct {
	httpClient *http.Client
	pat        string
	cache      *sync.Map
}

// NewClient creates a new GitHub API client
func NewClient(pat string) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		pat:   pat,
		cache: &sync.Map{},
	}
}

// GitHub API response structures
type githubRepo struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Owner    struct {
		Login string `json:"login"`
	} `json:"owner"`
	Private bool `json:"private"`
}

type githubCommit struct {
	SHA    string `json:"sha"`
	Commit struct {
		Author struct {
			Name  string `json:"name"`
			Email string `json:"email"`
			Date  string `json:"date"`
		} `json:"author"`
		Message string `json:"message"`
	} `json:"commit"`
	Author struct {
		Login string `json:"login"`
	} `json:"author"`
	Stats struct {
		Additions int `json:"additions"`
		Deletions int `json:"deletions"`
	} `json:"stats"`
}

// FetchCommits retrieves commit history for configured repositories
func (c *Client) FetchCommits(ctx context.Context, cfg *config.GitHubConfig, startYear, endYear int) ([]RepositoryData, error) {
	var allData []RepositoryData
	var errors []error

	for _, repo := range cfg.Repositories {
		select {
		case <-ctx.Done():
			return allData, ctx.Err()
		default:
		}

		data, err := c.fetchRepoCommits(ctx, repo.Owner, repo.Name, startYear, endYear)
		if err != nil {
			errors = append(errors, fmt.Errorf("failed to fetch %s/%s: %w", repo.Owner, repo.Name, err))
			continue
		}
		allData = append(allData, data)
	}

	if len(errors) > 0 && len(allData) == 0 {
		return nil, fmt.Errorf("all repositories failed: %v", errors)
	}

	return allData, nil
}

// fetchRepoCommits fetches commits for a single repository
func (c *Client) fetchRepoCommits(ctx context.Context, owner, name string, startYear, endYear int) (RepositoryData, error) {
	cacheKey := fmt.Sprintf("%s/%s/%d-%d", owner, name, startYear, endYear)
	if cached, ok := c.cache.Load(cacheKey); ok {
		return cached.(RepositoryData), nil
	}

	var commits []Commit
	page := 1
	perPage := 100

	startDate := fmt.Sprintf("%d-01-01T00:00:00Z", startYear)
	endDate := fmt.Sprintf("%d-12-31T23:59:59Z", endYear)

	for {
		select {
		case <-ctx.Done():
			return RepositoryData{}, ctx.Err()
		default:
		}

		url := fmt.Sprintf("%s/repos/%s/%s/commits?since=%s&until=%s&page=%d&per_page=%d",
			githubAPIURL, owner, name, startDate, endDate, page, perPage)

		var pageCommits []githubCommit
		resp, err := c.doRequest(ctx, url)
		if err != nil {
			return RepositoryData{}, err
		}

		if err := json.Unmarshal(resp, &pageCommits); err != nil {
			return RepositoryData{}, fmt.Errorf("failed to parse commits: %w", err)
		}

		if len(pageCommits) == 0 {
			break
		}

		commits = append(commits, c.normalizeCommits(pageCommits)...)

		// Check if there are more pages
		if len(pageCommits) < perPage {
			break
		}

		page++
	}

	data := RepositoryData{
		Owner:   owner,
		Name:    name,
		Commits: commits,
	}

	c.cache.Store(cacheKey, data)
	return data, nil
}

// normalizeCommits converts GitHub API commits to internal Commit model
func (c *Client) normalizeCommits(githubCommits []githubCommit) []Commit {
	commits := make([]Commit, 0, len(githubCommits))
	for _, gc := range githubCommits {
		timestamp, err := time.Parse(time.RFC3339, gc.Commit.Author.Date)
		if err != nil {
			continue
		}

		commits = append(commits, Commit{
			SHA:          gc.SHA,
			Author:       gc.Author.Login,
			Message:      gc.Commit.Message,
			Timestamp:    timestamp.Unix(),
			Additions:    gc.Stats.Additions,
			Deletions:    gc.Stats.Deletions,
			ChangedFiles: 0, // GitHub API doesn't provide this in basic commit endpoint
		})
	}
	return commits
}

// FetchRepositories lists repositories for a user
func (c *Client) FetchRepositories(ctx context.Context, username, pat string) ([]config.Repository, error) {
	var repos []config.Repository
	page := 1
	perPage := 100

	for {
		select {
		case <-ctx.Done():
			return repos, ctx.Err()
		default:
		}

		url := fmt.Sprintf("%s/users/%s/repos?page=%d&per_page=%d",
			githubAPIURL, username, page, perPage)

		var pageRepos []githubRepo
		resp, err := c.doRequestWithPAT(ctx, url, pat)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(resp, &pageRepos); err != nil {
			return nil, fmt.Errorf("failed to parse repositories: %w", err)
		}

		if len(pageRepos) == 0 {
			break
		}

		for _, repo := range pageRepos {
			repos = append(repos, config.Repository{
				Owner: repo.Owner.Login,
				Name:  repo.Name,
			})
		}

		if len(pageRepos) < perPage {
			break
		}

		page++
	}

	return repos, nil
}

// ValidateCredentials checks if GitHub credentials are valid
func (c *Client) ValidateCredentials(username, pat string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	url := fmt.Sprintf("%s/user", githubAPIURL)
	_, err := c.doRequestWithPAT(ctx, url, pat)
	return err
}

// doRequest makes an authenticated HTTP request with retry logic
func (c *Client) doRequest(ctx context.Context, url string) ([]byte, error) {
	return c.doRequestWithPAT(ctx, url, c.pat)
}

// doRequestWithPAT makes an HTTP request with a specific PAT
func (c *Client) doRequestWithPAT(ctx context.Context, url, pat string) ([]byte, error) {
	var lastErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return nil, err
		}

		if pat != "" {
			req.Header.Set("Authorization", fmt.Sprintf("token %s", pat))
		}
		req.Header.Set("Accept", "application/vnd.github.v3+json")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = err
			time.Sleep(baseDelay * time.Duration(attempt+1))
			continue
		}
		defer resp.Body.Close()

		// Check for rate limiting
		if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusTooManyRequests {
			resetTime := c.parseRateLimitReset(resp.Header)
			waitTime := time.Until(resetTime)
			if waitTime > 0 && waitTime < 5*time.Minute {
				time.Sleep(waitTime)
				continue
			}
			lastErr = fmt.Errorf("rate limited, reset at %v", resetTime)
			time.Sleep(baseDelay * time.Duration(attempt+1))
			continue
		}

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			lastErr = fmt.Errorf("GitHub API error: %s - %s", resp.Status, string(body))
			time.Sleep(baseDelay * time.Duration(attempt+1))
			continue
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			lastErr = err
			continue
		}

		return body, nil
	}

	return nil, fmt.Errorf("failed after %d attempts: %w", maxRetries, lastErr)
}

// parseRateLimitReset extracts the reset time from response headers
func (c *Client) parseRateLimitReset(headers http.Header) time.Time {
	resetStr := headers.Get("X-RateLimit-Reset")
	if resetStr == "" {
		return time.Now().Add(60 * time.Second)
	}

	resetUnix, err := strconv.ParseInt(resetStr, 10, 64)
	if err != nil {
		return time.Now().Add(60 * time.Second)
	}

	return time.Unix(resetUnix, 0)
}

// ClearCache clears the in-memory cache
func (c *Client) ClearCache() {
	c.cache = &sync.Map{}
}
