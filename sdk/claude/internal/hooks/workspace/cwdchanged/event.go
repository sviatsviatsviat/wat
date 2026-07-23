package cwdchanged

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
)

// Event is the CwdChanged hook event.
type Event struct {
	event.Envelope
	// NewCwd is the new working directory.
	NewCwd string `json:"new_cwd"`
	// OldCwd is the previous working directory.
	OldCwd string `json:"old_cwd,omitempty"`
}

// EventName returns the hook event name.
func (Event) EventName() string { return event.CwdChanged }

// Register registers this hook event decoder on c.
func Register(c *hookkit.Codec) {
	c.Register(event.CwdChanged, hookkit.EventDecoder[Event](c))
}
