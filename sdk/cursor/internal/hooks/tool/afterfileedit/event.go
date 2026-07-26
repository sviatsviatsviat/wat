package afterfileedit

import (
	"context"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/event"
)

// Event is the afterFileEdit hook event.
//
// Cursor documents no output fields for this event; handlers observe the edit
// and rely on side effects (for example formatting the file on disk).
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

// RegisterHandler registers an observe-only AfterFileEdit handler on d.
// Cursor documents no afterFileEdit output fields; handlers return only an error.
func RegisterHandler(d *hookkit.Dialect, fn func(context.Context, Event) error) {
	if fn == nil {
		return
	}
	register(d.Codec())
	hookkit.RegisterObserve(d, fn)
}
