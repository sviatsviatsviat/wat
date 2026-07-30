package subagentstop

// Results is the hook-scoped response builder supplied to SubagentStop handlers
// by registration.
type Results interface {
	// FollowUp blocks completion and feeds reason back to the subagent.
	// A valid block decision wins over ModifiedResponse on the same output.
	FollowUp(reason string) Output
	// ModifiedResponse replaces the response returned to the parent when the
	// subagent is allowed to complete. Not applicable to AgentStop.
	ModifiedResponse(text string) Output
	// Noop returns an empty response (silent stdout). Prefer nil from handlers
	// when not chaining With*.
	Noop() Output
	isResults()
}

type results struct{}

func (results) isResults() {}

// FollowUp blocks completion and feeds reason back to the subagent.
func (results) FollowUp(reason string) Output {
	return output{reason: reason}
}

// ModifiedResponse replaces the response returned to the parent when the
// subagent is allowed to complete.
func (results) ModifiedResponse(text string) Output {
	return output{modifiedResponse: text}
}

// Noop returns an empty response (silent stdout).
func (results) Noop() Output {
	return output{}
}

// NewResults returns a Results builder for SubagentStop handlers.
func NewResults() Results {
	return results{}
}
