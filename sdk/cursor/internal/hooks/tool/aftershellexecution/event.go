package aftershellexecution

import (
	"context"

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
}

// EventName returns the canonical hook event name.
func (Event) EventName() string { return event.AfterShellExecution }

// DurationMillis returns the execution duration in milliseconds.
// Prefer this helper over reading Duration or DurationMs directly: Cursor
// Hooks docs use `duration`, and DurationMillis falls back to `duration_ms`
// when `duration` is zero so alternate wire forms still decode.
func (e Event) DurationMillis() int64 {
	if e.Duration != 0 {
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
