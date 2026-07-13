package cursor

// PermissionDecision is a permission verdict label on permission-gating events.
type PermissionDecision string

const (
	// DecisionAllow permits the action.
	DecisionAllow PermissionDecision = "allow"
	// DecisionDeny blocks the action.
	DecisionDeny PermissionDecision = "deny"
	// DecisionAsk escalates to the user.
	DecisionAsk PermissionDecision = "ask"
)

// PermissionOutput is the response for permission-gating events.
type PermissionOutput struct {
	// Decision is the permission verdict (allow, deny, ask).
	Decision PermissionDecision
	// UserMessage is shown to the user.
	UserMessage string
	// AgentMessage is shown to the agent.
	AgentMessage string
	// UpdatedInput replaces tool arguments on preToolUse when set.
	UpdatedInput map[string]any
}

func (o PermissionOutput) isZero() bool {
	return o.Decision == "" && o.UserMessage == "" && o.AgentMessage == "" && o.UpdatedInput == nil
}

// BeforeSubmitPromptOutput is the response for beforeSubmitPrompt events.
type BeforeSubmitPromptOutput struct {
	// Continue is false to block prompt submission.
	Continue *bool
	// UserMessage is shown to the user when blocking.
	UserMessage string
}

func (o BeforeSubmitPromptOutput) isZero() bool {
	return o.Continue == nil && o.UserMessage == ""
}

// PostToolOutput is the response for post-tool events.
type PostToolOutput struct {
	// UpdatedMCPOutput replaces MCP tool output when set.
	UpdatedMCPOutput any
	// AdditionalContext injects model context.
	AdditionalContext string
}

func (o PostToolOutput) isZero() bool {
	return o.UpdatedMCPOutput == nil && o.AdditionalContext == ""
}

// StopOutput is the response for stop and subagentStop events.
type StopOutput struct {
	// FollowUpMessage is sent back to the agent as the next instruction.
	FollowUpMessage string
}

func (o StopOutput) isZero() bool {
	return o.FollowUpMessage == ""
}

// SessionStartOutput is the response for sessionStart events.
type SessionStartOutput struct {
	// Env sets environment variables for the session.
	Env map[string]string
	// AdditionalContext injects model context.
	AdditionalContext string
}

func (o SessionStartOutput) isZero() bool {
	return len(o.Env) == 0 && o.AdditionalContext == ""
}

// PreCompactOutput is the response for preCompact events.
type PreCompactOutput struct {
	// UserMessage is shown to the user.
	UserMessage string
}

func (o PreCompactOutput) isZero() bool {
	return o.UserMessage == ""
}
