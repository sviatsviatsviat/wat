package copilot

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
)

func mapSubagentStop(e sdkcopilot.SubagentStop, ev *model.Event) {
	ev.Subagent = &model.Subagent{
		Type: e.Name(),
		Task: e.DisplayName(),
	}
	ev.Turn = &model.TurnEnd{Status: e.Reason()}
}
