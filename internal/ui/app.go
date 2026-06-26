package ui

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"gotox/internal/config"
	"gotox/internal/pipeline"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// App represents the gotox application
type App struct {
	fyneApp    fyne.App
	window     fyne.Window
	pipeline   *pipeline.Pipeline
	isRunning  bool
	statusLock sync.Mutex

	// UI fields
	usernameEntry     *widget.Entry
	patEntry          *widget.Entry
	repoList          *widget.List
	repoEntry         *widget.Entry
	repoItems         []string
	selectedRepoIndex int
	startYearEntry    *widget.Entry
	endYearEntry      *widget.Entry
	seedEntry         *widget.Entry
	workHoursEntry    *widget.Entry
	complexityEntry   *widget.Entry
	bufferEntry       *widget.Entry
	statusLabel       *widget.Label
	progressBar       *widget.ProgressBar
	logOutput         *widget.TextGrid
}

// NewApp creates a new gotox application instance
func NewApp() *App {
	a := app.New()
	a.SetIcon(nil)

	w := a.NewWindow("gotox")
	w.Resize(fyne.NewSize(900, 700))

	appInstance := &App{
		fyneApp:   a,
		window:    w,
		pipeline:  nil,
		isRunning: false,
	}

	return appInstance
}

// SetPipeline sets the pipeline for the application
func (a *App) SetPipeline(p *pipeline.Pipeline) {
	a.pipeline = p
}

// Run starts the application and shows the main window
func (a *App) Run() {
	a.buildUI()
	a.window.ShowAndRun()
}

// buildUI constructs the user interface
func (a *App) buildUI() {
	// Description
	description := widget.NewLabel("Analyzes GitHub commit history and generates synthetic development schedules")
	description.Wrapping = fyne.TextWrapWord

	// GitHub Configuration Section
	githubSection := a.buildGitHubSection()

	// Analysis Configuration Section
	analysisSection := a.buildAnalysisSection()

	// Scheduling Configuration Section
	schedulingSection := a.buildSchedulingSection()

	// Status and Progress Section
	statusSection := a.buildStatusSection()

	// Generate Button
	generateButton := widget.NewButton("Generate", func() {
		a.onGenerate()
	})
	generateButton.Importance = widget.HighImportance

	// Layout
	content := container.NewVBox(
		description,
		widget.NewSeparator(),
		githubSection,
		widget.NewSeparator(),
		analysisSection,
		widget.NewSeparator(),
		schedulingSection,
		widget.NewSeparator(),
		statusSection,
		widget.NewSeparator(),
		generateButton,
	)

	scrollContainer := container.NewScroll(content)
	scrollContainer.SetMinSize(fyne.NewSize(800, 600))

	a.window.SetContent(container.NewPadded(scrollContainer))
}

// buildGitHubSection creates the GitHub configuration section
func (a *App) buildGitHubSection() *widget.Card {
	// Username field
	usernameEntry := widget.NewEntry()
	usernameEntry.SetPlaceHolder("GitHub username")

	// PAT field (masked)
	patEntry := widget.NewEntry()
	patEntry.SetPlaceHolder("GitHub Personal Access Token (optional)")
	patEntry.Password = true

	// Repository list
	repoList := widget.NewList(
		func() int { return len(a.repoItems) },
		func() fyne.CanvasObject { return widget.NewLabel("repository") },
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(a.repoItems[i])
		},
	)
	repoList.OnSelected = func(id widget.ListItemID) {
		a.selectedRepoIndex = id
	}

	// Add repository entry and button
	repoEntry := widget.NewEntry()
	repoEntry.SetPlaceHolder("owner/repo (e.g., octocat/hello-world)")

	addRepoButton := widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
		repo := repoEntry.Text
		if repo != "" {
			a.repoItems = append(a.repoItems, repo)
			repoList.Refresh()
			repoEntry.SetText("")
		}
	})

	removeRepoButton := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
		if a.selectedRepoIndex >= 0 && a.selectedRepoIndex < len(a.repoItems) {
			a.repoItems = append(a.repoItems[:a.selectedRepoIndex], a.repoItems[a.selectedRepoIndex+1:]...)
			a.selectedRepoIndex = -1
			repoList.Refresh()
		}
	})

	repoButtonContainer := container.NewHBox(addRepoButton, removeRepoButton)

	// Store references for later access
	a.usernameEntry = usernameEntry
	a.patEntry = patEntry
	a.repoList = repoList
	a.repoEntry = repoEntry
	a.repoItems = make([]string, 0)
	a.selectedRepoIndex = -1

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "GitHub Username", Widget: usernameEntry},
			{Text: "Personal Access Token", Widget: patEntry},
		},
	}

	repoContainer := container.NewVBox(
		widget.NewLabel("Target Repositories:"),
		container.NewBorder(nil, repoButtonContainer, nil, nil, repoEntry),
		repoList,
	)

	section := container.NewVBox(
		widget.NewLabelWithStyle("GitHub Configuration", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		form,
		widget.NewSeparator(),
		repoContainer,
	)

	return widget.NewCard("", "", section)
}

