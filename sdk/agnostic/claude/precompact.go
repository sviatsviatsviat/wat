package claude

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
)

func mapPreCompact(e sdkclaude.PreCompact, ev *model.Event) {
	ev.Compact = &model.CompactInfo{Trigger: e.Trigger, CustomInstructions: e.CustomInstructions}
}
