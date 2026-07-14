package claude

// PreToolUseResults constructs PreToolUse hook responses. Handlers receive
// this interface as the third parameter of Chain.PreToolUse handlers.
type PreToolUseResults interface {
	// Allow returns an allow verdict.
	Allow() PreToolUseOutput
	// Deny returns a deny verdict with an agent-facing reason.
	Deny(reason string) PreToolUseOutput
	// Ask returns an ask verdict with an agent-facing reason.
	Ask(reason string) PreToolUseOutput
	// Defer returns a defer verdict.
	Defer() PreToolUseOutput
}

type preToolUseResults struct{}

// Allow returns an allow verdict.
func (preToolUseResults) Allow() PreToolUseOutput {
	return PreToolUseOutput{Decision: DecisionAllow}
}

// Deny returns a deny verdict with an agent-facing reason.
func (preToolUseResults) Deny(reason string) PreToolUseOutput {
	return PreToolUseOutput{Decision: DecisionDeny, Reason: reason}
}

// Ask returns an ask verdict with an agent-facing reason.
func (preToolUseResults) Ask(reason string) PreToolUseOutput {
	return PreToolUseOutput{Decision: DecisionAsk, Reason: reason}
}

// Defer returns a defer verdict.
func (preToolUseResults) Defer() PreToolUseOutput {
	return PreToolUseOutput{Decision: DecisionDefer}
}

// PostToolUseResults constructs PostToolUse and PostToolUseFailure hook responses.
type PostToolUseResults interface {
	// Context returns a context-injection-only PostToolUse result.
	Context(text string) PostToolUseOutput
	// Block returns a block result with an agent-facing reason.
	Block(reason string) PostToolUseOutput
}

type postToolUseResults struct{}

// Context returns a context-injection-only PostToolUse result.
func (postToolUseResults) Context(text string) PostToolUseOutput {
	return PostToolUseOutput{AdditionalContext: text}
}

// Block returns a block result with an agent-facing reason.
func (postToolUseResults) Block(reason string) PostToolUseOutput {
	return PostToolUseOutput{Block: true, Reason: reason}
}

// PermissionRequestResults constructs PermissionRequest hook responses.
type PermissionRequestResults interface {
	// Allow returns an allow verdict.
	Allow() PermissionRequestOutput
	// Deny returns a deny verdict with a permission message.
	Deny(message string) PermissionRequestOutput
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

// UserPromptSubmitResults constructs UserPromptSubmit hook responses.
type UserPromptSubmitResults interface {
	// Block returns a block result with an agent-facing reason.
	Block(reason string) UserPromptSubmitOutput
}

type userPromptSubmitResults struct{}

// Block returns a block result with an agent-facing reason.
func (userPromptSubmitResults) Block(reason string) UserPromptSubmitOutput {
	return UserPromptSubmitOutput{Block: true, Reason: reason}
}

// StopResults constructs Stop and SubagentStop hook responses.
type StopResults interface {
	// FollowUp blocks completion and feeds reason back to Claude.
	FollowUp(reason string) StopOutput
}

type stopResults struct{}

// FollowUp blocks completion and feeds reason back to Claude.
func (stopResults) FollowUp(reason string) StopOutput {
	return StopOutput{Block: true, Reason: reason}
}

// SessionStartResults constructs SessionStart hook responses.
type SessionStartResults interface {
	// Context returns a context-injection-only SessionStart result.
	Context(text string) SessionStartOutput
}

type sessionStartResults struct{}

// Context returns a context-injection-only SessionStart result.
func (sessionStartResults) Context(text string) SessionStartOutput {
	return SessionStartOutput{AdditionalContext: text}
}

// CommonResults constructs CommonOutput hook responses for observe-only events.
type CommonResults interface {
	// Context returns a context-injection-only CommonOutput result.
	Context(text string) CommonOutput
}

type commonResults struct{}

// Context returns a context-injection-only CommonOutput result.
func (commonResults) Context(text string) CommonOutput {
	return CommonOutput{AdditionalContext: text}
}