// buildAnalysisSection creates the analysis configuration section
func (a *App) buildAnalysisSection() *widget.Card {
	// Start year
	startYearEntry := widget.NewEntry()
	startYearEntry.SetPlaceHolder("2020")
	startYearEntry.Validator = func(s string) error {
		if s == "" {
			return nil
		}
		_, err := strconv.Atoi(s)
		return err
	}

	// End year
	endYearEntry := widget.NewEntry()
	endYearEntry.SetPlaceHolder("2024")
	endYearEntry.Validator = func(s string) error {
		if s == "" {
			return nil
		}
		_, err := strconv.Atoi(s)
		return err
	}

	// Random seed
	seedEntry := widget.NewEntry()
	seedEntry.SetPlaceHolder("42")
	seedEntry.Validator = func(s string) error {
		if s == "" {
			return nil
		}
		_, err := strconv.ParseInt(s, 10, 64)
		return err
	}

	// Store references
	a.startYearEntry = startYearEntry
	a.endYearEntry = endYearEntry
	a.seedEntry = seedEntry

	yearContainer := container.NewHBox(
		widget.NewLabel("Start Year:"),
		startYearEntry,
		widget.NewLabel("End Year:"),
		endYearEntry,
	)

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Year Range", Widget: yearContainer},
			{Text: "Random Seed", Widget: seedEntry},
		},
	}

	section := container.NewVBox(
		widget.NewLabelWithStyle("Analysis Configuration", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		form,
	)

	return widget.NewCard("", "", section)
}

// buildSchedulingSection creates the scheduling configuration section
func (a *App) buildSchedulingSection() *widget.Card {
	// Work hours per day
	workHoursEntry := widget.NewEntry()
	workHoursEntry.SetPlaceHolder("8")
	workHoursEntry.Validator = func(s string) error {
		if s == "" {
			return nil
		}
		_, err := strconv.Atoi(s)
		return err
	}

	// Task complexity
	complexityEntry := widget.NewEntry()
	complexityEntry.SetPlaceHolder("1.0")
	complexityEntry.Validator = func(s string) error {
		if s == "" {
			return nil
		}
		_, err := strconv.ParseFloat(s, 64)
		return err
	}

	// Buffer percentage
	bufferEntry := widget.NewEntry()
	bufferEntry.SetPlaceHolder("0.1 (10%)")
	bufferEntry.Validator = func(s string) error {
		if s == "" {
			return nil
		}
		_, err := strconv.ParseFloat(s, 64)
		return err
	}

	// Store references
	a.workHoursEntry = workHoursEntry
	a.complexityEntry = complexityEntry
	a.bufferEntry = bufferEntry

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Work Hours/Day", Widget: workHoursEntry},
			{Text: "Task Complexity", Widget: complexityEntry},
			{Text: "Buffer %", Widget: bufferEntry},
		},
	}

	section := container.NewVBox(
		widget.NewLabelWithStyle("Scheduling Configuration", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		form,
	)

	return widget.NewCard("", "", section)
}

