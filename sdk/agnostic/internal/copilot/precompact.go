package copilot

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
)

// RegisterPreCompact registers fn on the Copilot PreCompact chain.
func RegisterPreCompact(fn model.PreCompactHandler) {
	if fn == nil {
		return
	}
	sdkcopilot.OnPreCompact(func(ctx context.Context, hook sdkcopilot.Hook[sdkcopilot.PreCompact]) error {
		return fn(ctx, model.NewPreCompactHook(hook.Invocation(), mapPreCompact(hook.Event, hook.Raw())))
	})
}

func mapPreCompact(e sdkcopilot.PreCompact, raw []byte) *model.PreCompactEvent {
	return &model.PreCompactEvent{
		Envelope: envelope(e, raw),
		Compact: &model.CompactInfo{
			Trigger:            e.Trigger,
			CustomInstructions: e.Instructions(),
		},
	}
}
