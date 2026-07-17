package copilot

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
)

// MapSessionEnd maps a Copilot SessionEnd hook into a unified Event.
func MapSessionEnd(e sdkcopilot.SessionEnd, raw []byte) *model.Event {
	ev := newEvent(e, raw, model.KindSessionEnd)
	ev.Life = &model.Lifecycle{Reason: e.Reason}
	return ev
}
