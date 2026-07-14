package cursor

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

func mapSubagentStart(e sdkcursor.SubagentStart, ev *model.Event) {
	ev.Subagent = &model.Subagent{ID: e.SubagentID, Type: e.SubagentType, Task: e.Task}
}
