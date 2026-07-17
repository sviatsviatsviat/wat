package copilot

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/adapter"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
)

func mapPostToolUse(e sdkcopilot.PostToolUse, ev *model.Event) {
	ev.Tool = adapter.NewToolCall(e.NativeToolName(), e.Input().Raw(), "")
	ev.Result = &model.ToolResult{Raw: adapter.CloneRaw(e.ResultRaw()), Text: e.ResultText()}
}
