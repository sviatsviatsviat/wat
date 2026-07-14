package agnostic

import "github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"

// PreToolResults constructs portable PreTool hook responses. Handlers receive
// this interface as the third parameter of PreToolHandler.
type PreToolResults interface {
	// Allow returns an allow verdict.
	Allow() PreToolResult
	// Deny returns a deny verdict with an agent-facing reason.
	Deny(reason string) PreToolResult
	// Ask returns an escalate-to-user verdict with an agent-facing reason.
	Ask(reason string) PreToolResult
}

type preToolResults struct{}

// Allow returns an allow verdict.
func (preToolResults) Allow() PreToolResult { return model.PreToolAllow() }

// Deny returns a deny verdict with an agent-facing reason.
func (preToolResults) Deny(reason string) PreToolResult { return model.PreToolDeny(reason) }

// Ask returns an escalate-to-user verdict with an agent-facing reason.
func (preToolResults) Ask(reason string) PreToolResult { return model.PreToolAsk(reason) }

// PostToolResults constructs portable PostTool hook responses. Handlers receive
// this interface as the third parameter of PostToolHandler.
type PostToolResults interface {
	// Context returns a context-injection-only PostTool result.
	Context(text string) PostToolResult
}

type postToolResults struct{}

// Context returns a context-injection-only PostTool result.
func (postToolResults) Context(text string) PostToolResult { return model.PostToolContext(text) }

// PostToolFailureResults constructs portable PostToolFailure hook responses.
// Handlers receive this interface as the third parameter of PostToolFailureHandler.
type PostToolFailureResults interface {
	// Context returns recovery guidance for PostToolFailure events.
	Context(text string) PostToolFailureResult
}

type postToolFailureResults struct{}

// Context returns recovery guidance for PostToolFailure events.
func (postToolFailureResults) Context(text string) PostToolFailureResult {
	return model.PostToolFailureContext(text)
}

// StopResults constructs portable Stop and SubagentStop hook responses. Handlers
// receive this interface as the third parameter of StopHandler.
type StopResults interface {
	// FollowUp returns a stop-gate result that blocks completion with a follow-up instruction.
	FollowUp(text string) StopResult
}

type stopResults struct{}

// FollowUp returns a stop-gate result that blocks completion with a follow-up instruction.
func (stopResults) FollowUp(text string) StopResult { return model.StopFollowUp(text) }

// SessionStartResults constructs portable SessionStart hook responses. Handlers
// receive this interface as the third parameter of SessionStartHandler.
type SessionStartResults interface {
	// Context returns a context-injection-only SessionStart result.
	Context(text string) SessionStartResult
}

type sessionStartResults struct{}

// Context returns a context-injection-only SessionStart result.
func (sessionStartResults) Context(text string) SessionStartResult {
	return model.SessionStartContext(text)
}
