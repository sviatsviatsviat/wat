package claude

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/adapter"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
)

// MapPreToolUse maps a Claude PreToolUse hook into a unified Event.
func MapPreToolUse(e sdkclaude.PreToolUse, raw []byte) *model.Event {
	ev := newEvent(e, raw, model.KindPreTool)
	ev.Tool = adapter.NewToolCall(e.ToolName, e.ToolInput.Raw(), e.ToolUseID)
	return ev
}
