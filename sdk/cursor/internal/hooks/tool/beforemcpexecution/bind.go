package beforemcpexecution

import (
	"context"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// RegisterHandler registers this event handler on d.
func RegisterHandler(d *hookkit.Dialect, fn func(context.Context, run.Hook[Event], event.PermissionResults) (event.PermissionOutput, error)) {
	hookkit.RegisterWith(d, event.NewPermissionResults(), fn)
}
