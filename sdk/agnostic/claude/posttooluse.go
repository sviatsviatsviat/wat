package claude

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/adapter"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
)

// MapPostToolUse maps a Claude PostToolUse hook into a unified Event.
func MapPostToolUse(e sdkclaude.PostToolUse, raw []byte) *model.Event {
	ev := newEvent(e, raw, model.KindPostTool)
	ev.Tool = adapter.NewToolCall(e.ToolName, e.ToolInput.Raw(), e.ToolUseID)
	ev.Result = &model.ToolResult{Raw: adapter.CloneRaw(e.ToolResponse), Text: adapter.RawToText(e.ToolResponse)}
	return ev
}
