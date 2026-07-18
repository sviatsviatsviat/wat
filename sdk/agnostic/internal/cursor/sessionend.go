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
	sdkcursor.OnSessionEnd(func(ctx context.Context, hook sdkcursor.Hook[sdkcursor.SessionEnd]) error {
		return fn(ctx, model.NewSessionEndHook(hook.Invocation(), mapSessionEnd(hook.Event)))
	})
}

func mapSessionEnd(e sdkcursor.SessionEnd) *model.SessionEndEvent {
	return &model.SessionEndEvent{
		Envelope: envelope(e),
		Life:     &model.Lifecycle{Reason: e.Reason, Background: e.IsBackgroundAgent},
	}
}
