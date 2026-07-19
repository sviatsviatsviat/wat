package agnostic

// Decision is the unified gate verdict for pre-events.
type Decision int

const (
	// DecisionUnset means the handler expressed no gate verdict.
	DecisionUnset Decision = iota
	// DecisionAllow permits the gated action to proceed.
	DecisionAllow
	// DecisionAsk escalates the decision to the user.
	DecisionAsk
	// DecisionDeny blocks the gated action.
	DecisionDeny
)

// String returns the agent-facing decision label ("allow", "deny", "ask")
// or an empty string for DecisionUnset.
func (d Decision) String() string {
	switch d {
	case DecisionAllow:
		return "allow"
	case DecisionDeny:
		return "deny"
	case DecisionAsk:
		return "ask"
	default:
		return ""
	}
}

// Result is a portable decision/context projection used by Merge helpers and tests.
//
// Runtime hook handlers return sealed kind-specific types (PreToolResult,
// PostToolResult, …), not this struct. Result is distinct from ToolResult,
// which carries incoming post-tool payload data on decoded events.
type Result struct {
	// Decision is the gate verdict for PreTool events.
	Decision Decision
	// Reason is an agent- or model-facing explanation.
	Reason string
	// UpdatedInput replaces tool arguments on PreTool events.
	UpdatedInput map[string]any
	// UpdatedOutput replaces tool result text on PostTool events.
	UpdatedOutput *string
	// Context is additional context injected for the
	Context string
	// FollowUp instructs the agent to keep working on Stop and SubagentStop events.
	FollowUp string
}

// Allow returns an allow verdict.
func Allow() Result { return Result{Decision: DecisionAllow} }

// Deny returns a deny verdict with an agent-facing reason.
func Deny(reason string) Result { return Result{Decision: DecisionDeny, Reason: reason} }

// Ask returns an escalate-to-user verdict with an agent-facing reason.
func Ask(reason string) Result { return Result{Decision: DecisionAsk, Reason: reason} }

// Context returns a context-injection-only result.
func Context(text string) Result { return Result{Context: text} }

// IsZero reports whether the result carries no instruction at all.
func (r Result) IsZero() bool {
	return r.Decision == DecisionUnset && r.Reason == "" &&
		r.UpdatedInput == nil && r.UpdatedOutput == nil && r.Context == "" &&
		r.FollowUp == ""
}

// Merge combines b into a. Deny outranks Ask outranks Allow; text fields are
// overridden by non-empty values from b, except Context which accumulates.
func Merge(a, b Result) Result {
	if b.Decision > a.Decision {
		a.Decision = b.Decision
	}
	if b.Reason != "" {
		a.Reason = b.Reason
	}
	if b.UpdatedInput != nil {
		a.UpdatedInput = b.UpdatedInput
	}
	if b.UpdatedOutput != nil {
		a.UpdatedOutput = b.UpdatedOutput
	}
	if b.Context != "" {
		if a.Context != "" {
			a.Context += "\n\n" + b.Context
		} else {
			a.Context = b.Context
		}
	}
	if b.FollowUp != "" {
		a.FollowUp = b.FollowUp
	}
	return a
}
