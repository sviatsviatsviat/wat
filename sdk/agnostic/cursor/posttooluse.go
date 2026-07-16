package cursor

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/adapter"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

func mapPostToolUse(e sdkcursor.PostToolUse, ev *model.Event) {
	ev.Tool = adapter.NewToolCall(e.ToolName, e.ToolInput.Raw(), e.ToolUseID)
	ev.Result = &model.ToolResult{Text: e.ToolOutput, DurationMs: e.DurationMillis()}
}

func mapPostToolOutput(res model.Result) any {
	var updated any
	if res.UpdatedOutput != nil {
		updated = *res.UpdatedOutput
	}
	if res.Context == "" && updated == nil {
		return nil
	}
	return sdkcursor.BuildPostToolOutput(res.Context, updated)
}
