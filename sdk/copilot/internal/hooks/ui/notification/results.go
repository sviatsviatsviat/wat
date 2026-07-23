package notification

// Results is the hook-scoped response builder supplied to handlers by registration.
type Results interface {
	// Context returns a context-injection-only Notification result.
	Context(text string) Output
	// Noop returns an empty response (silent stdout). Prefer nil from handlers when not chaining With*.
	Noop() Output
	isResults()
}

type results struct{}

func (results) isResults() {}

// Context returns a context-injection-only Notification result.
func (results) Context(text string) Output {
	return output{additionalContext: text}
}

// Noop returns an empty response (silent stdout).
func (results) Noop() Output {
	return output{}
}
