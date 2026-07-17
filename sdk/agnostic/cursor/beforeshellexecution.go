package cursor

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

// MapBeforeShellExecution maps a Cursor BeforeShellExecution hook into a unified Event.
func MapBeforeShellExecution(e sdkcursor.BeforeShellExecution, raw []byte) *model.Event {
	ev := newEvent(e, raw, model.KindPreTool)
	ev.Tool = &model.ToolCall{Name: model.ToolBash, Native: receivedName(e), Shell: e.Command}
	return ev
}
