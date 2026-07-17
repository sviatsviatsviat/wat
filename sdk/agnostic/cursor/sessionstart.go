package cursor

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

// MapSessionStart maps a Cursor SessionStart hook into a unified Event.
func MapSessionStart(e sdkcursor.SessionStart, raw []byte) *model.Event {
	ev := newEvent(e, raw, model.KindSessionStart)
	ev.Life = &model.Lifecycle{Model: e.Model, Background: e.IsBackgroundAgent}
	return ev
}
