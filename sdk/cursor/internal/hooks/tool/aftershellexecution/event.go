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
	event.DurationFields
	// Command is the shell command that ran.
	Command string `json:"command"`
	// Output is the terminal output.
	Output string `json:"output"`
	// Sandbox reports whether the command ran in a sandboxed environment.
	Sandbox bool `json:"sandbox"`
}

// EventName returns the canonical hook event name.
func (Event) EventName() string { return event.AfterShellExecution }

// register registers this hook event decoder on c.
func register(c *hookkit.Codec) {
	c.Register(event.AfterShellExecution, func(raw []byte) (hookkit.Event, error) {
		return hookkit.DecodeEvent(c, raw, func(e *Event, raw []byte) error {
			e.CaptureDurationPresent(raw)
			return nil
		})
	})
}
