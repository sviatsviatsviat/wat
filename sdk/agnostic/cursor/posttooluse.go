package cursor

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/adapter"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

// MapPostToolUse maps a Cursor PostToolUse hook into a unified Event.
func MapPostToolUse(e sdkcursor.PostToolUse, raw []byte) *model.Event {
	ev := newEvent(e, raw, model.KindPostTool)
	ev.Tool = adapter.NewToolCall(e.ToolName, e.ToolInput.Raw(), e.ToolUseID)
	ev.Result = &model.ToolResult{Text: e.ToolOutput, DurationMs: e.DurationMillis()}
	return ev
}
