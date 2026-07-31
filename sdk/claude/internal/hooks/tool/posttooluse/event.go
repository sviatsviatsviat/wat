package posttooluse

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
)

// Event is the PostToolUse hook event.
type Event struct {
	event.Envelope
	event.ToolFields
	// ToolResponse is the tool response JSON.
	ToolResponse json.RawMessage `json:"tool_response"`
	// DurationMs is the tool execution duration in milliseconds.
	DurationMs int64 `json:"duration_ms"`
}

// EventName returns the hook event name.
func (Event) EventName() string { return event.PostToolUse }

// register registers this hook event decoder on c.
func register(c *hookkit.Codec) {
	c.Register(event.PostToolUse, func(raw []byte) (hookkit.Event, error) {
		return hookkit.DecodeEvent(c, raw, func(e *Event, raw []byte) error {
			e.BindToolInput(raw)
			return nil
		})
	})
}
