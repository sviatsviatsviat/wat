package teammateidle

import (
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
)

// Results is the hook-scoped response builder for this event.
type Results interface {
	// Context returns a context-injection-only TeammateIdle result.
	Context(text string) event.ExitBlockOutput
	// Block prevents the teammate from going idle via Claude exit 2 and stderr.
	// Prefer Block over WithContinue(false); continue:false stops the teammate.
	Block(reason string) event.ExitBlockOutput
	isResults()
}

type results struct{}

func (results) isResults() {}

// Context returns a context-injection-only TeammateIdle result.
func (results) Context(text string) event.ExitBlockOutput {
	return event.ContextExitBlock(event.TeammateIdle, text)
}

// Block prevents the teammate from going idle via Claude exit 2 and stderr.
func (results) Block(reason string) event.ExitBlockOutput {
	return event.BlockExitBlock(event.TeammateIdle, reason)
}
