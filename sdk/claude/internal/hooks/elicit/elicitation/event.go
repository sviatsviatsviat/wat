package elicitation

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
)

// Event is the Elicitation hook event.
type Event struct {
	event.Envelope
	event.ElicitationFields
	// Message is the elicitation prompt shown to the user.
	Message string `json:"message"`
	// URL is the browser authentication URL for url-mode elicitation.
	URL string `json:"url"`
	// RequestedSchema is the requested form input schema JSON when provided.
	RequestedSchema json.RawMessage `json:"requested_schema"`
}

// EventName returns the hook event name.
func (Event) EventName() string { return event.Elicitation }

// register registers this hook event decoder on c.
func register(c *hookkit.Codec) {
	c.Register(event.Elicitation, hookkit.EventDecoder[Event](c))
}
