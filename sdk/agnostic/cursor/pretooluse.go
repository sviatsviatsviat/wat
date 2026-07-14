package cursor

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/adapter"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

func mapPreToolUse(e sdkcursor.PreToolUse, ev *model.Event) {
	ev.Tool = adapter.NewToolCall(e.ToolName, e.ToolInput, e.ToolUseID)
	if shell := e.ShellCommand(); shell != "" {
		ev.Tool.Shell = shell
	}
}

func mapPreToolOutput(ev *model.Event, res model.Result) any {
	out := sdkcursor.PermissionOutput{}
	if d := res.Decision.String(); d != "" {
		out.Decision = sdkcursor.PermissionDecision(d)
	}
	out.AgentMessage = res.Reason
	if res.UpdatedInput != nil && ev.Name == sdkcursor.EventPreToolUse {
		out.UpdatedInput = res.UpdatedInput
	}
	return out
}
