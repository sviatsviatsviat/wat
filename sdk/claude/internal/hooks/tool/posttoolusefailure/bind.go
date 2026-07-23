package posttoolusefailure

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/tool/posttooluse"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/runtime"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// RegisterHandler registers a PostToolUseFailure handler on reg.
func RegisterHandler(reg *run.Registry, fn func(context.Context, run.Hook[Event], Results) (posttooluse.Output, error)) {
	if fn == nil {
		return
	}
	reg.RegisterHandler(runtime.Dialect, run.Handler(func(ctx context.Context, hook run.Hook[Event]) (posttooluse.Output, error) {
		return fn(ctx, hook, results{})
	}))
}
