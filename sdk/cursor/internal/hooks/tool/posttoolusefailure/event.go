package posttoolusefailure

import (
	"context"

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
	// Duration is the Cursor Hooks docs duration field in milliseconds.
	Duration int64 `json:"duration"`
	// DurationMs is an alternate duration field some payloads may use.
	DurationMs int64 `json:"duration_ms"`
	// IsInterrupt is true when the failure was caused by a user interrupt
	// or cancellation.
	IsInterrupt bool `json:"is_interrupt"`
}

// EventName returns the canonical hook event name.
func (Event) EventName() string { return event.PostToolUseFailure }

// DurationMillis returns the failure duration in milliseconds.
// Prefer this helper over reading Duration or DurationMs directly: Cursor
// Hooks docs use `duration`, and DurationMillis falls back to `duration_ms`
// when `duration` is zero so alternate wire forms still decode.
func (e Event) DurationMillis() int64 {
	if e.Duration != 0 {
		return e.Duration
	}
	return e.DurationMs
}

// register registers this hook event decoder on c.
func register(c *hookkit.Codec) {
	c.Register(event.PostToolUseFailure, func(raw []byte) (hookkit.Event, error) {
		return hookkit.DecodeEvent(c, raw, func(e *Event, raw []byte) {
			e.ToolInput = tools.NewInputFromPayload(e.ToolName, raw, "tool_input")
		})
	})
}

// RegisterHandler registers a PostToolUseFailure observe handler on d.
func RegisterHandler(d *hookkit.Dialect, fn func(context.Context, Event) error) {
	if fn == nil {
		return
	}
	register(d.Codec())
	hookkit.RegisterObserve(d, fn)
}
