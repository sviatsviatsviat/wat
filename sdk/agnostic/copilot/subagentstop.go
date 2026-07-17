package copilot

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
)

// MapSubagentStop maps a Copilot SubagentStop hook into a unified Event.
func MapSubagentStop(e sdkcopilot.SubagentStop, raw []byte) *model.Event {
	ev := newEvent(e, raw, model.KindSubagentStop)
	ev.Subagent = &model.Subagent{
		Type:   e.Name(),
		Status: e.Reason(),
	}
	ev.Turn = &model.TurnEnd{Status: e.Reason()}
	return ev
}
