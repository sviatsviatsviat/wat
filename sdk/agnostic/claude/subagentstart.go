package claude

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
)

// MapSubagentStart maps a Claude SubagentStart hook into a unified Event.
func MapSubagentStart(e sdkclaude.SubagentStart, raw []byte) *model.Event {
	ev := newEvent(e, raw, model.KindSubagentStart)
	ev.Subagent = &model.Subagent{ID: e.AgentID, Type: e.AgentType}
	return ev
}
