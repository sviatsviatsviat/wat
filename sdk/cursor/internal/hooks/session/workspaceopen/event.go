package workspaceopen

import (
	"context"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/event"
)

// Event is the workspaceOpen hook event.
type Event struct {
	event.Envelope
}

// EventName returns the canonical hook event name.
func (Event) EventName() string { return event.WorkspaceOpen }

// register registers this hook event decoder on c.
func register(c *hookkit.Codec) {
	c.Register(event.WorkspaceOpen, hookkit.EventDecoder[Event](c))
}

// RegisterHandler registers a WorkspaceOpen observe handler on reg.
func RegisterHandler(d *hookkit.Dialect, fn func(context.Context, Event) error) {
	if fn == nil {
		return
	}
	register(d.Codec())
	hookkit.RegisterObserve(d, fn)
}
