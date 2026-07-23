package copilot

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/run"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
)

// RegisterPreCompact registers fn on the Copilot PreCompact chain.
func RegisterPreCompact(fn model.PreCompactHandler) run.Hooks {
	if fn == nil {
		return nil
	}
	return sdkcopilot.UseHooks().PreCompact(func(ctx context.Context, hook sdkcopilot.PreCompact) error {
		return fn(ctx, *mapPreCompact(hook))
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
