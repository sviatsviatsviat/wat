package aftermcpexecution

import (
	"context"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/event"
)

// RegisterHandler registers this event handler on d.
func RegisterHandler(d *hookkit.Dialect, fn func(context.Context, Event, event.PostToolResults) (event.PostToolOutput, error)) {
	hookkit.RegisterWith(d, event.NewPostToolResults(), fn)
}
