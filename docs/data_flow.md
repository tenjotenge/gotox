# gotox Data Flow

## Overview

This document describes how data flows through the gotox system from user input to final results.

## High-Level Flow

```
User Input (UI) → Configuration → Pipeline Orchestration → GitHub Ingestion → Analysis → Schedule Generation → Results Display
```

## Detailed Data Flow

### Phase 1: Configuration Input

**Input**: User enters data in Fyne UI
- GitHub username
- Personal Access Token (PAT)
- Repository list
- Year range (start/end)
- Random seed

**Processing**:
1. UI collects input fields
2. Configuration is assembled into `internal/config.Config` struct
3. Configuration validation occurs
4. Valid configuration is passed to subsystems

**Data Structure**:
```go
type Config struct {
    GitHub      GitHubConfig
    Analysis    AnalysisConfig
    Scheduling  SchedulingConfig
    RandomSeed  int64
}
```

### Phase 1.5: Pipeline Orchestration

**Input**: `Config` from configuration layer

**Processing**:
1. Pipeline receives configuration via `Pipeline.Run(ctx, cfg)`
2. Context is propagated to all subsystems
3. Progress handlers are registered for event reporting
4. Phases are executed sequentially with error handling
5. Results are aggregated into `pipeline.Result`

**Output**: `*pipeline.Result` containing:
- RepositoryData (artifacts from ingestion)
- AnalysisResult (artifacts from analysis)
- GeneratedSchedule (artifacts from schedule generation)
- ExecutionResults (artifacts from execution)
- Phase summaries (derived views)

**Error Handling**:
- Context cancellation stops all phases
- Phase failures halt pipeline and return error
- Partial results are preserved in Result struct

### Phase 2: GitHub Ingestion

**Input**: `GitHubConfig` from configuration

**Processing**:
1. `Ingestor.ValidateCredentials()` checks if credentials are valid
2. `Ingestor.FetchRepositories()` lists available repositories
3. `Ingestor.FetchCommits()` retrieves commit history for each repository
4. Commits are filtered by year range
5. Rate limiting is handled automatically

**Output**: `[]github.RepositoryData`
```go
type RepositoryData struct {
    Owner   string
    Name    string
    Commits []Commit
}

type Commit struct {
    SHA         string
    Author      string
    Message     string
    Timestamp   int64
    Additions   int
    Deletions   int
    ChangedFiles int
}
```

**Error Handling**:
- Invalid credentials: Return error, prompt user
- Rate limit exceeded: Wait and retry
- Network errors: Retry with exponential backoff
- Repository not found: Skip and log warning

### Phase 3: Analysis

**Input**: `[]github.RepositoryData` from ingestion

**Processing**:
1. `Analyzer.AnalyzeCommits()` processes all commit data
2. `Analyzer.ExtractDeveloperProfiles()` creates per-developer statistics
3. `Analyzer.ExtractRepositoryPatterns()` identifies repository-level patterns
4. Global patterns are extracted

**Output**: `*analysis.AnalysisResult`
```go
type AnalysisResult struct {
    DeveloperProfiles  map[string]DeveloperProfile
    RepositoryProfiles map[string]RepositoryProfile
    GlobalPatterns     map[string]interface{}
}

type DeveloperProfile struct {
    Username              string
    TotalCommits          int
    AverageCommitsPerDay  float64
    MostActiveHour        int
    MostActiveDay         string
    FileTypeDistribution  map[string]int
}
```

**Analysis Metrics**:
- Commit frequency per developer
- Active hours and days
- File type preferences
- Code churn (additions/deletions)
- Collaboration patterns

### Phase 4: Schedule Generation

**Input**: `*analysis.AnalysisResult` and random seed

**Processing**:
1. `Generator.GenerateSchedule()` creates tasks based on patterns
2. Tasks are estimated using historical data
3. Dependencies are identified
4. Timeline is constructed
5. `Generator.ValidateSchedule()` checks feasibility

**Output**: `*schedule.Schedule`
```go
type Schedule struct {
    Tasks         []Task
    StartDate     string
    EndDate       string
    TotalDuration int
    Metadata      map[string]interface{}
}

type Task struct {
    ID            string
    Title         string
    Description   string
    EstimatedHours float64
    Dependencies  []string
    Assignee      string
    Priority      int
}
```

**Generation Algorithm**:
- Uses random seed for reproducibility
- Estimates effort based on historical patterns
- Assigns tasks to developers based on profiles
- Creates realistic dependencies
- Adds buffer time based on configuration

### Phase 5: Execution (Optional)

**Input**: `*schedule.Schedule`

**Processing**:
1. `Executor.ValidateEnvironment()` checks execution readiness
2. `Executor.ExecuteTask()` runs individual tasks
3. Progress is tracked
4. Rollback available on failure

**Output**: `[]executor.ExecutionResult`
```go
type ExecutionResult struct {
    TaskID   string
    Success  bool
    Output   string
    Error    error
    Duration int64
}
```

### Phase 6: Results Display

**Input**: Analysis results and/or schedule

**Processing**:
1. Results are formatted for display
2. UI updates with visualizations
3. Export options are provided

**Output**: User-visible results in Fyne window

## Data Transformation Summary

| Phase | Input | Output | Transformation |
|-------|-------|--------|----------------|
| Config | User input | Config struct | Validation and struct assembly |
| Pipeline | Config | Result | Orchestration and aggregation |
| Ingestion | Config | RepositoryData | API calls and filtering |
| Analysis | RepositoryData | AnalysisResult | Pattern extraction and statistics |
| Generation | AnalysisResult | Schedule | Task creation and timeline building |
| Execution | Schedule | ExecutionResult | Task execution and tracking |
| Display | Results | UI | Formatting and visualization |

## Error Propagation

Errors flow upward through the system:
1. Low-level errors (network, parsing) are wrapped with context
2. Errors are logged at each level
3. User-facing errors are translated to friendly messages
4. Critical errors halt execution gracefully

## State Management

- Configuration is stateless (passed as parameters)
- Analysis results are immutable after generation
- Schedule can be regenerated with different seeds
- Execution state is tracked separately

## Caching Strategy

- GitHub API responses should be cached to avoid rate limits
- Analysis results can be cached for reuse
- Schedule generation is deterministic given same seed

## Concurrency Considerations

- GitHub ingestion: Concurrent repository fetching
- Analysis: Parallel processing of repositories
- Schedule generation: Sequential (depends on analysis completion)
- Execution: Parallel task execution where dependencies allow
