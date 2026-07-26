package precompact

// Results is the hook-scoped response builder supplied to handlers by registration.
type Results interface {
	// UserMessage returns a preCompact result with a user-facing message.
	// Compaction cannot be blocked; the message is observational only.
	UserMessage(text string) Output
	// Noop returns an empty response (silent stdout). Prefer nil from handlers when not chaining.
	Noop() Output
	isResults()
}

type results struct{}

func (results) isResults() {}

// UserMessage returns a preCompact result with a user-facing message.
// Compaction cannot be blocked; the message is observational only.
func (results) UserMessage(text string) Output {
	return output{userMessage: text}
}

// Noop returns an empty response (silent stdout).
func (results) Noop() Output {
	return output{}
}
