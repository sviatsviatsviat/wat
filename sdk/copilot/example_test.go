package copilot_test

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/copilot"
)

func ExampleMux() {
	mux := copilot.NewMux()
	copilot.On(mux, func(ctx context.Context, ev copilot.PreToolUse) (copilot.PreToolOutput, error) {
		if ev.NativeToolName() == "bash" {
			return copilot.PreToolOutput{Decision: copilot.DecisionDeny, Reason: "blocked"}, nil
		}
		return copilot.PreToolOutput{}, nil
	})
	// Hook entrypoint: mux.Main(copilot.WithEvent("preToolUse"))
}
