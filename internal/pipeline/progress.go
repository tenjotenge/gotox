package pipeline

// Phase represents a pipeline phase
type Phase string

const (
	PhaseIngestion   Phase = "ingestion"
	PhaseAnalysis    Phase = "analysis"
	PhaseSchedule    Phase = "schedule"
	PhaseExecution   Phase = "execution"
	PhaseVerification Phase = "verification"
)

// ProgressEvent represents a progress event from the pipeline
type ProgressEvent struct {
	Phase    Phase
	Type     EventType
	Message  string
	Metadata map[string]interface{}
}

// EventType represents the type of progress event
type EventType string

const (
	EventTypeStarted   EventType = "started"
	EventTypeCompleted EventType = "completed"
	EventTypeFailed    EventType = "failed"
)

// ProgressHandler is a callback function for progress events
type ProgressHandler func(event ProgressEvent)

// ProgressReporter handles progress reporting
type ProgressReporter struct {
	handlers []ProgressHandler
}

// NewProgressReporter creates a new progress reporter
func NewProgressReporter() *ProgressReporter {
	return &ProgressReporter{
		handlers: make([]ProgressHandler, 0),
	}
}

// AddHandler adds a progress event handler
func (r *ProgressReporter) AddHandler(handler ProgressHandler) {
	r.handlers = append(r.handlers, handler)
}

// Report reports a progress event to all handlers
func (r *ProgressReporter) Report(event ProgressEvent) {
	for _, handler := range r.handlers {
		handler(event)
	}
}

// PhaseStarted reports that a phase has started
func (r *ProgressReporter) PhaseStarted(phase Phase, message string) {
	r.Report(ProgressEvent{
		Phase:    phase,
		Type:     EventTypeStarted,
		Message:  message,
		Metadata: make(map[string]interface{}),
	})
}

// PhaseCompleted reports that a phase has completed
func (r *ProgressReporter) PhaseCompleted(phase Phase, message string) {
	r.Report(ProgressEvent{
		Phase:    phase,
		Type:     EventTypeCompleted,
		Message:  message,
		Metadata: make(map[string]interface{}),
	})
}

// PhaseFailed reports that a phase has failed
func (r *ProgressReporter) PhaseFailed(phase Phase, message string, err error) {
	r.Report(ProgressEvent{
		Phase:    phase,
		Type:     EventTypeFailed,
		Message:  message,
		Metadata: map[string]interface{}{
			"error": err.Error(),
		},
	})
}
