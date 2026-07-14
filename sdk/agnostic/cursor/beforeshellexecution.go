package cursor

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

func mapBeforeShellExecution(e sdkcursor.BeforeShellExecution, ev *model.Event, name string) {
	ev.Tool = &model.ToolCall{Name: model.ToolBash, Native: name, Shell: e.Command}
}
