package copilot

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/run"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
)

// RegisterSessionEnd registers fn on the Copilot SessionEnd chain.
func RegisterSessionEnd(fn model.SessionEndHandler) run.Hooks {
	if fn == nil {
		return nil
	}
	return sdkcopilot.UseHooks().SessionEnd(func(ctx context.Context, hook sdkcopilot.SessionEnd) error {
		return fn(ctx, *mapSessionEnd(hook))
	})
}

func mapSessionEnd(e sdkcopilot.SessionEnd) *model.SessionEndEvent {
	return &model.SessionEndEvent{
		Envelope: envelope(e.Envelope, e.EventName()),
		Life:     &model.Lifecycle{Reason: e.Reason},
	}
}
