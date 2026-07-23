package posttoolusefailure

import (
	"context"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/tool/posttooluse"
)

// RegisterHandler registers this event handler on d.
func RegisterHandler(d *hookkit.Dialect, fn func(context.Context, Event, Results) (posttooluse.Output, error)) {
	if fn == nil {
		return
	}
	register(d.Codec())
	hookkit.RegisterWith(d, Results(results{}), fn)
}
