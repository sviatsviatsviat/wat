package claude

import (
	"context"
	"encoding/json"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// TaskCreated is the TaskCreated hook event.
type TaskCreated struct {
	Envelope
	// Task is the task payload JSON.
	Task json.RawMessage `json:"task"`
}

// EventName returns the hook event name.
func (TaskCreated) EventName() string { return EventTaskCreated }

func init() {
	registerDecoder(EventTaskCreated, decodeAs[TaskCreated])
}

// TaskCreatedResults is the hook-scoped response builder supplied to Chain handlers by registration.
type TaskCreatedResults interface {
	// Context returns a context-injection-only TaskCreated result.
	Context(text string) CommonOutput
	isTaskCreatedResults()
}

type taskCreatedResults struct{}

func (taskCreatedResults) isTaskCreatedResults() {}

// Context returns a context-injection-only TaskCreated result.
func (taskCreatedResults) Context(text string) CommonOutput {
	return commonOutput{additionalContext: text}
}

// TaskCreated registers a TaskCreated handler.
func (c *Chain) TaskCreated(fn func(context.Context, Hook[TaskCreated], TaskCreatedResults) (CommonOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev TaskCreated) (CommonOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), taskCreatedResults{})
	})
	return c
}
