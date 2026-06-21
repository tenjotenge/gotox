# Schedule Model

## Overview

The schedule model (`internal/schedule`) is responsible for generating synthetic development schedules based on commit history analysis. This module implements the `Generator` interface.

## Interface Definition

```go
type Generator interface {
    GenerateSchedule(analysis *analysis.AnalysisResult, seed int64) (*Schedule, error)
    ValidateSchedule(schedule *Schedule) error
    ExportSchedule(schedule *Schedule, format string) (string, error)
}
```

## Data Structures

### Task
```go
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

### Schedule
```go
type Schedule struct {
    Tasks         []Task
    StartDate     string
    EndDate       string
    TotalDuration int
    Metadata      map[string]interface{}
}
```

## Generation Algorithm

### Phase 1: Task Extraction
- Analyze commit patterns to identify work types
- Group commits into logical tasks
- Estimate task complexity based on:
  - Number of commits
  - Lines changed
  - File types involved
  - Historical duration

### Phase 2: Effort Estimation
- Use developer profiles to estimate individual capacity
- Apply complexity multipliers from configuration
- Account for historical velocity
- Add buffer percentage from configuration

### Phase 3: Dependency Resolution
- Identify task dependencies based on:
  - File dependencies
  - Sequential commit patterns
  - Module boundaries
- Build dependency graph
- Topological sort for execution order

### Phase 4: Timeline Construction
- Assign start dates based on dependencies
- Distribute tasks across available developers
- Respect work hours per day from configuration
- Handle parallel execution where possible

### Phase 5: Validation
- Check for circular dependencies
- Verify timeline feasibility
- Ensure no resource conflicts
- Validate total duration

## Random Seed Usage

The random seed ensures reproducible schedule generation:
- Seed is used for random number generator
- Same seed + same analysis = same schedule
- Allows for comparison of different scenarios
- Useful for testing and debugging

## Configuration Parameters

```go
type SchedulingConfig struct {
    WorkHoursPerDay    int     // Default: 8
    TaskComplexity     float64 // Default: 1.0 (multiplier)
    BufferPercentage   float64 // Default: 0.1 (10% buffer)
}
```

### Work Hours Per Day
- Defines available work hours per developer per day
- Used for timeline calculation
- Affects total schedule duration

### Task Complexity
- Multiplier for effort estimation
- Higher values = more conservative estimates
- Can be adjusted based on project complexity

### Buffer Percentage
- Extra time added to estimates
- Accounts for uncertainty
- Default 10% can be increased for riskier projects

## Validation Rules

### Schedule Validation
- No circular dependencies
- All dependencies reference valid tasks
- Estimated hours are positive
- Start date is before end date
- Total duration matches task sum

### Task Validation
- Task ID is unique
- Title is not empty
- Estimated hours > 0
- Dependencies are valid task IDs
- Assignee exists in developer profiles

## Export Formats

### JSON
```json
{
  "tasks": [...],
  "start_date": "2024-01-01",
  "end_date": "2024-12-31",
  "total_duration": 250,
  "metadata": {}
}
```

### CSV
```csv
id,title,description,estimated_hours,dependencies,assignee,priority,start_date
TASK-001,Feature A,Implement feature A,8.0,TASK-000,developer1,1,2024-01-01
```

### Markdown
```markdown
# Development Schedule

## Task: TASK-001
- **Title**: Feature A
- **Assignee**: developer1
- **Estimated Hours**: 8.0
- **Dependencies**: TASK-000
- **Priority**: 1
```

## Implementation Notes

### Dependency Graph
- Use adjacency list representation
- Implement topological sort
- Detect cycles during construction
- Allow for parallel independent tasks

### Effort Estimation Formula
```
EstimatedHours = BaseHours × ComplexityMultiplier × (1 + BufferPercentage)

BaseHours = (Commits × AvgCommitTime) + (LinesChanged / LinesPerHour)
```

### Timeline Algorithm
- Greedy assignment to available developers
- Respect dependency constraints
- Minimize total duration
- Balance workload across developers

## Testing Considerations

### Unit Tests
- Test generation with known inputs
- Test validation logic
- Test export formats
- Test dependency resolution

### Integration Tests
- Test with real analysis results
- Test reproducibility with same seed
- Test different configuration values
- Test edge cases (empty tasks, single task)

## Dependencies

- `internal/analysis` for analysis results
- `internal/util` for logging and error handling
- Random number generator (seeded)
- Date/time handling

## Future Enhancements

- Support for multiple export formats
- Schedule optimization algorithms
- Resource leveling
- Critical path analysis
- What-if scenario analysis
- Machine learning for better estimation
