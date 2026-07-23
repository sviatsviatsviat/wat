package aftershellexecution

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/event"
)

// Event is the afterShellExecution hook event.
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
}

// EventName returns the canonical hook event name.
func (Event) EventName() string { return event.AfterShellExecution }

// DurationMillis returns the execution duration in milliseconds.
func (e Event) DurationMillis() int64 {
	if e.DurationMs != 0 {
		return e.DurationMs
	}
	return e.Duration
}

// register registers this hook event decoder on c.
func register(c *hookkit.Codec) {
	c.Register(event.AfterShellExecution, hookkit.EventDecoder[Event](c))
}
