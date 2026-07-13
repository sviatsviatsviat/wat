package copilothook

// PermissionDecision is a pre-tool permission verdict label.
type PermissionDecision string

const (
	// DecisionAllow permits the tool call.
	DecisionAllow PermissionDecision = "allow"
	// DecisionDeny blocks the tool call.
	DecisionDeny PermissionDecision = "deny"
	// DecisionAsk escalates to the user.
	DecisionAsk PermissionDecision = "ask"
)

// PreToolOutput is the response for preToolUse events.
type PreToolOutput struct {
	// Decision is the permission verdict (allow, deny, ask).
	Decision PermissionDecision
	// Reason is the agent-facing decision reason.
	Reason string
	// ModifiedArgs replaces tool arguments when set.
	ModifiedArgs map[string]any
}

func (o PreToolOutput) isZero() bool {
	return o.Decision == "" && o.Reason == "" && o.ModifiedArgs == nil
}

// PostToolOutput is the response for postToolUse events.
type PostToolOutput struct {
	// ModifiedResult replaces the tool result text when set.
	ModifiedResult string
	// AdditionalContext injects model context.
	AdditionalContext string
}

func (o PostToolOutput) isZero() bool {
	return o.ModifiedResult == "" && o.AdditionalContext == ""
}

// StopOutput is the response for agentStop and subagentStop events.
type StopOutput struct {
	// Reason is fed back to the agent as the next instruction.
	Reason string
}

func (o StopOutput) isZero() bool {
	return o.Reason == ""
}

// PermissionRequestOutput is the response for permissionRequest events.
type PermissionRequestOutput struct {
	// Behavior is allow or deny.
	Behavior string
	// Message is the permission message.
	Message string
	// Interrupt stops the session when true.
	Interrupt bool
	// SuppressWarnExit skips exit code 2 when Behavior is deny. Use for ask-style
	// responses that emit deny on the wire without warn-exit semantics.
	SuppressWarnExit bool
}

func (o PermissionRequestOutput) isZero() bool {
	return o.Behavior == "" && o.Message == "" && !o.Interrupt
}

// PostToolFailureOutput is the response for postToolUseFailure events.
type PostToolFailureOutput struct {
	// Context is recovery guidance written as raw stdout text.
	Context string
}

func (o PostToolFailureOutput) isZero() bool {
	return o.Context == ""
}

// SessionStartOutput is the response for sessionStart events.
type SessionStartOutput struct {
	// AdditionalContext injects model context.
	AdditionalContext string
}

func (o SessionStartOutput) isZero() bool {
	return o.AdditionalContext == ""
}

// SubagentStartOutput is the response for subagentStart events.
type SubagentStartOutput struct {
	// AdditionalContext injects model context.
	AdditionalContext string
}

func (o SubagentStartOutput) isZero() bool {
	return o.AdditionalContext == ""
}

// NotificationOutput is the response for notification events.
type NotificationOutput struct {
	// AdditionalContext injects model context.
	AdditionalContext string
}

func (o NotificationOutput) isZero() bool {
	return o.AdditionalContext == ""
}
