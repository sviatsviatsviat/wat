package stopevent

// Results is the hook-scoped response builder supplied to On* handlers by registration.
type Results interface {
	// FollowUp blocks completion and feeds a follow-up instruction to the agent
	// as followup_message. For subagentStop, Cursor only consumes the message
	// when the input status is "completed". For stop, a non-empty message is
	// always eligible. Host-side auto follow-up caps use the hooks.json
	// loop_limit handler option (default 5); that option is not an SDK field.
	FollowUp(text string) Output
	// Noop returns an empty response (silent stdout). Prefer nil from handlers when not chaining.
	Noop() Output
	isResults()
}

type results struct{}

func (results) isResults() {}

// FollowUp blocks completion and feeds a follow-up instruction to the agent
// as followup_message. For subagentStop, Cursor only consumes the message when
// the input status is "completed". For stop, a non-empty message is always
// eligible. Host-side auto follow-up caps use the hooks.json loop_limit
// handler option (default 5); that option is not an SDK field.
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
