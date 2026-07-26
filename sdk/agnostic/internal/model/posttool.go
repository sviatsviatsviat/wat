package model

import (
	"context"
)

// PostToolEvent is the normalized view of a PostTool hook invocation.
type PostToolEvent struct {
	Envelope
	Tool   *ToolCall
	Result *ToolResult
}

// PostToolHandler handles portable PostTool events.
type PostToolHandler func(ctx context.Context, event PostToolEvent, results PostToolResults) (PostToolResult, error)

// PostToolResult is the portable hook response for PostTool events.
// Construct only via PostToolResults (Context), then With*.
// A nil value is a no-op.
type PostToolResult interface {
	// IsZero reports whether the result carries no instruction.
	IsZero() bool
	// WithUpdatedOutput replaces tool result text when set.
	// On Cursor this maps to updated_mcp_tool_output on generic postToolUse for
	// MCP tools only; dedicated afterMCPExecution is observe-only.
	WithUpdatedOutput(output string) PostToolResult
}

// PostToolResults is the hook-scoped response builder for PostTool handlers.
type PostToolResults interface {
	// Context returns a context-injection-only PostTool result.
	Context(text string) PostToolResult
}
