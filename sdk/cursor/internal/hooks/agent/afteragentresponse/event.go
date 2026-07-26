package afteragentresponse

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/event"
)

// Event is the afterAgentResponse hook event. It is observe-only: Cursor does
// not currently support output fields for this hook.
//
// Cursor's hooks.json matcher for afterAgentResponse is the fixed value
// AgentResponse. That filter is host configuration, not a field on this
// payload.
type Event struct {
	event.Envelope
	// Text is the assistant's final response text.
	Text string `json:"text"`
}

// EventName returns the canonical hook event name.
func (Event) EventName() string { return event.AfterAgentResponse }

// register registers this hook event decoder on c.
func register(c *hookkit.Codec) {
	c.Register(event.AfterAgentResponse, hookkit.EventDecoder[Event](c))
}
