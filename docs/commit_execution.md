# Commit Execution Module

## Overview

The commit execution module (`internal/executor`) is responsible for executing generated schedules by creating actual git commits. This module implements the `Executor` interface. This is an optional feature for applying generated schedules to real repositories.

## Interface Definition

```go
type Executor interface {
    ExecuteTask(task schedule.Task) (*ExecutionResult, error)
    ExecuteSchedule(schedule *schedule.Schedule) ([]ExecutionResult, error)
    ValidateEnvironment() error
    Rollback(taskID string) error
}
```

## Data Structures

### ExecutionResult
```go
type ExecutionResult struct {
    TaskID   string
    Success  bool
    Output   string
    Error    error
    Duration int64 // milliseconds
}
```

## Execution Model

### Task Execution Flow
1. Validate environment (git installed, repository accessible)
2. Parse task description to determine commit content
3. Create/checkout appropriate branch
4. Make changes (create files, modify files)
5. Stage changes
6. Create commit with message from task
7. Return execution result

### Schedule Execution Flow
1. Validate environment
2. Sort tasks by dependencies
3. Execute tasks in dependency order
4. Track progress and results
5. Handle failures with rollback option
6. Return aggregate results

## Environment Validation

### Required Checks
- Git is installed and accessible
- Current directory is a git repository
- Git remote is configured
- Write permissions are available
- No uncommitted changes (or handle them)

### Validation Process
```go
func (e *Executor) ValidateEnvironment() error {
    // Check git installation
    // Check repository status
    // Check remote configuration
    // Check write permissions
    return nil
}
```

## Task Execution

### Task Interpretation
Tasks from the schedule need to be interpreted as git operations:
- Task title → commit message subject
- Task description → commit message body
- Task metadata → file changes

### File Generation Strategy
- Use templates for common file types
- Generate synthetic code based on task description
- Create realistic file structures
- Respect project conventions

### Commit Creation
```go
func (e *Executor) ExecuteTask(task schedule.Task) (*ExecutionResult, error) {
    // 1. Parse task
    // 2. Create branch if needed
    // 3. Generate files
    // 4. Stage changes
    // 5. Create commit
    // 6. Return result
}
```

## Rollback Mechanism

### Rollback Strategy
- Track commits created by executor
- On failure, revert to pre-execution state
- Use git reset or revert
- Clean up temporary branches

### Rollback Process
```go
func (e *Executor) Rollback(taskID string) error {
    // Find commit(s) for task
    // Reset or revert commits
    // Clean up branches
    return nil
}
```

## Safety Features

### Dry Run Mode
- Execute without making actual changes
- Show what would be done
- Validate before real execution

### Confirmation Prompts
- Require user confirmation before execution
- Show summary of changes
- Allow cancellation

### Backup Creation
- Create backup branch before execution
- Tag pre-execution state
- Enable easy recovery

## Error Handling

### Execution Errors
- File creation failures
- Git operation failures
- Permission errors
- Merge conflicts

### Error Recovery
- Log detailed error information
- Suggest corrective actions
- Provide rollback option
- Continue with remaining tasks if possible

## Implementation Notes

### Git Operations
- Use `exec.Command` to run git commands
- Parse git output for status
- Handle git exit codes
- Use appropriate git flags

### Branch Strategy
- Create feature branch per task
- Or use single branch for all tasks
- Merge to main after successful execution
- Clean up temporary branches

### Commit Message Format
```
<task title>

<task description>

Task ID: <task.id>
Priority: <task.priority>
Assignee: <task.assignee>
```

## Testing Considerations

### Unit Tests
- Mock git commands
- Test task parsing
- Test error handling
- Test rollback logic

### Integration Tests
- Use test repository
- Execute real git operations
- Test rollback
- Test concurrent execution

### Safety Tests
- Test with protected branches
- Test with insufficient permissions
- Test with merge conflicts

## Dependencies

- `internal/schedule` for schedule data
- `internal/gitutil` for git operations
- `internal/util` for logging and error handling
- Git command-line tool

## Security Considerations

### Code Injection
- Validate task descriptions
- Sanitize file paths
- Prevent command injection in git commands

### Repository Safety
- Never execute on production without confirmation
- Require explicit opt-in
- Create backups before execution

### Access Control
- Respect git permissions
- Don't bypass security checks
- Validate PAT scope

## Future Enhancements

- Support for multiple git providers
- Parallel task execution
- Progress visualization
- Execution simulation
- Integration with CI/CD
- Custom commit templates
