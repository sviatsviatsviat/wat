package cursor

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

func mapPreCompact(e sdkcursor.PreCompact, ev *model.Event) {
	ev.Compact = &model.CompactInfo{Trigger: e.Trigger}
}
