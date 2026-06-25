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

### Phase 1: Foundation (Complete)
- Project scaffold and architecture
- Configuration structures
- Interface definitions
- Documentation

### Phase 2: UI Implementation (Complete)
- Fully functional Fyne-based desktop interface
- Complete configuration forms for all backend parameters
- GitHub username and PAT input (masked password field)
- Dynamic repository list management (add/remove)
- Analysis year range and random seed configuration
- Scheduling parameters (work hours, complexity, buffer)
- Real-time status reporting and progress tracking
- Error display and result summaries
- Pipeline execution with UI state management

### Phase 3: Backend Integration (Current)
- GitHub API integration with real client
- Commit history fetching and analysis
- Schedule generation with actual algorithms
- Execution and verification phases

### Phase 4: Visualization (Future)
- Analysis visualizations
- Schedule timeline display
- Export functionality

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

## Usage

### Getting Started

1. **Launch the Application**: Run `go run cmd/gotox` or execute the compiled binary
2. **Configure GitHub Settings**:
   - Enter your GitHub username
   - Optionally provide a GitHub Personal Access Token (PAT) for higher API rate limits
   - Add target repositories in `owner/repo` format (e.g., `octocat/hello-world`)
3. **Set Analysis Parameters**:
   - Specify the year range for commit history analysis
   - Set a random seed for reproducible results
4. **Adjust Scheduling** (optional):
   - Configure work hours per day
   - Set task complexity multiplier
   - Adjust buffer percentage for schedule padding
5. **Generate Schedule**:
   - Click the "Generate" button to start the pipeline
   - Monitor progress in the status section
   - View results in the completion dialog

### Configuration via UI

The application is designed to be fully configurable through the UI without requiring environment variables:

- **GitHub Username**: Required field for identifying the account to analyze
- **Personal Access Token**: Optional field for GitHub API authentication
  - Without PAT: Uses stub mode with limited functionality
  - With PAT: Enables real GitHub API integration
- **Target Repositories**: Dynamic list where you can add/remove repositories
- **Analysis Period**: Start and end years for commit history analysis
- **Random Seed**: Controls reproducibility of generated schedules
- **Scheduling Parameters**: Fine-tune schedule generation behavior

### Environment Variables (Optional)

While the UI provides complete configuration, environment variables can still be used:

- `GITHUB_PAT`: GitHub Personal Access Token (overrides UI setting if provided)

### Stub Mode vs. Real Mode

The application operates in two modes:

- **Stub Mode** (default): Uses mock implementations for all backend components. Useful for UI testing and development.
- **Real Mode**: Activated when a valid GitHub PAT is provided. Enables actual GitHub API integration, commit analysis, and schedule generation.

The mode is automatically detected based on whether a PAT is provided via the UI or environment variable.

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

**Phase 1 & 2 Complete**:
- ✅ Project scaffold with proper structure
- ✅ Configuration data structures
- ✅ Interface definitions for all major subsystems
- ✅ Fully functional Fyne UI with real input controls
- ✅ Logging and error handling utilities
- ✅ Comprehensive documentation
- ✅ Complete UI-to-backend pipeline integration
- ✅ Real-time progress reporting and error handling

**Current Phase**: Backend Integration
- GitHub API client implementation
- Analysis engine development
- Schedule generation algorithms

**Next Steps**:
- Complete GitHub API integration
- Implement analysis engine with real algorithms
- Add schedule generation logic
- Implement execution and verification phases

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
