package model

// PreToolResult is the portable hook response for PreTool events.
// Construct via PreToolAllow/PreToolDeny/PreToolAsk or agnostic.PreToolResults, then With*.
// A nil value is a no-op.
type PreToolResult interface {
	isPreToolResult()
	// IsZero reports whether the result carries no instruction.
	IsZero() bool
	// Result converts into the wire Result type for codecs.
	Result() Result
	// WithUpdatedInput replaces tool arguments when set.
	WithUpdatedInput(input map[string]any) PreToolResult
}

type preToolResult struct {
	decision     Decision
	reason       string
	updatedInput map[string]any
}

func (preToolResult) isPreToolResult() {}

// IsZero reports whether the result carries no instruction.
func (r preToolResult) IsZero() bool {
	return r.decision == DecisionUnset && r.reason == "" && r.updatedInput == nil
}

// Result converts r into the wire Result type for codecs.
func (r preToolResult) Result() Result {
	return Result{
		Decision:     r.decision,
		Reason:       r.reason,
		UpdatedInput: r.updatedInput,
	}
}

// WithUpdatedInput replaces tool arguments when set.
func (r preToolResult) WithUpdatedInput(input map[string]any) PreToolResult {
	r.updatedInput = input
	return r
}

// PreToolAllow returns an allow verdict for PreTool events.
func PreToolAllow() PreToolResult { return preToolResult{decision: DecisionAllow} }

// PreToolDeny returns a deny verdict with an agent-facing reason.
func PreToolDeny(reason string) PreToolResult {
	return preToolResult{decision: DecisionDeny, reason: reason}
}

// PreToolAsk returns an escalate-to-user verdict with an agent-facing reason.
func PreToolAsk(reason string) PreToolResult {
	return preToolResult{decision: DecisionAsk, reason: reason}
}
