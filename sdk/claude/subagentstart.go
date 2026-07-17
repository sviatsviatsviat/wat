package claude

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// SubagentStart is the SubagentStart hook event.
type SubagentStart struct {
	Envelope
}

// EventName returns the hook event name.
func (SubagentStart) EventName() string { return EventSubagentStart }

func init() {
	registerDecoder(EventSubagentStart, decodeAs[SubagentStart])
}

// SubagentStartResults is the hook-scoped response builder supplied to Chain handlers by registration.
type SubagentStartResults interface {
	// Context returns a context-injection-only SubagentStart result.
	Context(text string) CommonOutput
	isSubagentStartResults()
}

type subagentStartResults struct{}

func (subagentStartResults) isSubagentStartResults() {}

// Context returns a context-injection-only SubagentStart result.
func (subagentStartResults) Context(text string) CommonOutput {
	return commonOutput{additionalContext: text}
}

// SubagentStart registers a SubagentStart handler.
func (c *Chain) SubagentStart(fn func(context.Context, Hook[SubagentStart], SubagentStartResults) (CommonOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(c.registerOwner(), func(ctx context.Context, ev SubagentStart) (CommonOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), subagentStartResults{})
	})
	return c
}
