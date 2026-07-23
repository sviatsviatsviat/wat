package worktreecreate

import (
	"context"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
)

// RegisterHandler registers this event handler on d.
func RegisterHandler(d *hookkit.Dialect, fn func(context.Context, Event, Results) (Output, error)) {
	if fn == nil {
		return
	}
	register(d.Codec())
	hookkit.RegisterWith(d, Results(results{}), fn)
}
