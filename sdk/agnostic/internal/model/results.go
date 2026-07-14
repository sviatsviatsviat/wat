package model

// PreToolResult is the portable hook response for PreTool events.
type PreToolResult struct {
	// Decision is the gate verdict.
	Decision Decision
	// Reason is an agent- or model-facing explanation.
	Reason string
	// UpdatedInput replaces tool arguments when set.
	UpdatedInput map[string]any
}

// IsZero reports whether the result carries no instruction.
func (r PreToolResult) IsZero() bool {
	return r.Decision == DecisionUnset && r.Reason == "" && r.UpdatedInput == nil
}

// Result converts r into the wire Result type for codecs.
func (r PreToolResult) Result() Result {
	return Result{
		Decision:     r.Decision,
		Reason:       r.Reason,
		UpdatedInput: r.UpdatedInput,
	}
}

// PreToolAllow returns an allow verdict for PreTool events.
func PreToolAllow() PreToolResult { return PreToolResult{Decision: DecisionAllow} }

// PreToolDeny returns a deny verdict with an agent-facing reason.
func PreToolDeny(reason string) PreToolResult {
	return PreToolResult{Decision: DecisionDeny, Reason: reason}
}

// PreToolAsk returns an escalate-to-user verdict with an agent-facing reason.
func PreToolAsk(reason string) PreToolResult {
	return PreToolResult{Decision: DecisionAsk, Reason: reason}
}

// PostToolResult is the portable hook response for PostTool events.
type PostToolResult struct {
	// UpdatedOutput replaces tool result text when set.
	UpdatedOutput *string
	// Context is additional context injected for the model.
	Context string
}

// IsZero reports whether the result carries no instruction.
func (r PostToolResult) IsZero() bool {
	return r.UpdatedOutput == nil && r.Context == ""
}

// Result converts r into the wire Result type for codecs.
func (r PostToolResult) Result() Result {
	return Result{UpdatedOutput: r.UpdatedOutput, Context: r.Context}
}

// PostToolContext returns a context-injection-only PostTool result.
func PostToolContext(text string) PostToolResult { return PostToolResult{Context: text} }

// PostToolFailureResult is the portable hook response for PostToolFailure events.
type PostToolFailureResult struct {
	// Context is recovery guidance injected for the model.
	Context string
}

// IsZero reports whether the result carries no instruction.
func (r PostToolFailureResult) IsZero() bool { return r.Context == "" }

// Result converts r into the wire Result type for codecs.
func (r PostToolFailureResult) Result() Result { return Result{Context: r.Context} }

// PostToolFailureContext returns recovery guidance for PostToolFailure events.
func PostToolFailureContext(text string) PostToolFailureResult {
	return PostToolFailureResult{Context: text}
}

// StopResult is the portable hook response for Stop and SubagentStop events.
type StopResult struct {
	// FollowUp instructs the agent to keep working.
	FollowUp string
}

// IsZero reports whether the result carries no instruction.
func (r StopResult) IsZero() bool { return r.FollowUp == "" }

// Result converts r into the wire Result type for codecs.
func (r StopResult) Result() Result { return Result{FollowUp: r.FollowUp} }

// StopFollowUp returns a stop-gate result that blocks completion with a follow-up instruction.
func StopFollowUp(text string) StopResult { return StopResult{FollowUp: text} }

// SessionStartResult is the portable hook response for SessionStart events.
type SessionStartResult struct {
	// Context is additional context injected for the model.
	Context string
}

// IsZero reports whether the result carries no instruction.
func (r SessionStartResult) IsZero() bool { return r.Context == "" }

// Result converts r into the wire Result type for codecs.
func (r SessionStartResult) Result() Result { return Result{Context: r.Context} }

// SessionStartContext returns a context-injection-only SessionStart result.
func SessionStartContext(text string) SessionStartResult {
	return SessionStartResult{Context: text}
}
