package stopfailure

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
)

// Event is the StopFailure hook event.
type Event struct {
	event.Envelope
	// ErrorType is the error category (rate_limit, overloaded, …).
	ErrorType string `json:"error_type"`
	// Message is the error message when provided.
	Message string `json:"message"`
}

// EventName returns the hook event name.
func (Event) EventName() string { return event.StopFailure }

// Register registers this hook event decoder on c.
func Register(c *hookkit.Codec) {
	c.Register(event.StopFailure, hookkit.EventDecoder[Event](c))
}
