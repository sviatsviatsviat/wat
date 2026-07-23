package messagedisplay

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
)

// Event is the MessageDisplay hook event.
type Event struct {
	event.Envelope
	// TurnID is the turn identifier.
	TurnID string `json:"turn_id"`
	// MessageID is the message identifier.
	MessageID string `json:"message_id"`
	// Index is the message index in the turn.
	Index int `json:"index"`
	// Final is true when this is the final delta.
	Final bool `json:"final"`
	// Delta is the streamed message delta.
	Delta string `json:"delta"`
}

// EventName returns the hook event name.
func (Event) EventName() string { return event.MessageDisplay }

// Register registers this hook event decoder on c.
func Register(c *hookkit.Codec) {
	c.Register(event.MessageDisplay, hookkit.EventDecoder[Event](c))
}
