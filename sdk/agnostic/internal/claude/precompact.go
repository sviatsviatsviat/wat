package claude

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
)

// RegisterPreCompact registers fn on the Claude PreCompact chain.
func RegisterPreCompact(fn model.PreCompactHandler) {
	if fn == nil {
		return
	}
	sdkclaude.Adapter().PreCompact(func(ctx context.Context, hook sdkclaude.Hook[sdkclaude.PreCompact], _ sdkclaude.PreCompactResults) (sdkclaude.CommonOutput, error) {
		return nil, fn(ctx, model.NewPreCompactHook(hook.Invocation(), mapPreCompact(hook.Event, hook.Raw())))
	})
}

func mapPreCompact(e sdkclaude.PreCompact, raw []byte) *model.PreCompactEvent {
	return &model.PreCompactEvent{
		Envelope: envelope(e, raw),
		Compact:  &model.CompactInfo{Trigger: e.Trigger, CustomInstructions: e.CustomInstructions},
	}
}
