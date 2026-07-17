package cursor

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

// MapStop maps a Cursor Stop hook into a unified Event.
func MapStop(e sdkcursor.Stop, raw []byte) *model.Event {
	ev := newEvent(e, raw, model.KindStop)
	ev.Turn = &model.TurnEnd{Status: e.Status, LoopCount: e.LoopCount}
	return ev
}
