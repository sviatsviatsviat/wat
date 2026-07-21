package claude

import (
	"context"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// SubagentStart is the SubagentStart hook event.
type SubagentStart struct {
	Envelope
}

// EventName returns the hook event name.
func (SubagentStart) EventName() string { return EventSubagentStart }

func init() {
	codec.Register(EventSubagentStart, hookkit.EventDecoder[SubagentStart](codec))
}

// SubagentStartResults is the hook-scoped response builder supplied to On* handlers by registration.
type SubagentStartResults interface {
	// Context returns a context-injection-only SubagentStart result.
	Context(text string) CommonOutput
	isSubagentStartResults()
}

type subagentStartResults struct{}

func (subagentStartResults) isSubagentStartResults() {}

// Context returns a context-injection-only SubagentStart result.
func (subagentStartResults) Context(text string) CommonOutput {
	return commonOutput{eventName: EventSubagentStart, additionalContext: text}
}

// SubagentStart registers a SubagentStart handler on the chain.
func (c *chain) SubagentStart(fn func(context.Context, run.Hook[SubagentStart], SubagentStartResults) (CommonOutput, error)) *chain {
	if fn == nil {
		return c
	}
	c.reg.RegisterHandler(Dialect, run.Handler(func(ctx context.Context, hook run.Hook[SubagentStart]) (CommonOutput, error) {
		return fn(ctx, hook, subagentStartResults{})
	}))
	return c
}
