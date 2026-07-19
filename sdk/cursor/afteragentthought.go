package cursor

import (
	"context"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/run"
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
	codec.Register(EventAfterAgentThought, hookkit.EventDecoder[AfterAgentThought](codec))
}

// AfterAgentThought registers a AfterAgentThought handler on the chain.
func (c *chain) AfterAgentThought(fn func(context.Context, run.Hook[AfterAgentThought]) error) *chain {
	registerObserveHandler(c.reg, fn)
	return c
}
