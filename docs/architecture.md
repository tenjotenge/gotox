# gotox Architecture

## Overview

gotox is a modular desktop application built with Go and Fyne that analyzes GitHub commit history and generates synthetic development schedules. The architecture follows a clean separation of concerns with well-defined interfaces between subsystems.

## System Architecture

```
┌─────────────────────────────────────────────────────────┐
│                      UI Layer                           │
│                   (internal/ui)                         │
│  - Fyne-based desktop interface                         │
│  - Configuration input forms                             │
│  - Progress display and results visualization           │
└────────────────────┬────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────┐
│              Configuration Layer                         │
│                (internal/config)                         │
│  - GitHub credentials and repository list                │
│  - Analysis window parameters                            │
│  - Scheduling configuration                              │
└────────────────────┬────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────┐
│              Core Subsystems                            │
│                                                          │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐ │
│  │   GitHub     │  │  Analysis    │  │  Schedule    │ │
│  │  Ingestion   │  │   Engine     │  │  Generator   │ │
│  │              │  │              │  │              │ │
│  │ Fetches      │  │ Extracts     │  │ Creates      │ │
│  │ commits      │  │ patterns     │  │ timelines    │ │
│  └──────────────┘  └──────────────┘  └──────────────┘ │
│                                                          │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐ │
│  │   Executor   │  │  Git Utils   │  │   Verify     │ │
│  │              │  │              │  │              │ │
│  │ Executes     │  │ Parses git   │  │ Validates    │ │
│  │ schedules    │  │ data         │  │ integrity    │ │
│  └──────────────┘  └──────────────┘  └──────────────┘ │
└──────────────────────────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────┐
│              Utilities Layer                             │
│                (internal/util)                          │
│  - Logging infrastructure                                │
│  - Error handling helpers                                │
└──────────────────────────────────────────────────────────┘
```

## Module Boundaries

### internal/ui
**Purpose**: Fyne-based desktop interface
**Responsibilities**:
- Display application window
- Collect user configuration
- Show progress and results
- Handle user interactions

**Dependencies**: internal/config, internal/util

### internal/config
**Purpose**: Configuration data structures
**Responsibilities**:
- Define configuration schemas
- Provide default values
- Validate configuration

**Dependencies**: None (pure data structures)

### internal/github
**Purpose**: GitHub API integration
**Responsibilities**:
- Fetch commit history
- Validate credentials
- List repositories
- Rate limit handling

**Interface**: `Ingestor`
**Dependencies**: internal/config

### internal/analysis
**Purpose**: Commit history pattern analysis
**Responsibilities**:
- Extract developer profiles
- Identify repository patterns
- Analyze commit frequency
- Detect coding patterns

**Interface**: `Analyzer`
**Dependencies**: internal/github

### internal/schedule
**Purpose**: Schedule generation from analysis
**Responsibilities**:
- Create task timelines
- Estimate effort
- Handle dependencies
- Export schedules

**Interface**: `Generator`
**Dependencies**: internal/analysis

### internal/executor
**Purpose**: Schedule execution
**Responsibilities**:
- Execute scheduled tasks
- Track progress
- Handle rollbacks
- Validate environment

**Interface**: `Executor`
**Dependencies**: internal/schedule

### internal/gitutil
**Purpose**: Git-specific utilities
**Responsibilities**:
- Parse commit data
- Analyze diffs
- Git repository operations

**Interface**: `Parser`
**Dependencies**: None

### internal/verify
**Purpose**: Data validation and verification
**Responsibilities**:
- Validate schedules
- Verify commit integrity
- Check data consistency

**Interface**: `Verifier`
**Dependencies**: None

### internal/util
**Purpose**: Shared utilities
**Responsibilities**:
- Logging infrastructure
- Error handling helpers
- Common utilities

**Dependencies**: None

## Design Principles

1. **Interface-Based Design**: All major subsystems expose interfaces for testability and replaceability
2. **Dependency Injection**: Components receive dependencies through constructors
3. **Single Responsibility**: Each module has a clear, focused purpose
4. **Separation of Concerns**: UI, business logic, and data access are separated
5. **Extensibility**: New implementations can be added without modifying existing code

## Data Flow

1. User provides configuration through UI
2. Configuration is validated and passed to subsystems
3. GitHub Ingestor fetches commit data
4. Analysis Engine processes commits and extracts patterns
5. Schedule Generator creates timelines from patterns
6. Executor runs scheduled tasks (optional)
7. Results are displayed in UI

## Concurrency Model

- GitHub ingestion uses concurrent fetching for multiple repositories
- Analysis processes data in parallel where possible
- UI runs on main thread, background work uses goroutines
- Logging is thread-safe

## Error Handling

- Errors are wrapped with context using `internal/util.ErrorHelper`
- Errors are logged before being returned
- UI displays user-friendly error messages
- Critical errors halt execution gracefully
