package posttooluse

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/event"
)

// Event is the PostToolUse hook event.
type Event struct {
	event.Envelope
	event.ToolFields
	// ToolResult is the tool result.
	ToolResult event.ToolResult `json:"tool_result"`
}

// EventName returns the canonical hook event name.
func (Event) EventName() string { return event.PostToolUse }

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

// register registers this hook event decoder on c.
func register(c *hookkit.Codec) {
	c.Register(event.PostToolUse, func(raw []byte) (hookkit.Event, error) {
		return hookkit.DecodeEvent(c, raw, func(e *Event, payload []byte) {
			e.BindToolInput(payload)
		})
	})
}
