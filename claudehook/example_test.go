package claudehook_test

import (
	"context"

	"github.com/sviatsviatsviat/wat/claudehook"
)

func ExampleMux() {
	mux := claudehook.NewMux()
	claudehook.On(mux, func(ctx context.Context, ev claudehook.PreToolUse) (claudehook.PreToolUseOutput, error) {
		if ev.ToolName == "Bash" {
			return claudehook.PreToolUseOutput{Decision: claudehook.DecisionDeny, Reason: "blocked"}, nil
		}
		return claudehook.PreToolUseOutput{}, nil
	})
	// Hook entrypoint: mux.Main()
}
