package userpromptsubmit

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/claude/internal/runtime"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// RegisterHandler registers this event handler on reg.
func RegisterHandler(reg *run.Registry, fn func(context.Context, run.Hook[Event], Results) (Output, error)) {
	if fn == nil {
		return
	}
	reg.RegisterHandler(runtime.Dialect, run.Handler(func(ctx context.Context, hook run.Hook[Event]) (Output, error) {
		return fn(ctx, hook, results{})
	}))
}

// OnUserPromptExpansion registers a UserPromptExpansion handler on reg.
