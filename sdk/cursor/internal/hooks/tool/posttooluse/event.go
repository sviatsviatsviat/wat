package posttooluse

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/tools"
)

// Event is the postToolUse hook event.
//
// Cursor Hooks docs document tool fields plus duration in milliseconds.
// Outputs may include additional_context and updated_mcp_tool_output.
type Event struct {
	event.Envelope
	// ToolName is the tool name.
	ToolName string `json:"tool_name"`
	// ToolInput is the typed tool input for ToolName.
	ToolInput tools.Input `json:"-"`
	// ToolUseID is the tool use identifier.
	ToolUseID string `json:"tool_use_id"`
	// ToolOutput is the tool output text (JSON-stringified result payload).
	ToolOutput string `json:"tool_output"`
	// Duration is the Cursor Hooks docs duration field in milliseconds.
	Duration int64 `json:"duration"`
	// DurationMs is an alternate duration field some payloads may use.
	DurationMs int64 `json:"duration_ms"`
}

// EventName returns the canonical hook event name.
func (Event) EventName() string { return event.PostToolUse }

// DurationMillis returns the tool execution duration in milliseconds.
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
	c.Register(event.PostToolUse, func(raw []byte) (hookkit.Event, error) {
		return hookkit.DecodeEvent(c, raw, func(e *Event, raw []byte) {
			e.ToolInput = tools.NewInputFromPayload(e.ToolName, raw, "tool_input")
		})
	})
}
