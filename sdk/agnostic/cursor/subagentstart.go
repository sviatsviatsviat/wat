package cursor

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

// MapSubagentStart maps a Cursor SubagentStart hook into a unified Event.
func MapSubagentStart(e sdkcursor.SubagentStart, raw []byte) *model.Event {
	ev := newEvent(e, raw, model.KindSubagentStart)
	ev.Subagent = &model.Subagent{ID: e.SubagentID, Type: e.SubagentType, Task: e.Task}
	return ev
}
