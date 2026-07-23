package claude_test

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/claude"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

func ExampleUseHooks() {
	claude.UseHooks().PreToolUse(func(ctx context.Context, hook claude.PreToolUse, r claude.PreToolUseResults) (claude.PreToolUseOutput, error) {
		if _, ok := hook.ToolInput.AsBash(); ok {
			return r.Deny("blocked"), nil
		}
		return r.Noop(), nil
	})
	// Hook entrypoint: run.Main()
	_ = run.Main
}
