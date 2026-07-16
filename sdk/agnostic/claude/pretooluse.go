package claude

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/adapter"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
)

var preToolDecisions = map[model.Decision]sdkclaude.PermissionDecision{
	model.DecisionAllow: sdkclaude.DecisionAllow,
	model.DecisionDeny:  sdkclaude.DecisionDeny,
	model.DecisionAsk:   sdkclaude.DecisionAsk,
}

func mapPreToolUse(e sdkclaude.PreToolUse, ev *model.Event) {
	ev.Tool = adapter.NewToolCall(e.ToolName, e.ToolInput.Raw(), e.ToolUseID)
}

func mapPreToolOutput(res model.Result) any {
	decision, ok := preToolDecisions[res.Decision]
	if !ok {
		if res.UpdatedInput == nil {
			return nil
		}
		decision = sdkclaude.DecisionAllow
	}
	return sdkclaude.BuildPreToolUseOutput(decision, res.Reason, res.UpdatedInput)
}
