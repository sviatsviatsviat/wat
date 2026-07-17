package copilot

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
)

// RegisterSubagentStart registers fn on the Copilot SubagentStart chain.
func RegisterSubagentStart(fn model.SubagentStartHandler) {
	if fn == nil {
		return
	}
	new(sdkcopilot.Chain).SubagentStart(func(ctx context.Context, hook sdkcopilot.Hook[sdkcopilot.SubagentStart], _ sdkcopilot.SubagentStartResults) (sdkcopilot.SubagentStartOutput, error) {
		return nil, fn(ctx, model.NewSubagentStartHook(hook.Invocation(), mapSubagentStart(hook.Event, hook.Raw())))
	})
}

func mapSubagentStart(e sdkcopilot.SubagentStart, raw []byte) *model.SubagentStartEvent {
	return &model.SubagentStartEvent{
		Envelope: envelope(e, raw),
		Subagent: &model.Subagent{
			Type:    e.Name(),
			Task:    e.DisplayName(),
			Summary: e.AgentDescription,
		},
	}
}
