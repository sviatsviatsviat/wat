package posttoolusefailure

import (
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/tool/posttooluse"
)

// Results is the hook-scoped response builder for PostToolUseFailure.
type Results interface {
	// Context returns recovery guidance for PostToolUseFailure events.
	Context(text string) posttooluse.Output
	isResults()
}

type results struct{}

func (results) isResults() {}

// Context returns recovery guidance for PostToolUseFailure events.
func (results) Context(text string) posttooluse.Output {
	return posttooluse.FailureContext(text)
}
