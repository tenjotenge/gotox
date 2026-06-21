package ui

import (
	"context"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"gotox/internal/config"
	"gotox/internal/pipeline"
)

// App represents the gotox application
type App struct {
	fyneApp  fyne.App
	window   fyne.Window
	pipeline *pipeline.Pipeline
}

// NewApp creates a new gotox application instance
func NewApp() *App {
	a := app.New()
	a.SetIcon(nil) // Icon can be added later
	
	w := a.NewWindow("gotox")
	w.Resize(fyne.NewSize(800, 600))
	
	return &App{
		fyneApp:  a,
		window:   w,
		pipeline: nil,
	}
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
	// Description label
	description := widget.NewLabel(
		"Analyzes GitHub commit history and generates synthetic development schedules",
	)
	description.Wrapping = fyne.TextWrapWord
	
	// Placeholder labels (no input fields yet)
	usernameLabel := widget.NewLabel("GitHub username: [placeholder]")
	patLabel := widget.NewLabel("PAT: [placeholder]")
	repoListLabel := widget.NewLabel("Repo list: [placeholder]")
	yearRangeLabel := widget.NewLabel("Year range: [placeholder]")
	seedLabel := widget.NewLabel("Seed: [placeholder]")
	
	// Generate button - enabled to call pipeline
	generateButton := widget.NewButton("Generate", func() {
		a.onGenerate()
	})
	
	// Layout
	content := container.NewVBox(
		description,
		widget.NewSeparator(),
		usernameLabel,
		patLabel,
		repoListLabel,
		yearRangeLabel,
		seedLabel,
		widget.NewSeparator(),
		generateButton,
	)
	
	a.window.SetContent(container.NewPadded(content))
}

// onGenerate handles the Generate button click
func (a *App) onGenerate() {
	if a.pipeline == nil {
		return
	}

	// Use default configuration for now
	cfg := config.DefaultConfig()

	// Add progress handler to display events
	a.pipeline.AddProgressHandler(func(event pipeline.ProgressEvent) {
		// In future phases, this will update UI with progress
		// For now, just log the event
	})

	// Run pipeline in background
	go func() {
		ctx := context.Background()
		result, err := a.pipeline.Run(ctx, *cfg)
		if err != nil {
			// Handle error - in future phases, display in UI
			return
		}
		// Handle success - in future phases, display results in UI
		_ = result
	}()
}
