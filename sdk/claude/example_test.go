package claude_test

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/claude"
)

func ExampleMux() {
	mux := claude.NewMux()
	claude.On(mux, func(ctx context.Context, ev claude.PreToolUse) (claude.PreToolUseOutput, error) {
		if ev.ToolName == "Bash" {
			return claude.PreToolUseOutput{Decision: claude.DecisionDeny, Reason: "blocked"}, nil
		}
		return claude.PreToolUseOutput{}, nil
	})
	// Hook entrypoint: mux.Main()
}
