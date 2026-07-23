package beforeshellexecution

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/event"
)

// Event is the beforeShellExecution hook event.
type Event struct {
	event.Envelope
	// Command is the shell command about to run.
	Command string `json:"command"`
	// Sandbox reports whether the command runs in a sandbox.
	Sandbox bool `json:"sandbox"`
}

// EventName returns the canonical hook event name.
func (Event) EventName() string { return event.BeforeShellExecution }

// Register registers this hook event decoder on c.
func Register(c *hookkit.Codec) {
	c.Register(event.BeforeShellExecution, hookkit.EventDecoder[Event](c))
}
