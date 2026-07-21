package cursor

import (
	"context"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// AfterFileEdit is the afterFileEdit hook event.
type AfterFileEdit struct {
	Envelope
	// FilePath is the edited file path.
	FilePath string `json:"file_path"`
	// Edits are the applied edits.
	Edits []Edit `json:"edits"`
}

// EventName returns the canonical hook event name.
func (AfterFileEdit) EventName() string { return EventAfterFileEdit }

func init() {
	codec.Register(EventAfterFileEdit, hookkit.EventDecoder[AfterFileEdit](codec))
}

// AfterFileEdit registers a AfterFileEdit handler on the chain.
func (c *chain) AfterFileEdit(fn func(context.Context, run.Hook[AfterFileEdit], PostToolResults) (PostToolOutput, error)) *chain {
	if fn == nil {
		return c
	}
	c.reg.RegisterHandler(Dialect, run.Handler(func(ctx context.Context, hook run.Hook[AfterFileEdit]) (PostToolOutput, error) {
		return fn(ctx, hook, postToolResults{})
	}))
	return c
}
