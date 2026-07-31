package claude

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/run"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
)

// RegisterPreCompact registers fn on the Claude PreCompact chain.
func RegisterPreCompact(fn model.PreCompactHandler) run.Hooks {
	if fn == nil {
		return nil
	}
	return sdkclaude.UseHooks().PreCompact(func(ctx context.Context, hook sdkclaude.PreCompact, _ sdkclaude.PreCompactResults) (sdkclaude.DecisionOutput, error) {
		return nil, fn(ctx, *mapPreCompact(hook))
	})
}

func mapPreCompact(e sdkclaude.PreCompact) *model.PreCompactEvent {
	return &model.PreCompactEvent{
		Envelope: envelope(e.Envelope, e.EventName()),
		Compact:  &model.CompactInfo{Trigger: e.Trigger, CustomInstructions: e.CustomInstructions},
	}
}
