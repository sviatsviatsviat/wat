package cursor

import (
	"context"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/run"
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
	codec.Register(EventAfterAgentResponse, hookkit.EventDecoder[AfterAgentResponse](codec))
}

// OnAfterAgentResponse registers an observe-only afterAgentResponse handler.
func OnAfterAgentResponse(fn func(context.Context, run.Hook[AfterAgentResponse]) error) *chain {
	return (&chain{}).AfterAgentResponse(fn)
}

// AfterAgentResponse registers another AfterAgentResponse handler on the chain.
func (c *chain) AfterAgentResponse(fn func(context.Context, run.Hook[AfterAgentResponse]) error) *chain {
	registerObserveHandler(fn)
	return c
}
