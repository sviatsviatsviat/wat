package notification

import "github.com/sviatsviatsviat/wat/sdk/copilot/internal/event"

// Output is the response for notification events.
// Construct via Results builders. A nil value is a no-op.
type Output = event.ContextOutput

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
	return event.ContextResult(text)
}

// Noop returns an empty response (silent stdout).
func (results) Noop() Output {
	return event.ContextResult("")
}
