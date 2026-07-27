package posttoolusefailure

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/event"
)

// Event is the postToolUseFailure hook event.
//
// Cursor Hooks docs list input fields only (error_message, failure_type,
// duration, is_interrupt); the host does not document consumed stdout JSON for
// this event. Handlers are observe-only.
type Event struct {
	event.Envelope
	event.ToolFields
	event.DurationFields
	// ErrorMessage is the failure message.
	ErrorMessage string `json:"error_message"`
	// FailureType is the failure category.
	FailureType string `json:"failure_type"`
	// IsInterrupt is true when the failure was caused by a user interrupt
	// or cancellation.
	IsInterrupt bool `json:"is_interrupt"`
}

// EventName returns the canonical hook event name.
func (Event) EventName() string { return event.PostToolUseFailure }

// register registers this hook event decoder on c.
func register(c *hookkit.Codec) {
	c.Register(event.PostToolUseFailure, func(raw []byte) (hookkit.Event, error) {
		return hookkit.DecodeEvent(c, raw, func(e *Event, raw []byte) {
			e.BindToolInput(raw)
			e.CaptureDurationPresent(raw)
		})
	})
}
