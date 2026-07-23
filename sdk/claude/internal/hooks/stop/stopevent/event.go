package stopevent

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
)

// Event is the Stop hook event.
type Event struct {
	event.Envelope
	// StopHookActive is true when a stop hook already forced continuation.
	StopHookActive bool `json:"stop_hook_active"`
	// LastAssistantMessage is the final assistant text of the turn.
	LastAssistantMessage string `json:"last_assistant_message"`
}

// EventName returns the hook event name.
func (Event) EventName() string { return event.Stop }

// register registers this hook event decoder on c.
func register(c *hookkit.Codec) {
	c.Register(event.Stop, hookkit.EventDecoder[Event](c))
}
