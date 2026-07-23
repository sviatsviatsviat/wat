package model

import (
	"context"
)

// SessionEndEvent is the normalized view of a SessionEnd hook invocation.
type SessionEndEvent struct {
	Envelope
	Life *Lifecycle
}

// SessionEndHandler handles observe-only SessionEnd events.
type SessionEndHandler func(ctx context.Context, event SessionEndEvent) error