// buildStatusSection creates the status and progress section
func (a *App) buildStatusSection() *widget.Card {
	// Status label
	statusLabel := widget.NewLabel("Ready")
	statusLabel.Importance = widget.MediumImportance

	// Progress bar
	progressBar := widget.NewProgressBar()
	progressBar.Min = 0.0
	progressBar.Max = 1.0
	progressBar.Value = 0.0

	// Log output - use Entry for selectable text
	logOutput := widget.NewEntry()
	logOutput.MultiLine = true
	logOutput.ReadOnly = true
	logOutput.SetPlaceHolder("Log output will appear here...")

	// Store references
	a.statusLabel = statusLabel
	a.progressBar = progressBar
	a.logOutput = logOutput

	section := container.NewVBox(
		widget.NewLabelWithStyle("Status", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		statusLabel,
		progressBar,
		widget.NewLabel("Log Output:"),
		logOutput,
	)

	return widget.NewCard("", "", section)
}

// onGenerate handles the Generate button click
func (a *App) onGenerate() {
	if a.pipeline == nil {
		dialog.ShowError(errors.New("pipeline not initialized"), a.window)
		return
	}

	// Lock to prevent concurrent runs
	a.statusLock.Lock()
	if a.isRunning {
		a.statusLock.Unlock()
		return
	}
	a.isRunning = true
	a.statusLock.Unlock()

	// Clear previous log
	a.logOutput.SetText("")

	// Read configuration from UI
	cfg, err := a.readConfig()
	if err != nil {
		a.showError("Configuration Error", err)
		a.setRunning(false)
		return
	}

	// Update UI state
	a.setStatus("Initializing pipeline...", widget.MediumImportance)
	a.setProgress(0.0)
	a.setRunning(true)

	// Add progress handler
	a.pipeline.AddProgressHandler(func(event pipeline.ProgressEvent) {
		a.handleProgressEvent(event)
	})

	// Run pipeline in background
	go func() {
		ctx := context.Background()
		result, err := a.pipeline.Run(ctx, *cfg)

		// Update UI on main thread
		if err != nil {
			a.showError("Pipeline Failed", err)
			a.setStatus("Failed", widget.DangerImportance)
			a.setProgress(0.0)
		} else {
			a.setStatus("Completed successfully", widget.SuccessImportance)
			a.setProgress(1.0)
			a.appendLog("Pipeline completed successfully!")
			a.showResultSummary(result)
		}
		a.setRunning(false)
	}()
}

// readConfig reads configuration from UI fields
func (a *App) readConfig() (*config.Config, error) {
	cfg := config.DefaultConfig()

	// GitHub username
	cfg.GitHub.Username = a.usernameEntry.Text
	if cfg.GitHub.Username == "" {
		return nil, errors.New("GitHub username is required")
	}

	// PAT (optional)
	cfg.GitHub.PAT = a.patEntry.Text

	// Repositories
	for _, repoStr := range a.repoItems {
		parts := strings.Split(repoStr, "/")
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return nil, fmt.Errorf("invalid repository format: %s (expected owner/repo)", repoStr)
		}
		cfg.GitHub.Repositories = append(cfg.GitHub.Repositories, config.Repository{
			Owner: parts[0],
			Name:  parts[1],
		})
	}

	// Start year
	if a.startYearEntry.Text != "" {
		year, err := strconv.Atoi(a.startYearEntry.Text)
		if err != nil {
			return nil, fmt.Errorf("invalid start year: %v", err)
		}
		cfg.Analysis.StartYear = year
	}

	// End year
	if a.endYearEntry.Text != "" {
		year, err := strconv.Atoi(a.endYearEntry.Text)
		if err != nil {
			return nil, fmt.Errorf("invalid end year: %v", err)
		}
		cfg.Analysis.EndYear = year
	}

	// Random seed
	if a.seedEntry.Text != "" {
		seed, err := strconv.ParseInt(a.seedEntry.Text, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid random seed: %v", err)
		}
		cfg.RandomSeed = seed
	}

	// Work hours per day
	if a.workHoursEntry.Text != "" {
		hours, err := strconv.Atoi(a.workHoursEntry.Text)
		if err != nil {
			return nil, fmt.Errorf("invalid work hours: %v", err)
		}
		cfg.Scheduling.WorkHoursPerDay = hours
	}

	// Task complexity
	if a.complexityEntry.Text != "" {
		complexity, err := strconv.ParseFloat(a.complexityEntry.Text, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid task complexity: %v", err)
		}
		cfg.Scheduling.TaskComplexity = complexity
	}

	// Buffer percentage
	if a.bufferEntry.Text != "" {
		buffer, err := strconv.ParseFloat(a.bufferEntry.Text, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid buffer percentage: %v", err)
		}
		cfg.Scheduling.BufferPercentage = buffer
	}

	return cfg, nil
}

