package pretooluse

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
)

// Event is the PreToolUse hook event.
type Event struct {
	event.Envelope
	event.ToolFields
}

// EventName returns the hook event name.
func (Event) EventName() string { return event.PreToolUse }

// register registers this hook event decoder on c.
func register(c *hookkit.Codec) {
	c.Register(event.PreToolUse, func(raw []byte) (hookkit.Event, error) {
		return hookkit.DecodeEvent(c, raw, func(e *Event, raw []byte) error {
			e.BindToolInput(raw)
			return nil
		})
	})
}
