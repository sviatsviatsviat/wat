package copilot

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/run"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
)

// RegisterSubagentStart registers fn on the Copilot SubagentStart chain.
func RegisterSubagentStart(fn model.SubagentStartHandler) {
	if fn == nil {
		return
	}
	sdkcopilot.UseHooks().SubagentStart(func(ctx context.Context, hook run.Hook[sdkcopilot.SubagentStart], _ sdkcopilot.SubagentStartResults) (sdkcopilot.SubagentStartOutput, error) {
		return nil, fn(ctx, model.NewSubagentStartHook(hook.Invocation(), mapSubagentStart(hook.Event)))
	})
}

func mapSubagentStart(e sdkcopilot.SubagentStart) *model.SubagentStartEvent {
	return &model.SubagentStartEvent{
		Envelope: envelope(e.Envelope, e.EventName()),
		Subagent: &model.Subagent{
			Type:    e.Name(),
			Task:    e.DisplayName(),
			Summary: e.AgentDescription,
		},
	}
}
