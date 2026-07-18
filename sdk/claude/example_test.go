package claude_test

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/claude"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

func ExampleOnPreToolUse() {
	claude.OnPreToolUse(func(ctx context.Context, hook claude.Hook[claude.PreToolUse], r claude.PreToolUseResults) (claude.PreToolUseOutput, error) {
		if _, ok := hook.Event.ToolInput.AsBash(); ok {
			return r.Deny("blocked"), nil
		}
		return nil, nil
	})
	// Hook entrypoint: run.Main()
	_ = run.Main
}
