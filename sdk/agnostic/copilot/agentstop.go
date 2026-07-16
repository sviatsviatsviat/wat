package copilot

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
)

func mapAgentStop(e sdkcopilot.AgentStop, ev *model.Event) {
	ev.Turn = &model.TurnEnd{Status: e.Reason()}
}

func mapStopOutput(res model.Result) any {
	if res.FollowUp == "" {
		return nil
	}
	return sdkcopilot.BuildStopOutput(res.FollowUp)
}
