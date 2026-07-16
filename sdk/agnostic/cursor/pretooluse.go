package cursor

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/adapter"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

var preToolDecisions = map[model.Decision]sdkcursor.PermissionDecision{
	model.DecisionAllow: sdkcursor.DecisionAllow,
	model.DecisionDeny:  sdkcursor.DecisionDeny,
	model.DecisionAsk:   sdkcursor.DecisionAsk,
}

func mapPreToolUse(e sdkcursor.PreToolUse, ev *model.Event) {
	ev.Tool = adapter.NewToolCall(e.ToolName, e.ToolInput.Raw(), e.ToolUseID)
	if shell := e.ShellCommand(); shell != "" {
		ev.Tool.Shell = shell
	}
}

func mapPreToolOutput(ev *model.Event, res model.Result) any {
	decision, ok := preToolDecisions[res.Decision]
	if !ok {
		if res.UpdatedInput == nil {
			return nil
		}
		decision = sdkcursor.DecisionAllow
	}
	var updated map[string]any
	if res.UpdatedInput != nil && ev.Name == sdkcursor.EventPreToolUse {
		updated = res.UpdatedInput
	}
	return sdkcursor.BuildPermissionOutput(decision, res.Reason, updated)
}
