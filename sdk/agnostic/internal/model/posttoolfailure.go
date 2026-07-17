package model

import (
	"context"
	"encoding/json"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// PostToolFailureEvent is the normalized view of a PostToolFailure hook invocation.
type PostToolFailureEvent struct {
	Envelope
	Tool   *ToolCall
	Result *ToolResult
}

// PostToolFailureHook is the handler context for portable PostToolFailure events.
type PostToolFailureHook struct {
	PostToolFailureEvent
	inv run.Invocation
}

// NewPostToolFailureHook wraps ev with serve-time invocation settings.
func NewPostToolFailureHook(inv run.Invocation, ev *PostToolFailureEvent) PostToolFailureHook {
	h := PostToolFailureHook{inv: inv}
	if ev != nil {
		h.PostToolFailureEvent = *ev
	}
	return h
}

// Invocation returns serve-time settings for this hook invocation.
func (h PostToolFailureHook) Invocation() run.Invocation { return h.inv }

// Raw returns the untouched native JSON payload.
func (h PostToolFailureHook) Raw() json.RawMessage { return h.PostToolFailureEvent.Raw }

// PostToolFailureHandler handles portable PostToolFailure events.
type PostToolFailureHandler func(ctx context.Context, hook PostToolFailureHook, results PostToolFailureResults) (PostToolFailureResult, error)

// PostToolFailureResult is the portable hook response for PostToolFailure events.
// Construct only via PostToolFailureResults (Context).
// A nil value is a no-op.
type PostToolFailureResult interface {
	// IsZero reports whether the result carries no instruction.
	IsZero() bool
}

// PostToolFailureResults is the hook-scoped response builder for PostToolFailure handlers.
type PostToolFailureResults interface {
	// Context returns recovery guidance for PostToolFailure events.
	Context(text string) PostToolFailureResult
}
