package copilot

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
)

func mapSubagentStop(e sdkcopilot.SubagentStop, ev *model.Event) {
	ev.Subagent = &model.Subagent{
		Type:   e.Name(),
		Status: e.Reason(),
	}
	ev.Turn = &model.TurnEnd{Status: e.Reason()}
}
