package claude

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/adapter"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
)

func mapPreToolUse(e sdkclaude.PreToolUse, ev *model.Event) {
	ev.Tool = adapter.NewToolCall(e.ToolName, e.ToolInput.Raw(), e.ToolUseID)
}

func mapPreToolOutput(res model.Result) any {
	out := sdkclaude.PreToolUseOutput{}
	if d := res.Decision.String(); d != "" {
		out.Decision = sdkclaude.PermissionDecision(d)
		out.Reason = res.Reason
	}
	out.UpdatedInput = res.UpdatedInput
	return out
}
