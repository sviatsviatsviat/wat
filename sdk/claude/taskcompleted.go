package claude

import (
	"context"
	"encoding/json"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// TaskCompleted is the TaskCompleted hook event.
type TaskCompleted struct {
	Envelope
	// Task is the task payload JSON.
	Task json.RawMessage `json:"task"`
}

// EventName returns the hook event name.
func (TaskCompleted) EventName() string { return EventTaskCompleted }

func init() {
	codec.Register(EventTaskCompleted, hookkit.EventDecoder[TaskCompleted](codec))
}

// TaskCompletedResults is the hook-scoped response builder supplied to On* handlers by registration.
type TaskCompletedResults interface {
	// Context returns a context-injection-only TaskCompleted result.
	Context(text string) CommonOutput
	isTaskCompletedResults()
}

type taskCompletedResults struct{}

func (taskCompletedResults) isTaskCompletedResults() {}

// Context returns a context-injection-only TaskCompleted result.
func (taskCompletedResults) Context(text string) CommonOutput {
	return commonOutput{eventName: EventTaskCompleted, additionalContext: text}
}

// TaskCompleted registers a TaskCompleted handler on the chain.
func (c *chain) TaskCompleted(fn func(context.Context, run.Hook[TaskCompleted], TaskCompletedResults) (CommonOutput, error)) *chain {
	if fn == nil {
		return c
	}
	registerHandler(c.reg, func(ctx context.Context, ev TaskCompleted) (CommonOutput, error) {
		return fn(ctx, run.NewHook(run.InvocationFrom(ctx), ev), taskCompletedResults{})
	})
	return c
}
