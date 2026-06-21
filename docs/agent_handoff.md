# Agent Handoff Guide

## Overview

This document provides guidance for AI agents (Cursor, Devin, or others) working on the gotox project. It explains how to interact with the codebase, understand the architecture, and make effective contributions without requiring prior conversation context.

## Project Context

gotox is a Go + Fyne desktop application that:
- Analyzes GitHub commit history
- Extracts development patterns
- Generates synthetic development schedules
- Optionally executes schedules by creating git commits

The project is designed for multi-agent development over time, with clear module boundaries and comprehensive documentation.

## Quick Start for New Agents

### 1. Understand the Architecture
Read `docs/architecture.md` to understand:
- System architecture and module boundaries
- Design principles and data flow
- Concurrency model and error handling

### 2. Review Data Structures
Read `docs/data_flow.md` to understand:
- How data flows through the system
- Data structures at each phase
- Error propagation and state management

### 3. Check Module-Specific Docs
Review relevant module documentation:
- `docs/github_ingestion.md` - GitHub API integration
- `docs/schedule_model.md` - Schedule generation
- `docs/commit_execution.md` - Commit execution
- `docs/ui_spec.md` - UI implementation

### 4. Build and Run
```bash
go mod tidy
go build ./...
go run cmd/gotox
```

## Agent Roles and Responsibilities

### Cursor Agent
**Primary Focus**: Incremental development, UI work, bug fixes

**Typical Tasks**:
- Implement UI features (forms, buttons, visualizations)
- Add input validation
- Fix UI bugs and layout issues
- Improve user experience
- Add accessibility features
- Implement progress indicators

**Working Pattern**:
- Make small, incremental changes
- Test UI changes frequently
- Focus on user-facing features
- Maintain backward compatibility

**Key Files**:
- `internal/ui/app.go` - Main UI implementation
- `docs/ui_spec.md` - UI specifications

### Devin Agent
**Primary Focus**: Full-module generation, backend logic, complex features

**Typical Tasks**:
- Implement complete modules (GitHub ingestion, analysis, scheduling)
- Add new interfaces and implementations
- Implement complex algorithms
- Add comprehensive testing
- Handle edge cases and error scenarios

**Working Pattern**:
- Implement complete features end-to-end
- Add thorough documentation
- Include comprehensive tests
- Consider performance and scalability

**Key Files**:
- Module-specific implementation files
- Interface definitions
- Module-specific documentation

## Module Implementation Guidelines

### When Implementing a New Module

1. **Read the Interface Definition**
   - Find the interface in the appropriate `internal/*` package
   - Understand all methods and their contracts
   - Review data structures used

2. **Read Module Documentation**
   - Review the corresponding doc in `docs/`
   - Understand requirements and constraints
   - Check implementation notes

3. **Implement the Interface**
   - Create a struct that implements the interface
   - Implement all required methods
   - Follow Go best practices and idioms

4. **Add Error Handling**
   - Use `internal/util.ErrorHelper` for error wrapping
   - Log errors at appropriate levels
   - Provide meaningful error messages

5. **Add Logging**
   - Use `internal/util.Logger` for logging
   - Log at appropriate levels (Debug, Info, Warn, Error)
   - Include context in log messages

6. **Write Tests**
   - Add unit tests for all public methods
   - Mock external dependencies
   - Test error scenarios
   - Aim for high coverage

7. **Update Documentation**
   - Update relevant docs if behavior changes
   - Add examples if helpful
   - Document any deviations from spec

### When Modifying Existing Code

1. **Understand the Context**
   - Read the module documentation
   - Understand the data flow
   - Check for dependencies

2. **Make Minimal Changes**
   - Change only what's necessary
   - Maintain backward compatibility
   - Don't refactor without reason

3. **Test Thoroughly**
   - Run existing tests
   - Add new tests for changes
   - Test integration with other modules

4. **Update Documentation**
   - Update relevant docs
   - Note breaking changes
   - Update examples if needed

## Code Style Guidelines

