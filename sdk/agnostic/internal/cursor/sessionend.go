package cursor

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

// RegisterSessionEnd registers fn on the Cursor SessionEnd chain.
func RegisterSessionEnd(fn model.SessionEndHandler) {
	if fn == nil {
		return
	}
	sdkcursor.UseHooks().SessionEnd(func(ctx context.Context, hook sdkcursor.SessionEnd) error {
		return fn(ctx, *mapSessionEnd(hook))
	})
}

func mapSessionEnd(e sdkcursor.SessionEnd) *model.SessionEndEvent {
	return &model.SessionEndEvent{
		Envelope: envelope(e.Envelope, e.EventName()),
		Life:     &model.Lifecycle{Reason: e.Reason, Background: e.IsBackgroundAgent},
	}
}
