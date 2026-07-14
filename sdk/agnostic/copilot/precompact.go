package copilot

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
)

func mapPreCompact(e sdkcopilot.PreCompact, ev *model.Event) {
	ev.Compact = &model.CompactInfo{
		Trigger:            e.Trigger,
		CustomInstructions: e.Instructions(),
	}
}
