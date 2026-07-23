package subagentstop

import (
	"context"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/hooks/agent/agentstop"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// RegisterHandler registers this event handler on d.
func RegisterHandler(d *hookkit.Dialect, fn func(context.Context, run.Hook[Event], agentstop.Results) (agentstop.Output, error)) {
	if fn == nil {
		return
	}
	d.Register(hookkit.Handler(func(ctx context.Context, hook run.Hook[Event]) (agentstop.Output, error) {
		return fn(ctx, hook, agentstop.NewResults())
	}))
}
