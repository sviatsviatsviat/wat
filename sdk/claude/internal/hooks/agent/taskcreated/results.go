package taskcreated

import (
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
)

// Results is the hook-scoped response builder for this event.
type Results interface {
	// Context returns a context-injection-only TaskCreated result.
	Context(text string) event.ExitBlockOutput
	// Block rolls back task creation via Claude exit 2 and stderr reason.
	// Prefer Block over WithContinue(false); continue:false stops the teammate.
	Block(reason string) event.ExitBlockOutput
	isResults()
}

type results struct{}

func (results) isResults() {}

// Context returns a context-injection-only TaskCreated result.
func (results) Context(text string) event.ExitBlockOutput {
	return event.ContextExitBlock(event.TaskCreated, text)
}

// Block rolls back task creation via Claude exit 2 and stderr reason.
func (results) Block(reason string) event.ExitBlockOutput {
	return event.BlockExitBlock(event.TaskCreated, reason)
}
