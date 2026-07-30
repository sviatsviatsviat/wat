package stopfailure

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
)

// Event is the StopFailure hook event.
type Event struct {
	event.Envelope
	// Error is the error category used for matcher filtering (rate_limit, overloaded, …).
	Error string `json:"error"`
	// ErrorDetails holds additional error detail when available.
	ErrorDetails string `json:"error_details"`
	// LastAssistantMessage is the rendered API error text shown in the conversation.
	LastAssistantMessage string `json:"last_assistant_message"`
}

// EventName returns the hook event name.
func (Event) EventName() string { return event.StopFailure }

// register registers this hook event decoder on c.
func register(c *hookkit.Codec) {
	c.Register(event.StopFailure, hookkit.EventDecoder[Event](c))
}
