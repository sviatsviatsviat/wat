package copilot

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/adapter"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
)

var preToolDecisions = map[model.Decision]sdkcopilot.PermissionDecision{
	model.DecisionAllow: sdkcopilot.DecisionAllow,
	model.DecisionDeny:  sdkcopilot.DecisionDeny,
	model.DecisionAsk:   sdkcopilot.DecisionAsk,
}

func mapPreToolUse(e sdkcopilot.PreToolUse, ev *model.Event) {
	ev.Tool = adapter.NewToolCall(e.NativeToolName(), e.Input().Raw(), "")
}

func mapPreToolOutput(res model.Result) any {
	decision, ok := preToolDecisions[res.Decision]
	if !ok {
		if res.UpdatedInput == nil {
			return nil
		}
		decision = sdkcopilot.DecisionAllow
	}
	return sdkcopilot.BuildPreToolOutput(decision, res.Reason, res.UpdatedInput)
}
