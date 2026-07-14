package cursor

// PermissionResults constructs permission-gating hook responses. Handlers
// receive this interface as the third parameter of permission Chain methods.
type PermissionResults interface {
	// Allow returns an allow verdict.
	Allow() PermissionOutput
	// Deny returns a deny verdict with an agent-facing message.
	Deny(agentMessage string) PermissionOutput
}

type permissionResults struct{}

// Allow returns an allow verdict.
func (permissionResults) Allow() PermissionOutput {
	return PermissionOutput{Decision: DecisionAllow}
}

// Deny returns a deny verdict with an agent-facing message.
func (permissionResults) Deny(agentMessage string) PermissionOutput {
	return PermissionOutput{Decision: DecisionDeny, AgentMessage: agentMessage}
}

// PostToolResults constructs post-tool hook responses.
type PostToolResults interface {
	// Context returns a context-injection-only PostTool result.
	Context(text string) PostToolOutput
}

type postToolResults struct{}

// Context returns a context-injection-only PostTool result.
func (postToolResults) Context(text string) PostToolOutput {
	return PostToolOutput{AdditionalContext: text}
}

// BeforeSubmitPromptResults constructs beforeSubmitPrompt hook responses.
type BeforeSubmitPromptResults interface {
	// Block blocks prompt submission with a user-facing message.
	Block(userMessage string) BeforeSubmitPromptOutput
}

type beforeSubmitPromptResults struct{}

// Block blocks prompt submission with a user-facing message.
func (beforeSubmitPromptResults) Block(userMessage string) BeforeSubmitPromptOutput {
	cont := false
	return BeforeSubmitPromptOutput{Continue: &cont, UserMessage: userMessage}
}

// StopResults constructs stop and subagentStop hook responses.
type StopResults interface {
	// FollowUp blocks completion and feeds a follow-up instruction to the agent.
	FollowUp(text string) StopOutput
}

type stopResults struct{}

// FollowUp blocks completion and feeds a follow-up instruction to the agent.
func (stopResults) FollowUp(text string) StopOutput {
	return StopOutput{FollowUpMessage: text}
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
