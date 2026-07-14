package copilot_test

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/copilot"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

func ExampleOn() {
	copilot.On(func(ctx context.Context, ev copilot.PreToolUse) (copilot.PreToolOutput, error) {
		if ev.ToolName == "powershell" {
			return copilot.PreToolOutput{Decision: copilot.DecisionDeny, Reason: "blocked"}, nil
		}
		return copilot.PreToolOutput{}, nil
	})
	// Hook entrypoint: run.Main(run.WithEvent("preToolUse"))
	_ = run.Main
}
