package copilot_test

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/copilot"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

func ExampleChain() {
	new(copilot.Chain).PreToolUse(func(ctx context.Context, ev copilot.PreToolUse, r copilot.PreToolResults) (copilot.PreToolOutput, error) {
		if ev.ToolName == "powershell" {
			return r.Deny("blocked"), nil
		}
		return copilot.PreToolOutput{}, nil
	})
	// Hook entrypoint: run.Main(run.WithEvent("preToolUse"))
	_ = run.Main
}
