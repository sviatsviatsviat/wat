package aftershellexecution

import (
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

// EventName returns the canonical hook event name.
func (Event) EventName() string { return event.AfterShellExecution }

// DurationMillis returns the execution duration in milliseconds.
// Prefer this helper over reading Duration or DurationMs directly: Cursor
// Hooks docs use `duration`, and DurationMillis falls back to `duration_ms`
// only when `duration` is absent so an explicit `duration: 0` still wins.
func (e Event) DurationMillis() int64 {
	return event.PreferDurationField(e.Duration, e.DurationMs, e.durationPresent)
}

// register registers this hook event decoder on c.
func register(c *hookkit.Codec) {
	c.Register(event.AfterShellExecution, func(raw []byte) (hookkit.Event, error) {
		return hookkit.DecodeEvent(c, raw, func(e *Event, raw []byte) {
			e.durationPresent = hookkit.RawObjectField(raw, "duration") != nil
		})
	})
}
