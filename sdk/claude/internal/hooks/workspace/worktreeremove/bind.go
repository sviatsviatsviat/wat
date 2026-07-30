package worktreeremove

import (
	"context"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
)

// RegisterHandler registers this event handler on d.
func RegisterHandler(d *hookkit.Dialect, fn func(context.Context, Event) error) {
	if fn == nil {
		return
	}
	register(d.Codec())
	hookkit.RegisterObserve(d, fn)
}
