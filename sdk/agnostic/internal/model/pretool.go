package model

import (
	"context"
)

// PreToolEvent is the normalized view of a PreTool hook invocation.
type PreToolEvent struct {
	Envelope
	Tool *ToolCall
}

// PreToolHandler handles portable PreTool events.
type PreToolHandler func(ctx context.Context, event PreToolEvent, results PreToolResults) (PreToolResult, error)

// PreToolResult is the portable hook response for PreTool events.
// Construct only via PreToolResults (Allow/Deny/Ask), then With*.
// A nil value is a no-op.
type PreToolResult interface {
	// IsZero reports whether the result carries no instruction.
	IsZero() bool
	// WithUpdatedInput replaces tool arguments when set.
	// On Cursor, updated_input is emitted only for preToolUse (not beforeShellExecution,
	// beforeMCPExecution, or beforeReadFile).
	WithUpdatedInput(input map[string]any) PreToolResult
}

// PreToolResults is the hook-scoped response builder for PreTool handlers.
type PreToolResults interface {
	// Allow returns an allow verdict.
	Allow() PreToolResult
	// Deny returns a deny verdict with an agent-facing reason.
	Deny(reason string) PreToolResult
	// Ask returns an escalate-to-user verdict with an agent-facing reason.
	Ask(reason string) PreToolResult
}
