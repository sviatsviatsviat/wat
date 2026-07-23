package sessionend

import (
	"context"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/event"
)

// Event is the sessionEnd hook event.
type Event struct {
	event.Envelope
	// Reason is the session end reason.
	Reason string `json:"reason"`
	// IsBackgroundAgent reports whether this is a background agent session.
	IsBackgroundAgent bool `json:"is_background_agent"`
}

// EventName returns the canonical hook event name.
func (Event) EventName() string { return event.SessionEnd }

// register registers this hook event decoder on c.
func register(c *hookkit.Codec) {
	c.Register(event.SessionEnd, hookkit.EventDecoder[Event](c))
}

// RegisterHandler registers a SessionEnd observe handler on reg.
func RegisterHandler(d *hookkit.Dialect, fn func(context.Context, Event) error) {
	if fn == nil {
		return
	}
	register(d.Codec())
	hookkit.RegisterObserve(d, fn)
}
