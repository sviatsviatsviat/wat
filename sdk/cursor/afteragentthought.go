package cursor

import (
	"context"
	"github.com/sviatsviatsviat/wat/internal/hookkit"
)

// AfterAgentThought is the afterAgentThought hook event.
type AfterAgentThought struct {
	Envelope
	hookkit.RawPayload
	// Text is the agent thought text.
	Text string `json:"text"`
}

// EventName returns the canonical hook event name.
func (AfterAgentThought) EventName() string { return EventAfterAgentThought }

func init() {
	registerDecoder(EventAfterAgentThought, decodeAs[AfterAgentThought])
}

// OnAfterAgentThought registers an observe-only afterAgentThought handler.
func OnAfterAgentThought(fn func(context.Context, Hook[AfterAgentThought]) error) *chain {
	return (&chain{}).AfterAgentThought(fn)
}

// AfterAgentThought registers another AfterAgentThought handler on the chain.
func (c *chain) AfterAgentThought(fn func(context.Context, Hook[AfterAgentThought]) error) *chain {
	registerObserveHandler(fn)
	return c
}
