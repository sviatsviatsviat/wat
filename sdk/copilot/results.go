package copilot

// PreToolResults constructs preToolUse hook responses. Handlers receive this
// interface as the third parameter of Chain.PreToolUse handlers.
type PreToolResults interface {
	// Allow returns an allow verdict.
	Allow() PreToolOutput
	// Deny returns a deny verdict with an agent-facing reason.
	Deny(reason string) PreToolOutput
	// Ask returns an ask verdict with an agent-facing reason.
	Ask(reason string) PreToolOutput
}

type preToolResults struct{}

// Allow returns an allow verdict.
func (preToolResults) Allow() PreToolOutput {
	return PreToolOutput{Decision: DecisionAllow}
}

// Deny returns a deny verdict with an agent-facing reason.
func (preToolResults) Deny(reason string) PreToolOutput {
	return PreToolOutput{Decision: DecisionDeny, Reason: reason}
}

// Ask returns an ask verdict with an agent-facing reason.
func (preToolResults) Ask(reason string) PreToolOutput {
	return PreToolOutput{Decision: DecisionAsk, Reason: reason}
}

// PostToolResults constructs postToolUse hook responses.
type PostToolResults interface {
	// Context returns a context-injection-only PostTool result.
	Context(text string) PostToolOutput
}

type postToolResults struct{}

// Context returns a context-injection-only PostTool result.
func (postToolResults) Context(text string) PostToolOutput {
	return PostToolOutput{AdditionalContext: text}
}

// PostToolFailureResults constructs postToolUseFailure hook responses.
type PostToolFailureResults interface {
	// Context returns recovery guidance for postToolUseFailure events.
	Context(text string) PostToolFailureOutput
}

type postToolFailureResults struct{}

// Context returns recovery guidance for postToolUseFailure events.
func (postToolFailureResults) Context(text string) PostToolFailureOutput {
	return PostToolFailureOutput{Context: text}
}

// StopResults constructs agentStop and subagentStop hook responses.
type StopResults interface {
	// FollowUp blocks completion and feeds reason back to the agent.
	FollowUp(reason string) StopOutput
}

type stopResults struct{}

// FollowUp blocks completion and feeds reason back to the agent.
func (stopResults) FollowUp(reason string) StopOutput {
	return StopOutput{Reason: reason}
}

// PermissionRequestResults constructs permissionRequest hook responses.
type PermissionRequestResults interface {
	// Allow returns an allow verdict.
	Allow() PermissionRequestOutput
	// Deny returns a deny verdict with a permission message.
	Deny(message string) PermissionRequestOutput
	// Ask returns an ask-style deny that suppresses warn-exit semantics.
	Ask(message string) PermissionRequestOutput
}

type permissionRequestResults struct{}

// Allow returns an allow verdict.
func (permissionRequestResults) Allow() PermissionRequestOutput {
	return PermissionRequestOutput{Behavior: "allow"}
}

// Deny returns a deny verdict with a permission message.
func (permissionRequestResults) Deny(message string) PermissionRequestOutput {
	return PermissionRequestOutput{Behavior: "deny", Message: message}
}

// Ask returns an ask-style deny that suppresses warn-exit semantics.
func (permissionRequestResults) Ask(message string) PermissionRequestOutput {
	return PermissionRequestOutput{Behavior: "deny", Message: message, SuppressWarnExit: true}
}

// SessionStartResults constructs sessionStart hook responses.
type SessionStartResults interface {
	// Context returns a context-injection-only SessionStart result.
	Context(text string) SessionStartOutput
}

type sessionStartResults struct{}

// Context returns a context-injection-only SessionStart result.
func (sessionStartResults) Context(text string) SessionStartOutput {
	return SessionStartOutput{AdditionalContext: text}
}

// SubagentStartResults constructs subagentStart hook responses.
type SubagentStartResults interface {
	// Context returns a context-injection-only SubagentStart result.
	Context(text string) SubagentStartOutput
}

type subagentStartResults struct{}

// Context returns a context-injection-only SubagentStart result.
func (subagentStartResults) Context(text string) SubagentStartOutput {
	return SubagentStartOutput{AdditionalContext: text}
}

// NotificationResults constructs notification hook responses.
type NotificationResults interface {
	// Context returns a context-injection-only Notification result.
	Context(text string) NotificationOutput
}

type notificationResults struct{}

// Context returns a context-injection-only Notification result.
func (notificationResults) Context(text string) NotificationOutput {
	return NotificationOutput{AdditionalContext: text}
}
