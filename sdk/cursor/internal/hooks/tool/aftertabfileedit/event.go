package aftertabfileedit

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/event"
)

// Event is the afterTabFileEdit hook event.
type Event struct {
	event.Envelope
	// FilePath is the edited file path.
	FilePath string `json:"file_path"`
	// Edits are the applied Tab edits, including range and line context.
	Edits []event.TabEdit `json:"edits"`
}

// EventName returns the canonical hook event name.
func (Event) EventName() string { return event.AfterTabFileEdit }

// register registers this hook event decoder on c.
func register(c *hookkit.Codec) {
	c.Register(event.AfterTabFileEdit, hookkit.EventDecoder[Event](c))
}
