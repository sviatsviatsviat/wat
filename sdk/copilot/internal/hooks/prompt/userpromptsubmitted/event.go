package userpromptsubmitted

import (
	"context"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/event"
)

// Event is the UserPromptSubmitted hook event.
type Event struct {
	event.Envelope
	// Prompt is the submitted user prompt text.
	Prompt string `json:"prompt"`
}

// EventName returns the canonical hook event name.
func (Event) EventName() string { return event.UserPromptSubmitted }

// register registers this hook event decoder on c.
func register(c *hookkit.Codec) {
	c.Register(event.UserPromptSubmitted, hookkit.EventDecoder[Event](c))
}

// RegisterHandler registers a UserPromptSubmitted observe handler on reg.
func RegisterHandler(d *hookkit.Dialect, fn func(context.Context, Event) error) {
	if fn == nil {
		return
	}
	register(d.Codec())
	hookkit.RegisterObserve(d, fn)
}
