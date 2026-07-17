package copilot

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
)

// MapAgentStop maps a Copilot AgentStop hook into a unified Event.
func MapAgentStop(e sdkcopilot.AgentStop, raw []byte) *model.Event {
	ev := newEvent(e, raw, model.KindStop)
	ev.Turn = &model.TurnEnd{Status: e.Reason()}
	return ev
}
