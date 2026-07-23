package subagentstop

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/agent/stopevent"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/runtime"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// RegisterHandler registers this event handler on reg.
func RegisterHandler(reg *run.Registry, fn func(context.Context, run.Hook[Event], stopevent.Results) (stopevent.Output, error)) {
	if fn == nil {
		return
	}
	reg.RegisterHandler(runtime.Dialect, run.Handler(func(ctx context.Context, hook run.Hook[Event]) (stopevent.Output, error) {
		return fn(ctx, hook, stopevent.NewResults())
	}))
}
