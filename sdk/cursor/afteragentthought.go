package cursor

import (
	"context"
)

// AfterAgentThought is the afterAgentThought hook event.
type AfterAgentThought struct {
	Envelope
	// Text is the agent thought text.
	Text string `json:"text"`
}

// EventName returns the canonical hook event name.
func (AfterAgentThought) EventName() string { return EventAfterAgentThought }

func init() {
	registerDecoder(EventAfterAgentThought, decodeAs[AfterAgentThought])
}

// AfterAgentThought registers an observe-only afterAgentThought handler.
func (c *Chain) AfterAgentThought(fn func(context.Context, Hook[AfterAgentThought]) error) *Chain {
	registerObserveHandler(c.registerOwner(), fn)
	return c
}
