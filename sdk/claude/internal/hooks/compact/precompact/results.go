package precompact

import (
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
)

// Results is the hook-scoped response builder for this event.
type Results interface {
	// Context returns a context-injection-only PreCompact result.
	Context(text string) event.DecisionOutput
	// Block returns a decision:"block" result that prevents compaction.
	// Encodes with SuccessExit; Claude processes JSON only on exit 0.
	Block(reason string) event.DecisionOutput
	isResults()
}

type results struct{}

func (results) isResults() {}

// Context returns a context-injection-only PreCompact result.
func (results) Context(text string) event.DecisionOutput {
	return event.ContextDecision(event.PreCompact, text)
}

// Block returns a decision:"block" result that prevents compaction.
func (results) Block(reason string) event.DecisionOutput {
	return event.BlockDecision(event.PreCompact, reason)
}
