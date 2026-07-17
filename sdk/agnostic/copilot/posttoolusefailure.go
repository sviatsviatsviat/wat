package copilot

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/adapter"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
)

func mapPostToolUseFailure(e sdkcopilot.PostToolUseFailure, ev *model.Event) {
	ev.Tool = adapter.NewToolCall(e.NativeToolName(), e.Input().Raw(), "")
	if msg := e.ErrorMessage(); msg != "" {
		ev.Result = &model.ToolResult{Error: msg}
	}
}
