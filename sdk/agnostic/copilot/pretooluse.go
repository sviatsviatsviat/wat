package copilot

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/adapter"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
)

func mapPreToolUse(e sdkcopilot.PreToolUse, ev *model.Event) {
	ev.Tool = adapter.NewToolCall(e.NativeToolName(), e.Input().Raw(), "")
}

func mapPreToolOutput(res model.Result) any {
	out := sdkcopilot.PreToolOutput{}
	if d := res.Decision.String(); d != "" {
		out.Decision = sdkcopilot.PermissionDecision(d)
		out.Reason = res.Reason
	}
	out.ModifiedArgs = res.UpdatedInput
	return out
}
