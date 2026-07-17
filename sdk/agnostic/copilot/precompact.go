package copilot

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
)

// MapPreCompact maps a Copilot PreCompact hook into a unified Event.
func MapPreCompact(e sdkcopilot.PreCompact, raw []byte) *model.Event {
	ev := newEvent(e, raw, model.KindPreCompact)
	ev.Compact = &model.CompactInfo{
		Trigger:            e.Trigger,
		CustomInstructions: e.Instructions(),
	}
	return ev
}
