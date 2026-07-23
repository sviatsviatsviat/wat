package copilot

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/run"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
)

// RegisterPreCompact registers fn on the Copilot PreCompact chain.
func RegisterPreCompact(fn model.PreCompactHandler) {
	if fn == nil {
		return
	}
	sdkcopilot.UseHooks().PreCompact(func(ctx context.Context, hook run.Hook[sdkcopilot.PreCompact]) error {
		return fn(ctx, model.NewPreCompactHook(hook.Invocation(), mapPreCompact(hook.Event)))
	})
}

func mapPreCompact(e sdkcopilot.PreCompact) *model.PreCompactEvent {
	return &model.PreCompactEvent{
		Envelope: envelope(e.Envelope, e.EventName()),
		Compact: &model.CompactInfo{
			Trigger:            e.Trigger,
			CustomInstructions: e.Instructions(),
		},
	}
}
