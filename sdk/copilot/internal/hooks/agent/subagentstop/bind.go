package subagentstop

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/hooks/agent/agentstop"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/runtime"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// RegisterHandler registers this event handler on reg.
func RegisterHandler(reg *run.Registry, fn func(context.Context, run.Hook[Event], agentstop.Results) (agentstop.Output, error)) {
	if fn == nil {
		return
	}
	reg.RegisterHandler(runtime.Dialect, run.Handler(func(ctx context.Context, hook run.Hook[Event]) (agentstop.Output, error) {
		return fn(ctx, hook, agentstop.NewResults())
	}))
}
