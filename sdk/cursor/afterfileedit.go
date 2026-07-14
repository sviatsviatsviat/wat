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

// AfterFileEdit registers an afterFileEdit handler.
func (c *Chain) AfterFileEdit(fn func(context.Context, AfterFileEditHook, PostToolResults) (PostToolOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev AfterFileEdit) (PostToolOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), postToolResults{})
	})
	return &Chain{}
}
