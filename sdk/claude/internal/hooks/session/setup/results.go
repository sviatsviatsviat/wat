package setup

// Results is the hook-scoped response builder supplied to Setup handlers.
type Results interface {
	// Context returns a context-injection-only Setup result.
	Context(text string) Output
	// Noop returns an empty response (silent stdout). Prefer nil from handlers when not chaining With*.
	Noop() Output
	isResults()
}

type results struct{}

func (results) isResults() {}

// Context returns a context-injection-only Setup result.
func (results) Context(text string) Output {
	return output{additionalContext: text}
}

// Noop returns an empty response (silent stdout).
func (results) Noop() Output {
	return output{}
}
