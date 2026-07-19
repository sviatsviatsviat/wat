package cursor

import (
	"context"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// AfterTabFileEdit is the afterTabFileEdit hook event.
type AfterTabFileEdit struct {
	Envelope
	// FilePath is the edited file path.
	FilePath string `json:"file_path"`
	// Edits are the applied edits.
	Edits []Edit `json:"edits"`
}

// EventName returns the canonical hook event name.
func (AfterTabFileEdit) EventName() string { return EventAfterTabFileEdit }

func init() {
	codec.Register(EventAfterTabFileEdit, hookkit.EventDecoder[AfterTabFileEdit](codec))
}

// AfterTabFileEdit registers a AfterTabFileEdit handler on the chain.
func (c *chain) AfterTabFileEdit(fn func(context.Context, run.Hook[AfterTabFileEdit]) error) *chain {
	registerObserveHandler(c.reg, fn)
	return c
}
