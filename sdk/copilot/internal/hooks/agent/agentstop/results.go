package agentstop

// Results is the hook-scoped response builder supplied to On* handlers by registration.
type Results interface {
	// FollowUp blocks completion and feeds reason back to the agent.
	FollowUp(reason string) Output
	// Noop returns an empty response (silent stdout). Prefer nil from handlers when not chaining With*.
	Noop() Output
	isResults()
}

type results struct{}

func (results) isResults() {}

// FollowUp blocks completion and feeds reason back to the agent.
func (results) FollowUp(reason string) Output {
	return output{reason: reason}
}

// Noop returns an empty response (silent stdout).
func (results) Noop() Output {
	return output{}
}

// NewResults returns a Results builder for handlers outside this package (e.g. SubagentStop).
func NewResults() Results {
	return results{}
}
