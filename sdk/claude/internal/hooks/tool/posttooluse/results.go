package posttooluse

import (
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
)

// Results is the hook-scoped response builder for PostToolUse.
type Results interface {
	// Context returns a context-injection-only PostToolUse result.
	Context(text string) Output
	// Block returns a block result with an agent-facing reason.
	Block(reason string) Output
	isResults()
}

type results struct{}

func (results) isResults() {}

// Context returns a context-injection-only PostToolUse result.
func (results) Context(text string) Output {
	return output{eventName: event.PostToolUse, additionalContext: text}
}

// Block returns a block result with an agent-facing reason.
func (results) Block(reason string) Output {
	return output{eventName: event.PostToolUse, block: true, reason: reason}
}

// FailureContext returns a PostToolUseFailure-stamped context result.
func FailureContext(text string) Output {
	return output{eventName: event.PostToolUseFailure, additionalContext: text}
}
