package copilot

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
)

// RegisterSessionEnd registers fn on the Copilot SessionEnd chain.
func RegisterSessionEnd(fn model.SessionEndHandler) {
	if fn == nil {
		return
	}
	sdkcopilot.Adapter().SessionEnd(func(ctx context.Context, hook sdkcopilot.Hook[sdkcopilot.SessionEnd]) error {
		return fn(ctx, model.NewSessionEndHook(hook.Invocation(), mapSessionEnd(hook.Event, hook.Raw())))
	})
}

func mapSessionEnd(e sdkcopilot.SessionEnd, raw []byte) *model.SessionEndEvent {
	return &model.SessionEndEvent{
		Envelope: envelope(e, raw),
		Life:     &model.Lifecycle{Reason: e.Reason},
	}
}
