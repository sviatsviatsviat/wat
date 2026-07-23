package userpromptexpansion

import (
	"context"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
)

// RegisterHandler registers this event handler on d.
func RegisterHandler(d *hookkit.Dialect, fn func(context.Context, Event, Results) (event.CommonOutput, error)) {
	hookkit.RegisterWith(d, Results(results{}), fn)
}
