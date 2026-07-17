package model

import (
	"context"
	"encoding/json"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// PreToolEvent is the normalized view of a PreTool hook invocation.
type PreToolEvent struct {
	Envelope
	Tool *ToolCall
}

// PreToolHook is the handler context for portable PreTool events.
type PreToolHook struct {
	PreToolEvent
	inv run.Invocation
}

// NewPreToolHook wraps ev with serve-time invocation settings.
func NewPreToolHook(inv run.Invocation, ev *PreToolEvent) PreToolHook {
	h := PreToolHook{inv: inv}
	if ev != nil {
		h.PreToolEvent = *ev
	}
	return h
}

// Invocation returns serve-time settings for this hook invocation.
func (h PreToolHook) Invocation() run.Invocation { return h.inv }

// Raw returns the untouched native JSON payload.
func (h PreToolHook) Raw() json.RawMessage { return h.PreToolEvent.Raw }

// PreToolHandler handles portable PreTool events.
type PreToolHandler func(ctx context.Context, hook PreToolHook, results PreToolResults) (PreToolResult, error)

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
