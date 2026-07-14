package claude_test

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/claude"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

func ExampleOn() {
	claude.On(func(ctx context.Context, ev claude.PreToolUse) (claude.PreToolUseOutput, error) {
		if ev.ToolName == "Bash" {
			return claude.PreToolUseOutput{Decision: claude.DecisionDeny, Reason: "blocked"}, nil
		}
		return claude.PreToolUseOutput{}, nil
	})
	// Hook entrypoint: run.Main()
	_ = run.Main
}
