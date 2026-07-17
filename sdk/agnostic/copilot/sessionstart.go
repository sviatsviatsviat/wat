package copilot

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
)

// MapSessionStart maps a Copilot SessionStart hook into a unified Event.
func MapSessionStart(e sdkcopilot.SessionStart, raw []byte) *model.Event {
	ev := newEvent(e, raw, model.KindSessionStart)
	ev.Life = &model.Lifecycle{Source: e.Source, InitialPrompt: e.InitialPrompt()}
	return ev
}
