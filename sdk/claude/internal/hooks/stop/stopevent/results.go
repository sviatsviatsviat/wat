package stopevent

// Results is the hook-scoped response builder for this event.
type Results interface {
	// Context returns non-blocking feedback that continues the conversation.
	Context(text string) Output
	// FollowUp blocks completion and feeds reason back to Claude.
	FollowUp(reason string) Output
	isResults()
}

type results struct {
	eventName string
}

func (results) isResults() {}

// Context returns non-blocking feedback that continues the conversation.
func (r results) Context(text string) Output {
	return output{eventName: r.eventName, additionalContext: text}
}

// FollowUp blocks completion and feeds reason back to Claude.
func (r results) FollowUp(reason string) Output {
	return output{eventName: r.eventName, block: true, reason: reason}
}
