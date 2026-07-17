package cursor

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

// MapSessionEnd maps a Cursor SessionEnd hook into a unified Event.
func MapSessionEnd(e sdkcursor.SessionEnd, raw []byte) *model.Event {
	ev := newEvent(e, raw, model.KindSessionEnd)
	ev.Life = &model.Lifecycle{Reason: e.Reason, Background: e.IsBackgroundAgent}
	return ev
}
