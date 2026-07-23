package posttoolusefailure

// Results is the hook-scoped response builder supplied to handlers by registration.
type Results interface {
	// Context returns recovery guidance for PostToolUseFailure events.
	Context(text string) Output
	// Noop returns an empty response (silent stdout). Prefer nil from handlers when not chaining With*.
	Noop() Output
	isResults()
}

type results struct{}

func (results) isResults() {}

// Context returns recovery guidance for PostToolUseFailure events.
func (results) Context(text string) Output {
	return output{context: text}
}

// Noop returns an empty response (silent stdout).
func (results) Noop() Output {
	return output{}
}
