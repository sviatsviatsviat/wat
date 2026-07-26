package afteragentthought

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/event"
)

// Event is the afterAgentThought hook event.
type Event struct {
	event.Envelope
	// Text is the agent thought text.
	Text string `json:"text"`
}

// EventName returns the canonical hook event name.
func (Event) EventName() string { return event.AfterAgentThought }

// register registers this hook event decoder on c.
func register(c *hookkit.Codec) {
	c.Register(event.AfterAgentThought, hookkit.EventDecoder[Event](c))
}
