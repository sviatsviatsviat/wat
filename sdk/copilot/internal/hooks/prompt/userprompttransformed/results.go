package userprompttransformed

// Results is the hook-scoped response builder supplied to handlers by registration.
type Results interface {
	// Modified returns a result that replaces the model-facing transformed prompt.
	Modified(text string) Output
	// Noop returns an empty response (silent stdout). Prefer nil from handlers when not chaining With*.
	Noop() Output
	isResults()
}

type results struct{}

func (results) isResults() {}

// Modified returns a result that replaces the model-facing transformed prompt.
func (results) Modified(text string) Output {
	return output{modifiedTransformedPrompt: text}
}

// Noop returns an empty response (silent stdout).
func (results) Noop() Output {
	return output{}
}
