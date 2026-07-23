package precompact

import (
	"context"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// Event is the PreCompact hook event.
type Event struct {
	event.Envelope
	// Trigger is the compaction trigger.
	Trigger string `json:"trigger"`
	// CustomInstructions are user-provided compaction instructions.
	CustomInstructions string `json:"custom_instructions"`
}

// EventName returns the canonical hook event name.
func (Event) EventName() string { return event.PreCompact }

// Instructions returns custom compaction instructions.
func (e Event) Instructions() string {
	return e.CustomInstructions
}

// Register registers this hook event decoder on c.
func Register(c *hookkit.Codec) {
	c.Register(event.PreCompact, hookkit.EventDecoder[Event](c))
}

// RegisterHandler registers a PreCompact observe handler on d.
func RegisterHandler(d *hookkit.Dialect, fn func(context.Context, run.Hook[Event]) error) {
	if fn == nil {
		return
	}
	d.Register(hookkit.ObserveHandler(fn))
}
