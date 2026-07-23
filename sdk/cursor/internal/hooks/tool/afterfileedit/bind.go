package afterfileedit

import (
	"context"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// RegisterHandler registers this event handler on d.
func RegisterHandler(d *hookkit.Dialect, fn func(context.Context, run.Hook[Event], event.PostToolResults) (event.PostToolOutput, error)) {
	if fn == nil {
		return
	}
	d.Register(hookkit.Handler(func(ctx context.Context, hook run.Hook[Event]) (event.PostToolOutput, error) {
		return fn(ctx, hook, event.NewPostToolResults())
	}))
}
