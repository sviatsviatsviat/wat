package copilot

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
)

// MapSubagentStart maps a Copilot SubagentStart hook into a unified Event.
func MapSubagentStart(e sdkcopilot.SubagentStart, raw []byte) *model.Event {
	ev := newEvent(e, raw, model.KindSubagentStart)
	ev.Subagent = &model.Subagent{
		Type:    e.Name(),
		Task:    e.DisplayName(),
		Summary: e.AgentDescription,
	}
	return ev
}
