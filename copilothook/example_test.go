package copilothook_test

import (
	"context"

	"github.com/sviatsviatsviat/wat/copilothook"
)

func ExampleMux() {
	mux := copilothook.NewMux()
	copilothook.On(mux, func(ctx context.Context, ev copilothook.PreToolUse) (copilothook.PreToolOutput, error) {
		if ev.NativeToolName() == "bash" {
			return copilothook.PreToolOutput{Decision: copilothook.DecisionDeny, Reason: "blocked"}, nil
		}
		return copilothook.PreToolOutput{}, nil
	})
	// Hook entrypoint: mux.Main(copilothook.WithEvent("preToolUse"))
}
