package claude_test

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/claude"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

func ExampleChain() {
	new(claude.Chain).PreToolUse(func(ctx context.Context, hook claude.PreToolUseHook, r claude.PreToolUseResults) (claude.PreToolUseOutput, error) {
		if hook.Event.ToolName == "Bash" {
			return r.Deny("blocked"), nil
		}
		return claude.PreToolUseOutput{}, nil
	})
	// Hook entrypoint: run.Main()
	_ = run.Main
}
