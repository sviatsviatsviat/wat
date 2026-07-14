package claude

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
)

func mapSessionEnd(e sdkclaude.SessionEnd, ev *model.Event) {
	ev.Life = &model.Lifecycle{Reason: e.Reason}
}
