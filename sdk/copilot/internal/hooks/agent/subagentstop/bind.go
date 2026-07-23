package subagentstop

import (
	"context"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/hooks/agent/agentstop"
)

// RegisterHandler registers this event handler on d.
func RegisterHandler(d *hookkit.Dialect, fn func(context.Context, Event, agentstop.Results) (agentstop.Output, error)) {
	if fn == nil {
		return
	}
	register(d.Codec())
	hookkit.RegisterWith(d, agentstop.NewResults(), fn)
}
