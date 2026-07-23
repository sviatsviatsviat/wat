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
	sdkcursor.UseHooks().SubagentStart(func(ctx context.Context, hook sdkcursor.SubagentStart, _ sdkcursor.SubagentStartResults) (sdkcursor.PermissionOutput, error) {
		return nil, fn(ctx, *mapSubagentStart(hook))
	})
}

func mapSubagentStart(e sdkcursor.SubagentStart) *model.SubagentStartEvent {
	return &model.SubagentStartEvent{
		Envelope: envelope(e.Envelope, e.EventName()),
		Subagent: &model.Subagent{ID: e.SubagentID, Type: e.SubagentType, Task: e.Task},
	}
}
