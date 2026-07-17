package copilot

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
)

func mapAgentStop(e sdkcopilot.AgentStop, ev *model.Event) {
	ev.Turn = &model.TurnEnd{Status: e.Reason()}
}
