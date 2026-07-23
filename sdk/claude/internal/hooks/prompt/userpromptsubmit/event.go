package userpromptsubmit

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
)

// Event is the UserPromptSubmit hook event.
type Event struct {
	event.Envelope
	// Prompt is the submitted user prompt text.
	Prompt string `json:"prompt"`
}

// EventName returns the hook event name.
func (Event) EventName() string { return event.UserPromptSubmit }

// register registers this hook event decoder on c.
func register(c *hookkit.Codec) {
	c.Register(event.UserPromptSubmit, hookkit.EventDecoder[Event](c))
}
