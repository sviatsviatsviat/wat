package cursor

import (
	"context"
)

// AfterAgentResponse is the afterAgentResponse hook event.
type AfterAgentResponse struct {
	Envelope
	// Text is the agent response text.
	Text string `json:"text"`
}

// EventName returns the canonical hook event name.
func (AfterAgentResponse) EventName() string { return EventAfterAgentResponse }

func init() {
	registerDecoder(EventAfterAgentResponse, decodeAs[AfterAgentResponse])
}

// AfterAgentResponse registers an observe-only afterAgentResponse handler.
func (c *Chain) AfterAgentResponse(fn func(context.Context, Hook[AfterAgentResponse]) error) *Chain {
	registerObserveHandler(fn)
	return &Chain{}
}
