package model

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// SessionEndEvent is the normalized view of a SessionEnd hook invocation.
type SessionEndEvent struct {
	Envelope
	Life *Lifecycle
}

// SessionEndHook is the handler context for portable SessionEnd events.
type SessionEndHook struct {
	SessionEndEvent
	inv run.Invocation
}

// NewSessionEndHook wraps ev with serve-time invocation settings.
func NewSessionEndHook(inv run.Invocation, ev *SessionEndEvent) SessionEndHook {
	h := SessionEndHook{inv: inv}
	if ev != nil {
		h.SessionEndEvent = *ev
	}
	return h
}

// Invocation returns serve-time settings for this hook invocation.
func (h SessionEndHook) Invocation() run.Invocation { return h.inv }

// SessionEndHandler handles observe-only SessionEnd events.
type SessionEndHandler func(ctx context.Context, hook SessionEndHook) error
