package stopevent

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/event"
)

// Event is the stop hook event.
type Event struct {
	event.Envelope
	// Status is the stop status.
	Status string `json:"status"`
	// LoopCount is the stop-loop iteration count.
	LoopCount int `json:"loop_count"`
}

// EventName returns the canonical hook event name.
func (Event) EventName() string { return event.Stop }

// register registers this hook event decoder on c.
func register(c *hookkit.Codec) {
	c.Register(event.Stop, hookkit.EventDecoder[Event](c))
}
