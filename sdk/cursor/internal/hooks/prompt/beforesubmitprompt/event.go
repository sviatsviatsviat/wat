package beforesubmitprompt

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/event"
)

// Event is the beforeSubmitPrompt hook event. Cursor calls it after the user
// hits send and before the backend request. In hooks.json, matchers for this
// event are matched against the value UserPromptSubmit (a config filter string,
// not the wire hook_event_name).
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
