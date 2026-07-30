package taskcompleted

import (
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
)

// Results is the hook-scoped response builder for this event.
type Results interface {
	// Context returns a context-injection-only TaskCompleted result.
	Context(text string) event.ExitBlockOutput
	// Block prevents marking the task completed via Claude exit 2 and stderr.
	// Prefer Block over WithContinue(false); continue:false stops the teammate.
	Block(reason string) event.ExitBlockOutput
	isResults()
}

type results struct{}

func (results) isResults() {}

// Context returns a context-injection-only TaskCompleted result.
func (results) Context(text string) event.ExitBlockOutput {
	return event.ContextExitBlock(event.TaskCompleted, text)
}

// Block prevents marking the task completed via Claude exit 2 and stderr.
func (results) Block(reason string) event.ExitBlockOutput {
	return event.BlockExitBlock(event.TaskCompleted, reason)
}
