package claude

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

// CommonOutput is a Common-only response for observe-only events.
type CommonOutput struct {
	Common
	// AdditionalContext injects model context.
	AdditionalContext string
}

func (o CommonOutput) isZero() bool {
	return o.Common.isZero() && o.AdditionalContext == ""
}
