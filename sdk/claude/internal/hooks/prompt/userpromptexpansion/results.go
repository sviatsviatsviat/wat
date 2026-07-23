package userpromptexpansion

import (
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
)

// Results is the hook-scoped response builder for this event.
type Results interface {
	// Context returns a context-injection-only UserPromptExpansion result.
	Context(text string) event.CommonOutput
	isResults()
}

type results struct{}

func (results) isResults() {}

// Context returns a context-injection-only UserPromptExpansion result.
func (results) Context(text string) event.CommonOutput {
	return event.ContextOutput(event.UserPromptExpansion, text)
}
