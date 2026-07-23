package worktreecreate

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
)

// Event is the WorktreeCreate hook event.
type Event struct {
	event.Envelope
}

// EventName returns the hook event name.
func (Event) EventName() string { return event.WorktreeCreate }

// register registers this hook event decoder on c.
func register(c *hookkit.Codec) {
	c.Register(event.WorktreeCreate, hookkit.EventDecoder[Event](c))
}
