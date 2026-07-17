package claude

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/adapter"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
)

func mapPreToolUse(e sdkclaude.PreToolUse, ev *model.Event) {
	ev.Tool = adapter.NewToolCall(e.ToolName, e.ToolInput.Raw(), e.ToolUseID)
}
