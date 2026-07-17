package copilot

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
)

func mapSessionStart(e sdkcopilot.SessionStart, ev *model.Event) {
	ev.Life = &model.Lifecycle{Source: e.Source, InitialPrompt: e.InitialPrompt()}
}
