package model

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// SessionStartEvent is the normalized view of a SessionStart hook invocation.
type SessionStartEvent struct {
	Envelope
	Life *Lifecycle
}

// SessionStartHook is the handler context for portable SessionStart events.
type SessionStartHook struct {
	SessionStartEvent
	inv run.Invocation
}

// NewSessionStartHook wraps ev with serve-time invocation settings.
func NewSessionStartHook(inv run.Invocation, ev *SessionStartEvent) SessionStartHook {
	h := SessionStartHook{inv: inv}
	if ev != nil {
		h.SessionStartEvent = *ev
	}
	return h
}

// Invocation returns serve-time settings for this hook invocation.
func (h SessionStartHook) Invocation() run.Invocation { return h.inv }

// SessionStartHandler handles portable SessionStart events.
type SessionStartHandler func(ctx context.Context, hook SessionStartHook, results SessionStartResults) (SessionStartResult, error)

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
