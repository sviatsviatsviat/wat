package claude

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
)

// RegisterSubagentStart registers fn on the Claude SubagentStart chain.
func RegisterSubagentStart(fn model.SubagentStartHandler) {
	if fn == nil {
		return
	}
	sdkclaude.OnSubagentStart(func(ctx context.Context, hook sdkclaude.Hook[sdkclaude.SubagentStart], _ sdkclaude.SubagentStartResults) (sdkclaude.CommonOutput, error) {
		return nil, fn(ctx, model.NewSubagentStartHook(hook.Invocation(), mapSubagentStart(hook.Event)))
	})
}

func mapSubagentStart(e sdkclaude.SubagentStart) *model.SubagentStartEvent {
	return &model.SubagentStartEvent{
		Envelope: envelope(e),
		Subagent: &model.Subagent{ID: e.AgentID, Type: e.AgentType},
	}
}
