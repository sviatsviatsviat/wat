package userpromptexpansion

import (
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
)

// Results is the hook-scoped response builder for this event.
type Results interface {
	// Context returns a context-injection-only UserPromptExpansion result.
	Context(text string) event.DecisionOutput
	// Block returns a decision:"block" result with an agent-facing reason.
	// Encodes with SuccessExit; Claude processes JSON only on exit 0.
	Block(reason string) event.DecisionOutput
	isResults()
}

type results struct{}

func (results) isResults() {}

// Context returns a context-injection-only UserPromptExpansion result.
func (results) Context(text string) event.DecisionOutput {
	return event.ContextDecision(event.UserPromptExpansion, text)
}

// Block returns a decision:"block" result with an agent-facing reason.
func (results) Block(reason string) event.DecisionOutput {
	return event.BlockDecision(event.UserPromptExpansion, reason)
}