### Go Conventions
- Follow standard Go formatting (`go fmt`)
- Use meaningful variable names
- Keep functions focused and small
- Use interfaces for abstraction
- Prefer composition over inheritance

### Error Handling
- Always handle errors
- Wrap errors with context
- Don't ignore errors
- Log errors before returning

### Documentation
- Add package comments
- Document exported functions
- Document complex algorithms
- Keep comments up to date

### Testing
- Write tests for all exported functions
- Use table-driven tests for multiple cases
- Mock external dependencies
- Test error scenarios

## Common Patterns

### Dependency Injection
```go
type MyService struct {
    logger *util.Logger
    helper *util.ErrorHelper
}

func NewMyService(logger *util.Logger) *MyService {
    return &MyService{
        logger: logger,
        helper: util.NewErrorHelper(),
    }
}
```

### Error Handling
```go
func (s *MyService) DoSomething() error {
    result, err := s.someOperation()
    if err != nil {
        return s.helper.Wrap(err, "failed to do something")
    }
    return nil
}
```

### Logging
```go
func (s *MyService) DoSomething() {
    s.logger.Info("Starting operation")
    // ... do work
    s.logger.Info("Operation completed")
}
```

## Testing Guidelines

### Unit Tests
- Test individual functions in isolation
- Mock external dependencies
- Test both success and error cases
- Use descriptive test names

### Integration Tests
- Test module interactions
- Use real implementations where appropriate
- Test data flow between modules
- Test error propagation

### Running Tests
```bash
go test ./...
go test -v ./...
go test -cover ./...
```

## Build and Verification

### Before Submitting Changes
1. Run `go build ./...` - Ensure code compiles
2. Run `go test ./...` - Ensure tests pass
3. Run `go fmt ./...` - Ensure formatting is correct
4. Run `go vet ./...` - Check for common issues
5. Test the application manually if UI changes

### Common Build Issues
- Missing imports: Run `go mod tidy`
- Type errors: Check interface implementations
- Dependency issues: Check go.mod and go.sum

## Communication Guidelines

### When Asking Questions
- Include relevant code snippets
- Reference specific documentation
- Describe what you've tried
- Explain the expected vs actual behavior

### When Reporting Issues
- Include error messages
- Describe reproduction steps
- Include environment details
- Suggest potential fixes

## Module-Specific Notes

### GitHub Ingestion
- Respect rate limits
- Handle pagination correctly
- Cache responses appropriately
- Test with real GitHub API

### Analysis Engine
- Process data efficiently
- Handle edge cases (empty repos, single commits)
- Validate analysis results
- Consider performance for large datasets

### Schedule Generation
- Ensure reproducibility with seed
- Validate dependencies
- Check timeline feasibility
- Test with various configurations

### Commit Execution
- Implement safety features
- Require user confirmation
- Provide rollback capability
- Test with real git operations

## Getting Help

### Internal Resources
- Check module documentation in `docs/`
- Review interface definitions in `internal/*`
- Check existing implementations for patterns
- Review test files for usage examples

### External Resources
- Go documentation: https://golang.org/doc/
- Fyne documentation: https://docs.fyne.io/
- GitHub API docs: https://docs.github.com/en/rest

## Checklist Before Handoff

- [ ] Code compiles without errors
- [ ] All tests pass
- [ ] Code is properly formatted
- [ ] Documentation is updated
- [ ] Error handling is complete
- [ ] Logging is appropriate
- [ ] Tests cover new functionality
- [ ] No breaking changes (or documented)
- [ ] Module boundaries are respected
- [ ] Dependencies are minimal

## Version Control

### Commit Messages
- Use clear, descriptive messages
- Reference relevant issues if applicable
- Keep messages concise but informative
- Use imperative mood

### Branch Strategy
- Create feature branches for new work
- Keep branches focused on single features
- Merge frequently to avoid conflicts
- Delete merged branches

## Future Work

The project is designed for iterative development. Future phases include:
- Phase 2: Input forms and real generation
- Phase 3: Progress display and visualizations
- Phase 4: Results display and export
- Phase 5: Optional execution features

When implementing future phases, refer to this document and the phase-specific documentation.
