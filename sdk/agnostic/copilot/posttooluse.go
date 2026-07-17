package copilot

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/adapter"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
)

// MapPostToolUse maps a Copilot PostToolUse hook into a unified Event.
func MapPostToolUse(e sdkcopilot.PostToolUse, raw []byte) *model.Event {
	ev := newEvent(e, raw, model.KindPostTool)
	ev.Tool = adapter.NewToolCall(e.NativeToolName(), e.Input().Raw(), "")
	ev.Result = &model.ToolResult{Raw: adapter.CloneRaw(e.ResultRaw()), Text: e.ResultText()}
	return ev
}
