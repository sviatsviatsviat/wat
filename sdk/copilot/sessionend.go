package copilot

import (
	"context"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// SessionEnd is the sessionEnd hook event.
type SessionEnd struct {
	Envelope
	// Reason is the session end reason.
	Reason string `json:"reason"`
}

// EventName returns the canonical hook event name.
func (SessionEnd) EventName() string { return EventSessionEnd }

func init() {
	codec.Register(EventSessionEnd, hookkit.EventDecoder[SessionEnd](codec))
}

// OnSessionEnd registers an observe-only SessionEnd handler.
func OnSessionEnd(fn func(context.Context, run.Hook[SessionEnd]) error) *chain {
	return (&chain{}).SessionEnd(fn)
}

// SessionEnd registers another SessionEnd handler on the chain.
func (c *chain) SessionEnd(fn func(context.Context, run.Hook[SessionEnd]) error) *chain {
	registerObserveHandler(fn)
	return c
}
