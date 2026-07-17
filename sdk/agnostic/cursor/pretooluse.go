package cursor

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/adapter"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

// MapPreToolUse maps a Cursor PreToolUse hook into a unified Event.
func MapPreToolUse(e sdkcursor.PreToolUse, raw []byte) *model.Event {
	ev := newEvent(e, raw, model.KindPreTool)
	ev.Tool = adapter.NewToolCall(e.ToolName, e.ToolInput.Raw(), e.ToolUseID)
	if shell := e.ShellCommand(); shell != "" {
		ev.Tool.Shell = shell
	}
	return ev
}
