package aftershellexecution

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/runtime"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// RegisterHandler registers this event handler on reg.
func RegisterHandler(reg *run.Registry, fn func(context.Context, run.Hook[Event], event.PostToolResults) (event.PostToolOutput, error)) {
	if fn == nil {
		return
	}
	reg.RegisterHandler(runtime.Dialect, run.Handler(func(ctx context.Context, hook run.Hook[Event]) (event.PostToolOutput, error) {
		return fn(ctx, hook, event.NewPostToolResults())
	}))
}
