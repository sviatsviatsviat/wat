package model

import (
	"context"
)

// SessionStartEvent is the normalized view of a SessionStart hook invocation.
type SessionStartEvent struct {
	Envelope
	Life *Lifecycle
}

// SessionStartHandler handles portable SessionStart events.
type SessionStartHandler func(ctx context.Context, event SessionStartEvent, results SessionStartResults) (SessionStartResult, error)

// SessionStartResult is the portable hook response for SessionStart events.
// Construct only via SessionStartResults (Context).
// A nil value is a no-op.
type SessionStartResult interface {
	// IsZero reports whether the result carries no instruction.
	IsZero() bool
}

// SessionStartResults is the hook-scoped response builder for SessionStart handlers.
type SessionStartResults interface {
	// Context returns a context-injection-only SessionStart result.
	Context(text string) SessionStartResult
}
