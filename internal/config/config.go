package config

// Config holds all application configuration
type Config struct {
	GitHub      GitHubConfig      `json:"github"`
	Analysis    AnalysisConfig    `json:"analysis"`
	Scheduling  SchedulingConfig  `json:"scheduling"`
	RandomSeed  int64             `json:"random_seed"`
}

// GitHubConfig holds GitHub account and repository configuration
type GitHubConfig struct {
	Username   string   `json:"username"`
	PAT        string   `json:"pat"` // Personal Access Token
	Repositories []Repository `json:"repositories"`
}

// Repository represents a GitHub repository to analyze
type Repository struct {
	Owner string `json:"owner"`
	Name  string `json:"name"`
}

// AnalysisConfig holds analysis window and parameters
type AnalysisConfig struct {
	StartYear int `json:"start_year"`
	EndYear   int `json:"end_year"`
}

// SchedulingConfig holds scheduling generation parameters
type SchedulingConfig struct {
	// Placeholder parameters for future scheduling logic
	WorkHoursPerDay    int     `json:"work_hours_per_day"`
	TaskComplexity     float64 `json:"task_complexity"`
	BufferPercentage   float64 `json:"buffer_percentage"`
}

// DefaultConfig returns a configuration with sensible defaults
func DefaultConfig() *Config {
	return &Config{
		GitHub: GitHubConfig{
			Username:   "",
			PAT:        "",
			Repositories: []Repository{},
		},
		Analysis: AnalysisConfig{
			StartYear: 2020,
			EndYear:   2024,
		},
		Scheduling: SchedulingConfig{
			WorkHoursPerDay:   8,
			TaskComplexity:    1.0,
			BufferPercentage:  0.1,
		},
		RandomSeed: 42,
	}
}
