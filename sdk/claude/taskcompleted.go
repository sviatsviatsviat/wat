package claude

import (
	"context"
	"encoding/json"

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
	registerDecoder(EventTaskCompleted, decodeAs[TaskCompleted])
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
	return commonOutput{additionalContext: text}
}

// OnTaskCompleted registers a TaskCompleted handler.
func OnTaskCompleted(fn func(context.Context, Hook[TaskCompleted], TaskCompletedResults) (CommonOutput, error)) *chain {
	return (&chain{}).TaskCompleted(fn)
}

// TaskCompleted registers another TaskCompleted handler on the chain.
func (c *chain) TaskCompleted(fn func(context.Context, Hook[TaskCompleted], TaskCompletedResults) (CommonOutput, error)) *chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev TaskCompleted) (CommonOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), taskCompletedResults{})
	})
	return c
}
