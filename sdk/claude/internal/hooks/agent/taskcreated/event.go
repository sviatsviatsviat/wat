package taskcreated

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
)

// Event is the TaskCreated hook event.
type Event struct {
	event.Envelope
	// Task is the task payload JSON.
	Task json.RawMessage `json:"task"`
}

// EventName returns the hook event name.
func (Event) EventName() string { return event.TaskCreated }

// register registers this hook event decoder on c.
func register(c *hookkit.Codec) {
	c.Register(event.TaskCreated, hookkit.EventDecoder[Event](c))
}
