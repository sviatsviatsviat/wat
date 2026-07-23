package beforesubmitprompt

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/event"
)

// Event is the beforeSubmitPrompt hook event.
type Event struct {
	event.Envelope
	// Prompt is the user prompt about to be submitted.
	Prompt string `json:"prompt"`
	// Attachments are context attachments associated with the prompt.
	Attachments []event.Attachment `json:"attachments"`
}

// EventName returns the canonical hook event name.
func (Event) EventName() string { return event.BeforeSubmitPrompt }

// register registers this hook event decoder on c.
func register(c *hookkit.Codec) {
	c.Register(event.BeforeSubmitPrompt, hookkit.EventDecoder[Event](c))
}
