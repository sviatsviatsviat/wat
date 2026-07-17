package claude

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
)

// MapPreCompact maps a Claude PreCompact hook into a unified Event.
func MapPreCompact(e sdkclaude.PreCompact, raw []byte) *model.Event {
	ev := newEvent(e, raw, model.KindPreCompact)
	ev.Compact = &model.CompactInfo{Trigger: e.Trigger, CustomInstructions: e.CustomInstructions}
	return ev
}
