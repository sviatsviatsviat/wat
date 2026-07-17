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
	new(sdkclaude.Chain).SessionEnd(func(ctx context.Context, hook sdkclaude.Hook[sdkclaude.SessionEnd]) error {
		return fn(ctx, model.NewSessionEndHook(hook.Invocation(), mapSessionEnd(hook.Event, hook.Raw())))
	})
}

func mapSessionEnd(e sdkclaude.SessionEnd, raw []byte) *model.SessionEndEvent {
	return &model.SessionEndEvent{
		Envelope: envelope(e, raw),
		Life:     &model.Lifecycle{Reason: e.Reason},
	}
}
