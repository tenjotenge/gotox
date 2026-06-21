# gotox

A desktop application for analyzing GitHub commit history and generating synthetic development schedules.

## Purpose

gotox analyzes historical commit patterns from GitHub repositories to create realistic development schedules. This tool is particularly useful for filling in commit history for projects that cannot be publicly represented, such as contract work under NDA, private client projects, or other confidential work. By analyzing your public GitHub activity, gotox can generate realistic schedules that reflect your actual development patterns and work habits.

## Architecture Overview

gotox follows a modular architecture with clear separation of concerns:

- **UI Layer** (internal/ui): Fyne-based desktop interface
- **Configuration Layer** (internal/config): Application configuration and settings
- **Core Subsystems**: GitHub ingestion, analysis engine, schedule generator, commit executor
- **Utilities Layer** (internal/util): Logging, error handling, and common utilities

The system is designed with comprehensive documentation in the `docs/` folder serving as technical reference.

## Project Structure

```
gotox/
├── cmd/
│   └── gotox/           # Application entry point
│       └── main.go
├── internal/
│   ├── ui/              # Fyne desktop interface
│   ├── config/          # Configuration data structures
│   ├── github/          # GitHub API integration
│   ├── analysis/        # Commit history analysis and session extraction
│   ├── schedule/        # Schedule generation
│   ├── executor/        # Schedule execution
│   ├── gitutil/         # Git utilities
│   ├── verify/          # Data validation
│   └── util/            # Logging and error handling
├── docs/                # Technical documentation
│   ├── architecture.md
│   ├── data_flow.md
│   ├── github_ingestion.md
│   ├── schedule_model.md
│   ├── session_model.md
│   ├── commit_execution.md
│   └── ui_spec.md
├── go.mod
├── go.sum
└── README.md
```

## Development Phases

### Phase 1: Foundation (Current)
- Project scaffold and architecture
- Basic UI shell with placeholder labels
- Configuration structures
- Interface definitions
- Documentation

### Phase 2: Input and Generation
- Real input forms for GitHub credentials
- GitHub API integration
- Commit history fetching
- Basic analysis engine
- Schedule generation

### Phase 3: Visualization
- Progress indicators
- Analysis visualizations
- Schedule timeline display
- Export functionality

### Phase 4: Execution (Optional)
- Schedule execution capabilities
- Git commit generation
- Rollback mechanisms
- Safety features

## Build Instructions

### Prerequisites
- Go 1.26.4 or later
- Git (for execution features)

### Setup
```bash
# Install dependencies
go mod tidy

# Build all packages
go build ./...

# Run the application
go run cmd/gotox
```

### Build Output
The application compiles to an executable that opens a Fyne window titled "gotox".

## Documentation

Technical documentation is available in the `docs/` folder:

- **architecture.md**: System architecture, module boundaries, design principles
- **data_flow.md**: How data flows through the system
- **github_ingestion.md**: GitHub API integration details
- **schedule_model.md**: Schedule generation algorithm
- **session_model.md**: Work session extraction and behavioral profiling
- **commit_execution.md**: Commit execution implementation
- **ui_spec.md**: UI specifications and future phases

## Current Status

**Phase 1 Complete**:
- ✅ Project scaffold with proper structure
- ✅ Configuration data structures
- ✅ Interface definitions for all major subsystems
- ✅ Basic Fyne UI with placeholder labels
- ✅ Logging and error handling utilities
- ✅ Comprehensive documentation

**Next Steps**:
- Implement GitHub API integration
- Add real input forms to UI
- Implement commit history analysis
- Add schedule generation logic

## Dependencies

- **fyne.io/fyne/v2**: Cross-platform GUI framework
- **Go standard library**: For core functionality

## License

[Add license information here]

## Contributing

When contributing to this project:

1. Read the relevant documentation in `docs/`
2. Follow the interface definitions in `internal/*`
3. Maintain module boundaries
4. Add tests for new functionality
5. Update documentation as needed
