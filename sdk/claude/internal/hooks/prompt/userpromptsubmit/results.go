package userpromptsubmit

// Results is the hook-scoped response builder for this event.
type Results interface {
	// Context returns a non-blocking context-injection result.
	Context(text string) Output
	// Block returns a block result with an agent-facing reason.
	Block(reason string) Output
	// Noop returns an empty response (silent stdout). Prefer nil from handlers when not chaining With*.
	Noop() Output
	isResults()
}

type results struct{}

func (results) isResults() {}

// Context returns a non-blocking context-injection result.
func (results) Context(text string) Output {
	return output{additionalContext: text}
}

// Block returns a block result with an agent-facing reason.
func (results) Block(reason string) Output {
	return output{block: true, reason: reason}
}

// Noop returns an empty response (silent stdout).
func (results) Noop() Output {
	return output{}
}
