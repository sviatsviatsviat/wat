package cursor

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/run"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

// RegisterPreCompact registers fn on the Cursor PreCompact chain.
func RegisterPreCompact(r *run.Registry, fn model.PreCompactHandler) {
	if fn == nil {
		return
	}
	sdkcursor.UseHooks(r).PreCompact(func(ctx context.Context, hook run.Hook[sdkcursor.PreCompact], _ sdkcursor.PreCompactResults) (sdkcursor.PreCompactOutput, error) {
		return nil, fn(ctx, model.NewPreCompactHook(hook.Invocation(), mapPreCompact(hook.Event)))
	})
}

func mapPreCompact(e sdkcursor.PreCompact) *model.PreCompactEvent {
	return &model.PreCompactEvent{
		Envelope: envelope(e.Envelope, e.EventName()),
		Compact:  &model.CompactInfo{Trigger: e.Trigger},
	}
}
