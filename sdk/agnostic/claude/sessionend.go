package claude

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
)

// MapSessionEnd maps a Claude SessionEnd hook into a unified Event.
func MapSessionEnd(e sdkclaude.SessionEnd, raw []byte) *model.Event {
	ev := newEvent(e, raw, model.KindSessionEnd)
	ev.Life = &model.Lifecycle{Reason: e.Reason}
	return ev
}
