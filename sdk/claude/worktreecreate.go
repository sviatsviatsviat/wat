package claude

// WorktreeCreate is the WorktreeCreate hook event.
type WorktreeCreate struct {
	Envelope
}

// EventName returns the hook event name.
func (WorktreeCreate) EventName() string { return EventWorktreeCreate }

func init() {
	registerDecoder(EventWorktreeCreate, decodeAs[WorktreeCreate])
}

// WorktreeCreateOutput is the response for WorktreeCreate events.
type WorktreeCreateOutput struct {
	Common
	// WorktreePath is the created worktree path.
	WorktreePath string
}

func (o WorktreeCreateOutput) isZero() bool {
	return o.Common.isZero() && o.WorktreePath == ""
}
