package notification

import (
	"context"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// RegisterHandler registers this event handler on d.
func RegisterHandler(d *hookkit.Dialect, fn func(context.Context, run.Hook[Event], Results) (event.CommonOutput, error)) {
	if fn == nil {
		return
	}
	d.Register(hookkit.Handler(func(ctx context.Context, hook run.Hook[Event]) (event.CommonOutput, error) {
		return fn(ctx, hook, results{})
	}))
}

// OnMessageDisplay registers a MessageDisplay handler on reg.
