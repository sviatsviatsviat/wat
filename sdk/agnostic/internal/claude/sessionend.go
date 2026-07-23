package claude

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/run"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
)

// RegisterSessionEnd registers fn on the Claude SessionEnd chain.
func RegisterSessionEnd(fn model.SessionEndHandler) run.Hooks {
	if fn == nil {
		return nil
	}
	return sdkclaude.UseHooks().SessionEnd(func(ctx context.Context, hook sdkclaude.SessionEnd) error {
		return fn(ctx, *mapSessionEnd(hook))
	})
}

func mapSessionEnd(e sdkclaude.SessionEnd) *model.SessionEndEvent {
	return &model.SessionEndEvent{
		Envelope: envelope(e.Envelope, e.EventName()),
		Life:     &model.Lifecycle{Reason: e.Reason},
	}
}
