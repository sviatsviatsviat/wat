package posttoolusefailure

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/tools"
)

// Event is the PostToolUseFailure hook event.
type Event struct {
	event.Envelope
	// ToolName is the tool name.
	ToolName string `json:"tool_name"`
	// ToolInput is the typed tool input.
	ToolInput tools.Input `json:"-"`
	// Error is the failure payload (string or object).
	Error json.RawMessage `json:"error"`
}

// EventName returns the canonical hook event name.
func (Event) EventName() string { return event.PostToolUseFailure }

// NativeToolName returns the tool name.
func (e Event) NativeToolName() string {
	return e.ToolName
}

// Input returns tool input.
func (e Event) Input() tools.Input {
	return e.ToolInput
}

// ErrorMessage returns the failure message from the error field.
func (e Event) ErrorMessage() string {
	if len(e.Error) == 0 {
		return ""
	}
	var s string
	if json.Unmarshal(e.Error, &s) == nil {
		return s
	}
	var detail event.ErrorDetail
	if json.Unmarshal(e.Error, &detail) == nil {
		return detail.Message
	}
	return string(e.Error)
}

// register registers this hook event decoder on c.
func register(c *hookkit.Codec) {
	c.Register(event.PostToolUseFailure, func(raw []byte) (hookkit.Event, error) {
		return hookkit.DecodeEvent(c, raw, func(e *Event, payload []byte) {
			e.ToolInput = tools.NewInputFromPayload(e.ToolName, payload, "tool_input")
		})
	})
}
