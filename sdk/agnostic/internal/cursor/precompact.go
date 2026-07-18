package cursor

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

// RegisterPreCompact registers fn on the Cursor PreCompact chain.
func RegisterPreCompact(fn model.PreCompactHandler) {
	if fn == nil {
		return
	}
	sdkcursor.OnPreCompact(func(ctx context.Context, hook sdkcursor.Hook[sdkcursor.PreCompact], _ sdkcursor.PreCompactResults) (sdkcursor.PreCompactOutput, error) {
		return nil, fn(ctx, model.NewPreCompactHook(hook.Invocation(), mapPreCompact(hook.Event, hook.Raw())))
	})
}

func mapPreCompact(e sdkcursor.PreCompact, raw []byte) *model.PreCompactEvent {
	return &model.PreCompactEvent{
		Envelope: envelope(e, raw),
		Compact:  &model.CompactInfo{Trigger: e.Trigger},
	}
}
