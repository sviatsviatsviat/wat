package elicitationresult

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
)

// Event is the ElicitationResult hook event.
type Event struct {
	event.Envelope
	event.ElicitationFields
	// Action is the user action (accept, decline, cancel).
	Action string `json:"action"`
	// Content is the response content JSON when provided.
	Content json.RawMessage `json:"content"`
}

// EventName returns the hook event name.
func (Event) EventName() string { return event.ElicitationResult }

// register registers this hook event decoder on c.
func register(c *hookkit.Codec) {
	c.Register(event.ElicitationResult, hookkit.EventDecoder[Event](c))
}
