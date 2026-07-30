package posttoolbatch

// Results is the hook-scoped response builder for PostToolBatch.
type Results interface {
	// Context returns a context-injection-only PostToolBatch result.
	Context(text string) Output
	// Block returns a block result with an agent-facing reason.
	Block(reason string) Output
	isResults()
}

type results struct{}

func (results) isResults() {}

// Context returns a context-injection-only PostToolBatch result.
func (results) Context(text string) Output {
	return output{additionalContext: text}
}

// Block returns a block result with an agent-facing reason.
func (results) Block(reason string) Output {
	return output{block: true, reason: reason}
}
