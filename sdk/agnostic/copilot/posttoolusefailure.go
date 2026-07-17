package copilot

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/adapter"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
)

// MapPostToolUseFailure maps a Copilot PostToolUseFailure hook into a unified Event.
func MapPostToolUseFailure(e sdkcopilot.PostToolUseFailure, raw []byte) *model.Event {
	ev := newEvent(e, raw, model.KindPostToolFailure)
	ev.Tool = adapter.NewToolCall(e.NativeToolName(), e.Input().Raw(), "")
	if msg := e.ErrorMessage(); msg != "" {
		ev.Result = &model.ToolResult{Error: msg}
	}
	return ev
}
