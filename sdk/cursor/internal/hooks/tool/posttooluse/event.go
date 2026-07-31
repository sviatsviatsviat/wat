package posttooluse

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/event"
)

// Event is the postToolUse hook event.
type Event struct {
	event.Envelope
	event.ToolFields
	event.DurationFields
	// ToolOutput is the tool output text.
	ToolOutput string `json:"tool_output"`
}

// EventName returns the canonical hook event name.
func (Event) EventName() string { return event.PostToolUse }

// register registers this hook event decoder on c.
func register(c *hookkit.Codec) {
	c.Register(event.PostToolUse, func(raw []byte) (hookkit.Event, error) {
		return hookkit.DecodeEvent(c, raw, func(e *Event, raw []byte) error {
			e.BindToolInput(raw)
			e.CaptureDurationPresent(raw)
			return nil
		})
	})
}
