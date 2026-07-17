package cursor

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/adapter"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

// MapPostToolUseFailure maps a Cursor PostToolUseFailure hook into a unified Event.
func MapPostToolUseFailure(e sdkcursor.PostToolUseFailure, raw []byte) *model.Event {
	ev := newEvent(e, raw, model.KindPostToolFailure)
	ev.Tool = adapter.NewToolCall(e.ToolName, e.ToolInput.Raw(), e.ToolUseID)
	ev.Result = &model.ToolResult{
		Error:       e.ErrorMessage,
		FailureType: e.FailureType,
		DurationMs:  e.DurationMillis(),
	}
	return ev
}
