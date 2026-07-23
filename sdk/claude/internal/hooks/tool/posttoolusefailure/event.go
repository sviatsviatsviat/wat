package posttoolusefailure

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/tools"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// Event is the PostToolUseFailure hook event.
type Event struct {
	event.Envelope
	// ToolName is the tool name.
	ToolName string `json:"tool_name"`
	// ToolInput is the typed tool input for ToolName.
	ToolInput tools.Input `json:"-"`
	// ToolUseID is the tool use identifier.
	ToolUseID string `json:"tool_use_id"`
	// Error is the failure message.
	Error string `json:"error"`
	// IsInterrupt is true when the failure was caused by an interrupt.
	IsInterrupt bool `json:"is_interrupt"`
	// DurationMs is the tool execution duration in milliseconds.
	DurationMs int64 `json:"duration_ms"`
}

// EventName returns the hook event name.
func (Event) EventName() string { return event.PostToolUseFailure }

// Register registers this hook event decoder on c.
func Register(c *hookkit.Codec) {
	c.Register(event.PostToolUseFailure, func(raw []byte) (run.Event, error) {
		return hookkit.DecodeEvent(c, raw, func(e *Event, raw []byte) {
			e.ToolInput = tools.NewInputFromPayload(e.ToolName, raw, "tool_input")
		})
	})
}
