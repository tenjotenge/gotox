# UI Specification

## Overview

The UI layer (`internal/ui`) provides a Fyne-based desktop interface for gotox. This document specifies the current UI implementation and future enhancements.

## Current Implementation (Phase 1)

### Window Properties
- **Title**: "gotox"
- **Default Size**: 800x600 pixels
- **Resizable**: Yes
- **Icon**: None (to be added later)

### Main Window Layout

```
┌─────────────────────────────────────────────────────────┐
│  gotox                                                  │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  Analyzes GitHub commit history and generates          │
│  synthetic development schedules                       │
│                                                         │
│  ─────────────────────────────────────────────────────  │
│                                                         │
│  GitHub username: [placeholder]                         │
│  PAT: [placeholder]                                     │
│  Repo list: [placeholder]                               │
│  Year range: [placeholder]                              │
│  Seed: [placeholder]                                    │
│                                                         │
│  ─────────────────────────────────────────────────────  │
│                                                         │
│  [ Generate ] (disabled)                                │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

### Components

#### Description Label
- **Type**: Static text label
- **Content**: "Analyzes GitHub commit history and generates synthetic development schedules"
- **Wrapping**: Word wrap enabled
- **Purpose**: Explain application purpose to users

#### Placeholder Labels
- **GitHub username**: Static label showing "[placeholder]"
- **PAT**: Static label showing "[placeholder]"
- **Repo list**: Static label showing "[placeholder]"
- **Year range**: Static label showing "[placeholder]"
- **Seed**: Static label showing "[placeholder]"

**Note**: These are currently static labels. In Phase 2, these will be replaced with input fields.

#### Generate Button
- **Type**: Button
- **Label**: "Generate"
- **State**: Disabled
- **Action**: None (to be implemented in Phase 2)

## Future UI Phases

### Phase 2: Input Forms
Replace placeholder labels with actual input fields:

#### GitHub Configuration Section
- **Username**: Text entry field
- **PAT**: Password entry field (masked)
- **Repositories**: Multi-select list or text area (one per line)

#### Analysis Configuration Section
- **Year Range**: Two number entry fields (start year, end year)
- **Seed**: Number entry field

#### Generate Button
- Enable when all required fields are filled
- Show loading state during generation
- Disable during generation

### Phase 3: Progress Display
Add progress indicators during generation:

#### Progress Bar
- Shows overall generation progress
- Updates as each phase completes

#### Status Label
- Shows current phase (e.g., "Fetching commits...", "Analyzing patterns...")
- Shows percentage complete

#### Log Output
- Scrollable text area showing detailed progress
- Timestamps for each log entry
- Color-coded for info/warning/error

### Phase 4: Results Display
Add visualization of results:

#### Tabbed Interface
- **Analysis Tab**: Developer profiles, repository patterns
- **Schedule Tab**: Task timeline, Gantt chart
- **Export Tab**: Export options (JSON, CSV, Markdown)

#### Analysis Visualizations
- Bar charts for commit frequency
- Heat maps for active hours/days
- Pie charts for file type distribution

#### Schedule Visualizations
- Gantt chart showing task timeline
- Task list with details
- Resource allocation view

#### Export Options
- Download buttons for different formats
- Preview of export content
- Copy to clipboard option

### Phase 5: Execution UI (Optional)
Add controls for schedule execution:

#### Execution Controls
- "Execute Schedule" button
- "Dry Run" checkbox
- "Rollback" button

#### Execution Progress
- Task-by-task progress
- Real-time status updates
- Success/failure indicators

#### Execution Results
- Summary of executed tasks
- List of created commits
- Error details if any

## UI Components Library

### Fyne Widgets Used
- `widget.Label`: Static text
- `widget.Entry`: Text input
- `widget.PasswordEntry`: Masked input
- `widget.Button`: Action buttons
- `widget.ProgressBar`: Progress indication
- `widget.Select`: Dropdown selection
- `widget.Check`: Checkbox
- `widget.TextGrid`: Multi-line text
- `container.VBox`: Vertical layout
- `container.HBox`: Horizontal layout
- `container.Padded`: Padding wrapper

### Custom Components (Future)
- `RepositorySelector`: Custom repository selection widget
- `GanttChart`: Custom schedule visualization
- `ProfileCard`: Developer profile display
- `LogViewer`: Formatted log display

## Layout Guidelines

### Spacing
- Use consistent padding (16px default)
- Add spacing between sections (24px)
- Maintain alignment of labels and inputs

### Typography
- Use Fyne's default theme fonts
- Bold for section headers
- Monospace for code/data
- Italic for placeholders

### Colors
- Use Fyne's default color scheme
- Green for success states
- Red for error states
- Yellow for warnings
- Blue for information

## Responsiveness
- Window should be resizable
- Layout should adapt to window size
- Minimum size: 600x400
- Scrollbars for overflow content

## Accessibility
- Keyboard navigation support
- Screen reader compatible labels
- High contrast mode support
- Sufficient touch targets (44px minimum)

## Implementation Notes

### State Management
- UI state should be separate from business logic
- Use callbacks to communicate with backend
- Update UI asynchronously for long operations

### Error Display
- Show user-friendly error messages
- Provide actionable error details
- Offer retry options where applicable

### Validation
- Validate inputs before generation
- Show inline validation errors
- Disable generate button until valid

## Dependencies

- `fyne.io/fyne/v2`: UI framework
- `internal/config`: Configuration data
- `internal/util`: Logging and error handling

## Future Enhancements

- Dark mode support
- Custom themes
- Keyboard shortcuts
- Undo/redo functionality
- Save/load configurations
- Recent configurations list
- Settings/preferences dialog
