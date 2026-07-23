package afteragentresponse

import (
	"context"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// Event is the afterAgentResponse hook event.
type Event struct {
	event.Envelope
	// Text is the agent response text.
	Text string `json:"text"`
}

// EventName returns the canonical hook event name.
func (Event) EventName() string { return event.AfterAgentResponse }

// Register registers this hook event decoder on c.
func Register(c *hookkit.Codec) {
	c.Register(event.AfterAgentResponse, hookkit.EventDecoder[Event](c))
}

// RegisterHandler registers an observe handler on reg.
func RegisterHandler(d *hookkit.Dialect, fn func(context.Context, run.Hook[Event]) error) {
	hookkit.RegisterObserve(d, fn)
}
