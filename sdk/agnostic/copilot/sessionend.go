package copilot

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
)

func mapSessionEnd(e sdkcopilot.SessionEnd, ev *model.Event) {
	ev.Life = &model.Lifecycle{Reason: e.Reason}
}
