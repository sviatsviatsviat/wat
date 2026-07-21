package claude

import (
	"context"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// SessionEnd is the SessionEnd hook event.
type SessionEnd struct {
	Envelope
	// Reason is the session end reason.
	Reason string `json:"reason"`
}

// EventName returns the hook event name.
func (SessionEnd) EventName() string { return EventSessionEnd }

func init() {
	codec.Register(EventSessionEnd, hookkit.EventDecoder[SessionEnd](codec))
}

// SessionEnd registers a SessionEnd handler on the chain.
func (c *chain) SessionEnd(fn func(context.Context, run.Hook[SessionEnd]) error) *chain {
	c.reg.RegisterObserveHandler(Dialect, run.ObserveHandler(fn))
	return c
}
