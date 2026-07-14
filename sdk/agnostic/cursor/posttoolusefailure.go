package cursor

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/adapter"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

func mapPostToolUseFailure(e sdkcursor.PostToolUseFailure, ev *model.Event) {
	ev.Tool = adapter.NewToolCall(e.ToolName, e.ToolInput, e.ToolUseID)
	ev.Result = &model.ToolResult{
		Error:       e.ErrorMessage,
		FailureType: e.FailureType,
		DurationMs:  e.DurationMillis(),
	}
}

func mapPostToolFailureOutput(res model.Result) any {
	return sdkcursor.PostToolOutput{AdditionalContext: res.Context}
}
