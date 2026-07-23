package posttooluse

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/tools"
)

// Event is the PostToolUse hook event.
type Event struct {
	event.Envelope
	// ToolName is the tool name.
	ToolName string `json:"tool_name"`
	// ToolInput is the typed tool input.
	ToolInput tools.Input `json:"-"`
	// ToolResult is the tool result.
	ToolResult event.ToolResult `json:"tool_result"`
}

// EventName returns the canonical hook event name.
func (Event) EventName() string { return event.PostToolUse }

// NativeToolName returns the tool name.
func (e Event) NativeToolName() string {
	return e.ToolName
}

// Input returns tool input.
func (e Event) Input() tools.Input {
	return e.ToolInput
}

// ResultText returns the textual tool result.
func (e Event) ResultText() string {
	return e.ToolResult.Text()
}

// ResultRaw returns the tool result JSON.
func (e Event) ResultRaw() json.RawMessage {
	if e.ToolResult.TextResultForLLM != "" || e.ToolResult.ResultType != "" {
		return event.MarshalToolResult(e.ToolResult)
	}
	return nil
}

// Register registers this hook event decoder on c.
func Register(c *hookkit.Codec) {
	c.Register(event.PostToolUse, func(raw []byte) (hookkit.Event, error) {
		return hookkit.DecodeEvent(c, raw, func(e *Event, payload []byte) {
			e.ToolInput = tools.NewInputFromPayload(e.ToolName, payload, "tool_input")
		})
	})
}
