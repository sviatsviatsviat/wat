package cursor

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

// MapAfterShellExecution maps a Cursor AfterShellExecution hook into a unified Event.
func MapAfterShellExecution(e sdkcursor.AfterShellExecution, raw []byte) *model.Event {
	ev := newEvent(e, raw, model.KindPostTool)
	ev.Tool = &model.ToolCall{Name: model.ToolBash, Native: receivedName(e), Shell: e.Command}
	ev.Result = &model.ToolResult{Text: e.Output, DurationMs: e.DurationMillis()}
	return ev
}
