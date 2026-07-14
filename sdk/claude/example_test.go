package claude_test

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/claude"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

func ExampleChain() {
	new(claude.Chain).PreToolUse(func(ctx context.Context, ev claude.PreToolUse, r claude.PreToolUseResults) (claude.PreToolUseOutput, error) {
		if ev.ToolName == "Bash" {
			return r.Deny("blocked"), nil
		}
		return claude.PreToolUseOutput{}, nil
	})
	// Hook entrypoint: run.Main()
	_ = run.Main
}
