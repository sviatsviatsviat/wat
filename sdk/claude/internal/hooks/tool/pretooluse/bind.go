package pretooluse

import (
	"context"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// RegisterHandler registers this event handler on d.
func RegisterHandler(d *hookkit.Dialect, fn func(context.Context, run.Hook[Event], Results) (Output, error)) {
	if fn == nil {
		return
	}
	d.Register(hookkit.Handler(func(ctx context.Context, hook run.Hook[Event]) (Output, error) {
		return fn(ctx, hook, results{})
	}))
}

// OnPostToolUse registers a PostToolUse handler on reg.
