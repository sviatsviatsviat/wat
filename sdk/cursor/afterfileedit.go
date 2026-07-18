package cursor

import (
	"context"

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
	registerDecoder(EventAfterFileEdit, decodeAs[AfterFileEdit])
}

// OnAfterFileEdit registers an afterFileEdit handler.
func OnAfterFileEdit(fn func(context.Context, Hook[AfterFileEdit], PostToolResults) (PostToolOutput, error)) *chain {
	return (&chain{}).AfterFileEdit(fn)
}

// AfterFileEdit registers another AfterFileEdit handler on the chain.
func (c *chain) AfterFileEdit(fn func(context.Context, Hook[AfterFileEdit], PostToolResults) (PostToolOutput, error)) *chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev AfterFileEdit) (PostToolOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), postToolResults{})
	})
	return c
}
