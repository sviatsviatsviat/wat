package aftertabfileedit

import (
	"context"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// Event is the afterTabFileEdit hook event.
type Event struct {
	event.Envelope
	// FilePath is the edited file path.
	FilePath string `json:"file_path"`
	// Edits are the applied edits.
	Edits []event.Edit `json:"edits"`
}

// EventName returns the canonical hook event name.
func (Event) EventName() string { return event.AfterTabFileEdit }

// Register registers this hook event decoder on c.
func Register(c *hookkit.Codec) {
	c.Register(event.AfterTabFileEdit, hookkit.EventDecoder[Event](c))
}

// RegisterHandler registers an observe handler on reg.
func RegisterHandler(d *hookkit.Dialect, fn func(context.Context, run.Hook[Event]) error) {
	if fn == nil {
		return
	}
	d.Register(hookkit.ObserveHandler(fn))
}
