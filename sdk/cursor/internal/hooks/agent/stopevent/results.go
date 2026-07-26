package stopevent

// Results is the hook-scoped response builder supplied to On* handlers by registration.
type Results interface {
	// FollowUp emits followup_message so Cursor auto-submits text as the next
	// user message. Auto follow-ups are limited by hooks.json loop_limit
	// (default 5; null means unlimited); see Event.LoopCount.
	FollowUp(text string) Output
	// Noop returns an empty response (silent stdout). Prefer nil from handlers when not chaining.
	Noop() Output
	isResults()
}

type results struct{}

func (results) isResults() {}

// FollowUp emits followup_message so Cursor auto-submits text as the next user
// message. Auto follow-ups are limited by hooks.json loop_limit (default 5;
// null means unlimited); see Event.LoopCount.
func (results) FollowUp(text string) Output {
	return output{followUpMessage: text}
}

// Noop returns an empty response (silent stdout).
func (results) Noop() Output {
	return output{}
}

// NewResults returns a Results builder for handlers outside this package (e.g. SubagentStop).
func NewResults() Results {
	return results{}
}
