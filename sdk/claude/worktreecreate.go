package claude

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

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

// WorktreeCreate registers a WorktreeCreate handler.
func (c *Chain) WorktreeCreate(fn func(context.Context, WorktreeCreateHook) (WorktreeCreateOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev WorktreeCreate) (WorktreeCreateOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev))
	})
	return &Chain{}
}
