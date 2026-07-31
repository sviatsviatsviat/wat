package configchange

import (
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
)

// Results is the hook-scoped response builder for this event.
type Results interface {
	// Context returns a context-injection-only ConfigChange result.
	Context(text string) event.DecisionOutput
	// Block returns a decision:"block" result that rejects the config change
	// (except policy_settings, which Claude does not allow hooks to block).
	// Encodes with SuccessExit; Claude processes JSON only on exit 0.
	Block(reason string) event.DecisionOutput
	isResults()
}

type results struct{}

func (results) isResults() {}

// Context returns a context-injection-only ConfigChange result.
func (results) Context(text string) event.DecisionOutput {
	return event.ContextDecision(event.ConfigChange, text)
}

// Block returns a decision:"block" result that rejects the config change.
func (results) Block(reason string) event.DecisionOutput {
	return event.BlockDecision(event.ConfigChange, reason)
}
