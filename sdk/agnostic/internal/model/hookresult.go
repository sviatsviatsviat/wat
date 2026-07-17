package model

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

// StopResult is the portable hook response for Stop and SubagentStop events.
// Construct only via StopResults (FollowUp).
// A nil value is a no-op.
type StopResult interface {
	// IsZero reports whether the result carries no instruction.
	IsZero() bool
}

// StopResults is the hook-scoped response builder for Stop handlers.
type StopResults interface {
	// FollowUp returns a stop-gate result that blocks completion with a follow-up instruction.
	FollowUp(text string) StopResult
}

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
