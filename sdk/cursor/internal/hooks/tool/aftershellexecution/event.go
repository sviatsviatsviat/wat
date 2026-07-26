package aftershellexecution

import (
	"context"
	"encoding/json"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/event"
)

// Event is the afterShellExecution hook event.
//
// Cursor Hooks docs list input fields only (command, output, duration, sandbox);
// the host does not document consumed stdout JSON for this event. Handlers are
// observe-only.
type Event struct {
	event.Envelope
	// Command is the shell command that ran.
	Command string `json:"command"`
	// Output is the terminal output.
	Output string `json:"output"`
	// Duration is the execution duration in milliseconds.
	Duration int64 `json:"duration"`
	// DurationMs is an alternate duration field in milliseconds.
	DurationMs int64 `json:"duration_ms"`
	// Sandbox reports whether the command ran in a sandboxed environment.
	Sandbox bool `json:"sandbox"`

	durationPresent bool
}

// UnmarshalJSON decodes the event and records whether the documented duration
// field was present so DurationMillis can distinguish explicit 0 from absent.
func (e *Event) UnmarshalJSON(data []byte) error {
	type wire struct {
		event.Envelope
		Command    string `json:"command"`
		Output     string `json:"output"`
		Duration   *int64 `json:"duration"`
		DurationMs int64  `json:"duration_ms"`
		Sandbox    bool   `json:"sandbox"`
	}
	var w wire
	if err := json.Unmarshal(data, &w); err != nil {
		return err
	}
	*e = Event{
		Envelope:   w.Envelope,
		Command:    w.Command,
		Output:     w.Output,
		DurationMs: w.DurationMs,
		Sandbox:    w.Sandbox,
	}
	if w.Duration != nil {
		e.durationPresent = true
		e.Duration = *w.Duration
	}
	return nil
}

// EventName returns the canonical hook event name.
func (Event) EventName() string { return event.AfterShellExecution }

// DurationMillis returns the execution duration in milliseconds.
// Prefer this helper over reading Duration or DurationMs directly: Cursor
// Hooks docs use `duration`, and DurationMillis falls back to `duration_ms`
// only when `duration` is absent so an explicit `duration: 0` still wins.
func (e Event) DurationMillis() int64 {
	if e.durationPresent {
		return e.Duration
	}
	return e.DurationMs
}

// register registers this hook event decoder on c.
func register(c *hookkit.Codec) {
	c.Register(event.AfterShellExecution, hookkit.EventDecoder[Event](c))
}

// RegisterHandler registers an AfterShellExecution observe handler on d.
func RegisterHandler(d *hookkit.Dialect, fn func(context.Context, Event) error) {
	if fn == nil {
		return
	}
	register(d.Codec())
	hookkit.RegisterObserve(d, fn)
}
