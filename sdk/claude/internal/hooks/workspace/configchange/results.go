package configchange

// Results is the hook-scoped response builder for ConfigChange.
type Results interface {
	// Block returns a block result that prevents the configuration change.
	Block(reason string) Output
	// Noop returns an empty response (silent stdout). Prefer nil from handlers when not chaining With*.
	Noop() Output
	isResults()
}

type results struct{}

func (results) isResults() {}

// Block returns a block result that prevents the configuration change.
func (results) Block(reason string) Output {
	return output{block: true, reason: reason}
}

// Noop returns an empty response (silent stdout).
func (results) Noop() Output {
	return output{}
}
