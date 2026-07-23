package subagentstop

import (
	"context"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/stop/stopevent"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// RegisterHandler registers this event handler on d.
func RegisterHandler(d *hookkit.Dialect, fn func(context.Context, run.Hook[Event], stopevent.Results) (stopevent.Output, error)) {
	if fn == nil {
		return
	}
	d.Register(hookkit.Handler(func(ctx context.Context, hook run.Hook[Event]) (stopevent.Output, error) {
		return fn(ctx, hook, stopevent.NewResults(event.SubagentStop))
	}))
}

// OnTaskCreated registers a TaskCreated handler on reg.
