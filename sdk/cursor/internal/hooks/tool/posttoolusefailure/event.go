package posttoolusefailure

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/tools"
)

// Event is the postToolUseFailure hook event.
//
// Cursor Hooks docs list input fields only (error_message, failure_type,
// duration, is_interrupt); the host does not document consumed stdout JSON for
// this event. Handlers are observe-only.
type Event struct {
	event.Envelope
	event.DurationFields
	// ToolName is the tool name.
	ToolName string `json:"tool_name"`
	// ToolInput is the typed tool input for ToolName.
	ToolInput tools.Input `json:"-"`
	// ToolUseID is the tool use identifier.
	ToolUseID string `json:"tool_use_id"`
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
			e.ToolInput = tools.NewInputFromPayload(e.ToolName, raw, "tool_input")
			e.CaptureDurationPresent(raw)
		})
	})
}
