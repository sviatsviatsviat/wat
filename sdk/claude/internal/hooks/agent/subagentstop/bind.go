package subagentstop

import (
	"context"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/stop/stopevent"
)

// RegisterHandler registers this event handler on d.
func RegisterHandler(d *hookkit.Dialect, fn func(context.Context, Event, stopevent.Results) (stopevent.Output, error)) {
	if fn == nil {
		return
	}
	register(d.Codec())
	hookkit.RegisterWith(d, stopevent.NewResults(event.SubagentStop), fn)
}
