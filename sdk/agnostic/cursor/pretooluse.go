package cursor

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/adapter"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

func mapPreToolUse(e sdkcursor.PreToolUse, ev *model.Event) {
	ev.Tool = adapter.NewToolCall(e.ToolName, e.ToolInput.Raw(), e.ToolUseID)
	if shell := e.ShellCommand(); shell != "" {
		ev.Tool.Shell = shell
	}
}
