package cursor

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

// RegisterSubagentStart registers fn on the Cursor SubagentStart chain.
func RegisterSubagentStart(fn model.SubagentStartHandler) {
	if fn == nil {
		return
	}
	sdkcursor.Adapter().SubagentStart(func(ctx context.Context, hook sdkcursor.Hook[sdkcursor.SubagentStart], _ sdkcursor.SubagentStartResults) (sdkcursor.PermissionOutput, error) {
		return nil, fn(ctx, model.NewSubagentStartHook(hook.Invocation(), mapSubagentStart(hook.Event, hook.Raw())))
	})
}

func mapSubagentStart(e sdkcursor.SubagentStart, raw []byte) *model.SubagentStartEvent {
	return &model.SubagentStartEvent{
		Envelope: envelope(e, raw),
		Subagent: &model.Subagent{ID: e.SubagentID, Type: e.SubagentType, Task: e.Task},
	}
}
