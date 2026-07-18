package claude

import "github.com/sviatsviatsviat/wat/internal/hookkit"

// WorktreeRemove is the WorktreeRemove hook event.
type WorktreeRemove struct {
	Envelope
	hookkit.RawPayload
	// WorktreePath is the worktree path being removed.
	WorktreePath string `json:"worktree_path"`
}

// EventName returns the hook event name.
func (WorktreeRemove) EventName() string { return EventWorktreeRemove }

func init() {
	registerDecoder(EventWorktreeRemove, decodeAs[WorktreeRemove])
}
