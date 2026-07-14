package cursor

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/adapter"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

func mapPostToolUse(e sdkcursor.PostToolUse, ev *model.Event) {
	ev.Tool = adapter.NewToolCall(e.ToolName, e.ToolInput, e.ToolUseID)
	ev.Result = &model.ToolResult{Text: e.ToolOutput, DurationMs: e.DurationMillis()}
}

func mapPostToolOutput(res model.Result) any {
	out := sdkcursor.PostToolOutput{AdditionalContext: res.Context}
	if res.UpdatedOutput != nil {
		out.UpdatedMCPOutput = *res.UpdatedOutput
	}
	return out
}
