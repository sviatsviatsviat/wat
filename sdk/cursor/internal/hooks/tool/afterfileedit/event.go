package afterfileedit

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/event"
)

// Event is the afterFileEdit hook event.
type Event struct {
	event.Envelope
	// FilePath is the edited file path.
	FilePath string `json:"file_path"`
	// Edits are the applied edits.
	Edits []event.Edit `json:"edits"`
}

// EventName returns the canonical hook event name.
func (Event) EventName() string { return event.AfterFileEdit }

// register registers this hook event decoder on c.
func register(c *hookkit.Codec) {
	c.Register(event.AfterFileEdit, hookkit.EventDecoder[Event](c))
}
