package cursor

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

// MapPreCompact maps a Cursor PreCompact hook into a unified Event.
func MapPreCompact(e sdkcursor.PreCompact, raw []byte) *model.Event {
	ev := newEvent(e, raw, model.KindPreCompact)
	ev.Compact = &model.CompactInfo{Trigger: e.Trigger}
	return ev
}
