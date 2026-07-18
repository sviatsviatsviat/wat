package copilot_test

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/copilot"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

func ExampleOnPreToolUse() {
	copilot.OnPreToolUse(func(ctx context.Context, hook copilot.Hook[copilot.PreToolUse], r copilot.PreToolResults) (copilot.PreToolOutput, error) {
		if hook.Event.ToolName == "powershell" {
			return r.Deny("blocked"), nil
		}
		return r.Noop(), nil
	})
	// Hook entrypoint: run.Main(run.WithEvent("preToolUse"))
	_ = run.Main
}
