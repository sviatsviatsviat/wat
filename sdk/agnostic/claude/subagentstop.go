package claude

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
)

func mapSubagentStop(e sdkclaude.SubagentStop, ev *model.Event) {
	ev.Subagent = &model.Subagent{ID: e.AgentID, Type: e.AgentType, Summary: e.LastAssistantMessage}
	ev.Turn = &model.TurnEnd{StopHookActive: e.StopHookActive, LastAssistantMessage: e.LastAssistantMessage}
}
