package claude

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/adapter"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
)

func mapPostToolUseFailure(e sdkclaude.PostToolUseFailure, ev *model.Event) {
	ev.Tool = adapter.NewToolCall(e.ToolName, e.ToolInput.Raw(), e.ToolUseID)
	ev.Result = &model.ToolResult{Error: e.Error}
}

func mapPostToolFailureOutput(res model.Result) any {
	if res.Context == "" {
		return nil
	}
	return sdkclaude.BuildPostToolUseFailureOutput(res.Context)
}
