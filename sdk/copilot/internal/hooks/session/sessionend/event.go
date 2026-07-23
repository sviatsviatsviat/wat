package sessionend

import (
	"context"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// Event is the SessionEnd hook event.
type Event struct {
	event.Envelope
	// Reason is the session end reason.
	Reason string `json:"reason"`
}

// EventName returns the canonical hook event name.
func (Event) EventName() string { return event.SessionEnd }

// Register registers this hook event decoder on c.
func Register(c *hookkit.Codec) {
	c.Register(event.SessionEnd, hookkit.EventDecoder[Event](c))
}

// RegisterHandler registers a SessionEnd observe handler on reg.
func RegisterHandler(d *hookkit.Dialect, fn func(context.Context, run.Hook[Event]) error) {
	if fn == nil {
		return
	}
	d.Register(hookkit.ObserveHandler(fn))
}
