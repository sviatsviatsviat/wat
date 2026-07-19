package copilot

import (
	"context"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// PreCompact is the PreCompact hook event.
type PreCompact struct {
	Envelope
	// Trigger is the compaction trigger.
	Trigger string `json:"trigger"`
	// CustomInstructions are user-provided compaction instructions.
	CustomInstructions string `json:"custom_instructions"`
}

// EventName returns the canonical hook event name.
func (PreCompact) EventName() string { return EventPreCompact }

// Instructions returns custom compaction instructions.
func (e PreCompact) Instructions() string {
	return e.CustomInstructions
}

func init() {
	codec.Register(EventPreCompact, hookkit.EventDecoder[PreCompact](codec))
}

// PreCompact registers a PreCompact handler on the chain.
func (c *chain) PreCompact(fn func(context.Context, run.Hook[PreCompact]) error) *chain {
	registerObserveHandler(c.reg, fn)
	return c
}
