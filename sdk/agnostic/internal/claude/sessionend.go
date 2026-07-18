package claude

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
)

// RegisterSessionEnd registers fn on the Claude SessionEnd chain.
func RegisterSessionEnd(fn model.SessionEndHandler) {
	if fn == nil {
		return
	}
	sdkclaude.OnSessionEnd(func(ctx context.Context, hook sdkclaude.Hook[sdkclaude.SessionEnd]) error {
		return fn(ctx, model.NewSessionEndHook(hook.Invocation(), mapSessionEnd(hook.Event)))
	})
}

func mapSessionEnd(e sdkclaude.SessionEnd) *model.SessionEndEvent {
	return &model.SessionEndEvent{
		Envelope: envelope(e),
		Life:     &model.Lifecycle{Reason: e.Reason},
	}
}
