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

// TaskCompletedResults is the hook-scoped response builder supplied to Chain handlers by registration.
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

// TaskCompleted registers a TaskCompleted handler.
func (c *Chain) TaskCompleted(fn func(context.Context, Hook[TaskCompleted], TaskCompletedResults) (CommonOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(c.registerOwner(), func(ctx context.Context, ev TaskCompleted) (CommonOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), taskCompletedResults{})
	})
	return c
}
