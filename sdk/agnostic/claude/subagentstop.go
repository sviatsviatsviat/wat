package claude

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
)

// MapSubagentStop maps a Claude SubagentStop hook into a unified Event.
func MapSubagentStop(e sdkclaude.SubagentStop, raw []byte) *model.Event {
	ev := newEvent(e, raw, model.KindSubagentStop)
	ev.Subagent = &model.Subagent{ID: e.AgentID, Type: e.AgentType, Summary: e.LastAssistantMessage}
	ev.Turn = &model.TurnEnd{StopHookActive: e.StopHookActive, LastAssistantMessage: e.LastAssistantMessage}
	return ev
}