// handleProgressEvent handles progress events from the pipeline
func (a *App) handleProgressEvent(event pipeline.ProgressEvent) {
	message := fmt.Sprintf("[%s] %s", event.Phase, event.Message)

	fyne.Do(func() {
		a.appendLog(message)

		switch event.Type {
		case pipeline.EventTypeStarted:
			a.setStatus(fmt.Sprintf("Running: %s", event.Phase), widget.MediumImportance)
		case pipeline.EventTypeCompleted:
			a.setStatus(fmt.Sprintf("Completed: %s", event.Phase), widget.SuccessImportance)
		case pipeline.EventTypeFailed:
			errorMsg := ""
			if err, ok := event.Metadata["error"].(string); ok {
				errorMsg = err
			}
			a.showError(fmt.Sprintf("%s failed", event.Phase), errors.New(errorMsg))
			a.setStatus(fmt.Sprintf("Failed: %s", event.Phase), widget.DangerImportance)
		}
	})
}

// setStatus updates the status label
func (a *App) setStatus(message string, importance widget.Importance) {
	fyne.Do(func() {
		a.statusLabel.SetText(message)
		a.statusLabel.Importance = importance
	})
}

// setProgress updates the progress bar
func (a *App) setProgress(value float64) {
	fyne.Do(func() {
		a.progressBar.SetValue(value)
	})
}

// setRunning updates the UI state for running/not running
func (a *App) setRunning(running bool) {
	a.statusLock.Lock()
	a.isRunning = running
	a.statusLock.Unlock()

	fyne.Do(func() {
		// Disable/enable all input fields
		widgets := []fyne.Disableable{
			a.usernameEntry,
			a.patEntry,
			a.repoEntry,
			a.startYearEntry,
			a.endYearEntry,
			a.seedEntry,
			a.workHoursEntry,
			a.complexityEntry,
			a.bufferEntry,
		}

		for _, w := range widgets {
			w.Disable()
		}

		if !running {
			for _, w := range widgets {
				w.Enable()
			}
		}
	})
}

// appendLog adds a line to the log output
func (a *App) appendLog(message string) {
	currentText := a.logOutput.Text()
	newText := currentText + message + "\n"
	a.logOutput.SetText(newText)
}

// showError displays an error dialog
func (a *App) showError(title string, err error) {
	fyne.Do(func() {
		dialog.ShowError(err, a.window)
	})
	a.appendLog(fmt.Sprintf("ERROR: %s: %v", title, err))
}

// showResultSummary displays a summary of the pipeline result
func (a *App) showResultSummary(result *pipeline.Result) {
	if result == nil {
		return
	}

	summary := "Pipeline Summary:\n"
	summary += fmt.Sprintf("- Repositories analyzed: %d\n", len(result.RepositoryData))
	summary += fmt.Sprintf("- Analysis completed: %v\n", result.Analysis != nil)
	summary += fmt.Sprintf("- Schedule generated: %v\n", result.GeneratedSchedule != nil)
	summary += fmt.Sprintf("- Execution completed: %v\n", result.Execution != nil)
	summary += fmt.Sprintf("- Verification: %v\n", result.Verification.Success)

	if result.Error != nil {
		summary += fmt.Sprintf("\nError: %v", result.Error)
	}

	fyne.Do(func() {
		dialog.ShowInformation("Generation Complete", summary, a.window)
	})
}
