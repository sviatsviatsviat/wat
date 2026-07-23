package worktreeremove

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
)

// Event is the WorktreeRemove hook event.
type Event struct {
	event.Envelope
	// WorktreePath is the worktree path being removed.
	WorktreePath string `json:"worktree_path"`
}

// EventName returns the hook event name.
func (Event) EventName() string { return event.WorktreeRemove }

// register registers this hook event decoder on c.
func register(c *hookkit.Codec) {
	c.Register(event.WorktreeRemove, hookkit.EventDecoder[Event](c))
}
