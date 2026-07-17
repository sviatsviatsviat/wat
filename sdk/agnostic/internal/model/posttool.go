package model

import (
	"context"
	"encoding/json"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// PostToolEvent is the normalized view of a PostTool hook invocation.
type PostToolEvent struct {
	Envelope
	Tool   *ToolCall
	Result *ToolResult
}

// PostToolHook is the handler context for portable PostTool events.
type PostToolHook struct {
	PostToolEvent
	inv run.Invocation
}

// NewPostToolHook wraps ev with serve-time invocation settings.
func NewPostToolHook(inv run.Invocation, ev *PostToolEvent) PostToolHook {
	h := PostToolHook{inv: inv}
	if ev != nil {
		h.PostToolEvent = *ev
	}
	return h
}

// Invocation returns serve-time settings for this hook invocation.
func (h PostToolHook) Invocation() run.Invocation { return h.inv }

// Raw returns the untouched native JSON payload.
func (h PostToolHook) Raw() json.RawMessage { return h.PostToolEvent.Raw }

// PostToolHandler handles portable PostTool events.
type PostToolHandler func(ctx context.Context, hook PostToolHook, results PostToolResults) (PostToolResult, error)

// PostToolResult is the portable hook response for PostTool events.
// Construct only via PostToolResults (Context), then With*.
// A nil value is a no-op.
type PostToolResult interface {
	// IsZero reports whether the result carries no instruction.
	IsZero() bool
	// WithUpdatedOutput replaces tool result text when set.
	// On Cursor this maps to updated_mcp_tool_output (MCP AfterMCP / post-tool only).
	WithUpdatedOutput(output string) PostToolResult
}

// PostToolResults is the hook-scoped response builder for PostTool handlers.
type PostToolResults interface {
	// Context returns a context-injection-only PostTool result.
	Context(text string) PostToolResult
}
