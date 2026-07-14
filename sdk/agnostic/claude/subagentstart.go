package claude

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
)

func mapSubagentStart(e sdkclaude.SubagentStart, ev *model.Event) {
	ev.Subagent = &model.Subagent{ID: e.AgentID, Type: e.AgentType}
}
