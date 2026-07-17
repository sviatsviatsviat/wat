package claude

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/adapter"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
)

// MapPostToolUseFailure maps a Claude PostToolUseFailure hook into a unified Event.
func MapPostToolUseFailure(e sdkclaude.PostToolUseFailure, raw []byte) *model.Event {
	ev := newEvent(e, raw, model.KindPostToolFailure)
	ev.Tool = adapter.NewToolCall(e.ToolName, e.ToolInput.Raw(), e.ToolUseID)
	ev.Result = &model.ToolResult{Error: e.Error}
	return ev
}
