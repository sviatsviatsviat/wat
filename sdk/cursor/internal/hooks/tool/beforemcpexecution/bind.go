package beforemcpexecution

import (
	"context"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/event"
)

// RegisterHandler registers this event handler on d.
func RegisterHandler(d *hookkit.Dialect, fn func(context.Context, Event, event.PermissionResults) (event.PermissionOutput, error)) {
	if fn == nil {
		return
	}
	register(d.Codec())
	hookkit.RegisterWith(d, event.NewPermissionResults(), fn)
}
