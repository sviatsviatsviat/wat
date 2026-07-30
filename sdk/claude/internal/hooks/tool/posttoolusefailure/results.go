package posttoolusefailure

import (
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/tool/posttooluse"
)

// Results is the hook-scoped response builder for PostToolUseFailure.
type Results interface {
	// Context returns recovery guidance for PostToolUseFailure events.
	Context(text string) posttooluse.Output
	// Block returns a decision:"block" result with an agent-facing reason.
	// The tool already failed; Claude feeds the reason next to the failure.
	// Encodes with SuccessExit; Claude processes JSON only on exit 0.
	Block(reason string) posttooluse.Output
	isResults()
}

type results struct{}

func (results) isResults() {}

// Context returns recovery guidance for PostToolUseFailure events.
func (results) Context(text string) posttooluse.Output {
	return posttooluse.FailureContext(text)
}

// Block returns a decision:"block" result with an agent-facing reason.
func (results) Block(reason string) posttooluse.Output {
	return posttooluse.FailureBlock(reason)
}
