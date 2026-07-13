package claudehook

// PermissionDecision is a pre-tool permission verdict label.
type PermissionDecision string

const (
	// DecisionAllow permits the tool call.
	DecisionAllow PermissionDecision = "allow"
	// DecisionDeny blocks the tool call.
	DecisionDeny PermissionDecision = "deny"
	// DecisionAsk escalates to the user.
	DecisionAsk PermissionDecision = "ask"
	// DecisionDefer defers the permission decision.
	DecisionDefer PermissionDecision = "defer"
)

// Common holds output fields shared across Claude Code hook responses.
type Common struct {
	// Continue when false stops Claude entirely.
	Continue *bool
	// StopReason explains why the session was stopped.
	StopReason string
	// SuppressOutput suppresses hook output when true.
	SuppressOutput bool
	// SystemMessage is a user-visible system message.
	SystemMessage string
	// TerminalSequence is an OSC terminal sequence (allowlisted).
	TerminalSequence string
}

func (c Common) isZero() bool {
	return c.Continue == nil && c.StopReason == "" && !c.SuppressOutput &&
		c.SystemMessage == "" && c.TerminalSequence == ""
}

// PreToolUseOutput is the response for PreToolUse events.
type PreToolUseOutput struct {
	Common
	// Decision is the permission verdict (allow, deny, ask, defer).
	Decision PermissionDecision
	// Reason is the agent-facing decision reason.
	Reason string
	// UpdatedInput replaces tool arguments when set.
	UpdatedInput map[string]any
	// AdditionalContext injects model context.
	AdditionalContext string
}

func (o PreToolUseOutput) isZero() bool {
	return o.Common.isZero() && o.Decision == "" && o.Reason == "" &&
		o.UpdatedInput == nil && o.AdditionalContext == ""
}

// PermissionRequestOutput is the response for PermissionRequest events.
type PermissionRequestOutput struct {
	Common
	// Behavior is allow or deny.
	Behavior string
	// UpdatedInput replaces tool arguments when set.
	UpdatedInput map[string]any
	// Message is the permission message.
	Message string
	// Interrupt stops the session when true.
	Interrupt bool
	// AdditionalContext injects model context.
	AdditionalContext string
}

func (o PermissionRequestOutput) isZero() bool {
	return o.Common.isZero() && o.Behavior == "" && o.UpdatedInput == nil &&
		o.Message == "" && !o.Interrupt && o.AdditionalContext == ""
}

// PostToolUseOutput is the response for PostToolUse and PostToolUseFailure events.
type PostToolUseOutput struct {
	Common
	// Block prompts Claude with Reason when true.
	Block bool
	// Reason is the block reason.
	Reason string
	// AdditionalContext injects model context.
	AdditionalContext string
	// UpdatedToolOutput replaces the tool result when set.
	UpdatedToolOutput any
}

func (o PostToolUseOutput) isZero() bool {
	return o.Common.isZero() && !o.Block && o.Reason == "" &&
		o.AdditionalContext == "" && o.UpdatedToolOutput == nil
}

// UserPromptSubmitOutput is the response for UserPromptSubmit events.
type UserPromptSubmitOutput struct {
	Common
	// Block rejects the submitted prompt when true.
	Block bool
	// Reason is the block reason.
	Reason string
	// AdditionalContext injects model context.
	AdditionalContext string
	// SessionTitle sets the session title.
	SessionTitle string
	// SuppressOriginalPrompt suppresses the original prompt when true.
	SuppressOriginalPrompt bool
}

func (o UserPromptSubmitOutput) isZero() bool {
	return o.Common.isZero() && !o.Block && o.Reason == "" &&
		o.AdditionalContext == "" && o.SessionTitle == "" && !o.SuppressOriginalPrompt
}

// StopOutput is the response for Stop and SubagentStop events.
type StopOutput struct {
	Common
	// Block keeps the agent working when true.
	Block bool
	// Reason is fed back to Claude as the next instruction.
	Reason string
	// AdditionalContext is non-error feedback that continues the conversation.
	AdditionalContext string
}

func (o StopOutput) isZero() bool {
	return o.Common.isZero() && !o.Block && o.Reason == "" && o.AdditionalContext == ""
}

// SessionStartOutput is the response for SessionStart events.
type SessionStartOutput struct {
	Common
	// AdditionalContext injects model context.
	AdditionalContext string
	// InitialUserMessage sets the initial user message.
	InitialUserMessage string
	// SessionTitle sets the session title.
	SessionTitle string
	// WatchPaths registers filesystem watch paths.
	WatchPaths []string
	// ReloadSkills reloads skills when true.
	ReloadSkills bool
	// Env carries session environment variables written to CLAUDE_ENV_FILE.
	Env map[string]string
}

func (o SessionStartOutput) isZero() bool {
	return o.Common.isZero() && o.AdditionalContext == "" && o.InitialUserMessage == "" &&
		o.SessionTitle == "" && len(o.WatchPaths) == 0 && !o.ReloadSkills && len(o.Env) == 0
}

// MessageDisplayOutput is the response for MessageDisplay events.
type MessageDisplayOutput struct {
	Common
	// DisplayContent overrides displayed content when set.
	DisplayContent *string
}

func (o MessageDisplayOutput) isZero() bool {
	return o.Common.isZero() && o.DisplayContent == nil
}

// PermissionDeniedOutput is the response for PermissionDenied events.
type PermissionDeniedOutput struct {
	Common
	// Retry requests a permission retry when true.
	Retry bool
}

func (o PermissionDeniedOutput) isZero() bool {
	return o.Common.isZero() && !o.Retry
}

// ElicitationOutput is the response for Elicitation events.
type ElicitationOutput struct {
	Common
	// Action is the elicitation action.
	Action string
	// Content is the elicitation response content.
	Content map[string]any
}

func (o ElicitationOutput) isZero() bool {
	return o.Common.isZero() && o.Action == "" && o.Content == nil
}

// WorktreeCreateOutput is the response for WorktreeCreate events.
type WorktreeCreateOutput struct {
	Common
	// WorktreePath is the created worktree path.
	WorktreePath string
}

func (o WorktreeCreateOutput) isZero() bool {
	return o.Common.isZero() && o.WorktreePath == ""
}

// CommonOutput is a Common-only response for observe-only events.
type CommonOutput struct {
	Common
	// AdditionalContext injects model context.
	AdditionalContext string
}

func (o CommonOutput) isZero() bool {
	return o.Common.isZero() && o.AdditionalContext == ""
}
