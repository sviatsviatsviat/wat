package copilot

import (
	"context"
	"github.com/sviatsviatsviat/wat/internal/hookkit"
)

// PreCompact is the preCompact hook event.
type PreCompact struct {
	Envelope
	hookkit.RawPayload
	// Trigger is the compaction trigger.
	Trigger string `json:"trigger"`
	// CustomInstructions are user-provided compaction instructions (VS Code).
	CustomInstructions string `json:"custom_instructions"`
	// CustomInstructionsCamel are user-provided compaction instructions (camelCase).
	CustomInstructionsCamel string `json:"customInstructions"`
}

// EventName returns the canonical hook event name.
func (PreCompact) EventName() string { return EventPreCompact }

// Instructions returns custom compaction instructions from either wire format.
func (e PreCompact) Instructions() string {
	if e.CustomInstructions != "" {
		return e.CustomInstructions
	}
	return e.CustomInstructionsCamel
}

func init() {
	registerDecoder(EventPreCompact, decodeAs[PreCompact])
}

// OnPreCompact registers an observe-only preCompact handler.
func OnPreCompact(fn func(context.Context, Hook[PreCompact]) error) *chain {
	return (&chain{}).PreCompact(fn)
}

// PreCompact registers another PreCompact handler on the chain.
func (c *chain) PreCompact(fn func(context.Context, Hook[PreCompact]) error) *chain {
	registerObserveHandler(fn)
	return c
}
