package afteragentresponse

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/event"
)

// Event is the afterAgentResponse hook event.
type Event struct {
	event.Envelope
	// Text is the agent response text.
	Text string `json:"text"`
}

// EventName returns the canonical hook event name.
func (Event) EventName() string { return event.AfterAgentResponse }

// register registers this hook event decoder on c.
func register(c *hookkit.Codec) {
	c.Register(event.AfterAgentResponse, hookkit.EventDecoder[Event](c))
}
