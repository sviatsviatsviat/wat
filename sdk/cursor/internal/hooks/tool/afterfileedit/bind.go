package afterfileedit

import (
	"context"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
)

// RegisterHandler registers an observe-only AfterFileEdit handler on d.
// Cursor documents no afterFileEdit output fields; handlers return only an error.
func RegisterHandler(d *hookkit.Dialect, fn func(context.Context, Event) error) {
	if fn == nil {
		return
	}
	register(d.Codec())
	hookkit.RegisterObserve(d, fn)
}
