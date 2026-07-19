package cursor

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/run"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

// RegisterSubagentStart registers fn on the Cursor SubagentStart chain.
func RegisterSubagentStart(r *run.Registry, fn model.SubagentStartHandler) {
	if fn == nil {
		return
	}
	sdkcursor.UseHooks(r).SubagentStart(func(ctx context.Context, hook run.Hook[sdkcursor.SubagentStart], _ sdkcursor.SubagentStartResults) (sdkcursor.PermissionOutput, error) {
		return nil, fn(ctx, model.NewSubagentStartHook(hook.Invocation(), mapSubagentStart(hook.Event)))
	})
}

func mapSubagentStart(e sdkcursor.SubagentStart) *model.SubagentStartEvent {
	return &model.SubagentStartEvent{
		Envelope: envelope(e.Envelope, e.EventName()),
		Subagent: &model.Subagent{ID: e.SubagentID, Type: e.SubagentType, Task: e.Task},
	}
}
