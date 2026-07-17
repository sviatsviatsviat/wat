package claude

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// PreCompact is the PreCompact hook event.
type PreCompact struct {
	Envelope
	// Trigger is the compaction trigger (manual, auto).
	Trigger string `json:"trigger"`
	// CustomInstructions are user-provided compaction instructions.
	CustomInstructions string `json:"custom_instructions"`
}

// EventName returns the hook event name.
func (PreCompact) EventName() string { return EventPreCompact }

func init() {
	registerDecoder(EventPreCompact, decodeAs[PreCompact])
}

// PreCompactResults is the hook-scoped response builder supplied to Chain handlers by registration.
type PreCompactResults interface {
	// Context returns a context-injection-only PreCompact result.
	Context(text string) CommonOutput
	isPreCompactResults()
}

type preCompactResults struct{}

func (preCompactResults) isPreCompactResults() {}

// Context returns a context-injection-only PreCompact result.
func (preCompactResults) Context(text string) CommonOutput {
	return commonOutput{additionalContext: text}
}

// PreCompact registers a PreCompact handler.
func (c *Chain) PreCompact(fn func(context.Context, Hook[PreCompact], PreCompactResults) (CommonOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev PreCompact) (CommonOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), preCompactResults{})
	})
	return c
}
