package cursor

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

func mapAfterShellExecution(e sdkcursor.AfterShellExecution, ev *model.Event, name string) {
	ev.Tool = &model.ToolCall{Name: model.ToolBash, Native: name, Shell: e.Command}
	ev.Result = &model.ToolResult{Text: e.Output, DurationMs: e.DurationMillis()}
}
