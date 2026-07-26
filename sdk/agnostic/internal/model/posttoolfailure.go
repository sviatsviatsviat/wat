package model

import (
	"context"
)

// PostToolFailureEvent is the normalized view of a PostToolFailure hook invocation.
type PostToolFailureEvent struct {
	Envelope
	Tool   *ToolCall
	Result *ToolResult
}

// PostToolFailureHandler handles portable PostToolFailure events.
type PostToolFailureHandler func(ctx context.Context, event PostToolFailureEvent, results PostToolFailureResults) (PostToolFailureResult, error)

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
	// On Cursor, Context is discarded because postToolUseFailure has no
	// documented output fields.
	Context(text string) PostToolFailureResult
}
