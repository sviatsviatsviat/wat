package copilot

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
)

func mapSubagentStart(e sdkcopilot.SubagentStart, ev *model.Event) {
	ev.Subagent = &model.Subagent{
		Type:    e.Name(),
		Task:    e.DisplayName(),
		Summary: e.AgentDescription,
	}
}
