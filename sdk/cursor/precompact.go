package cursor

import (
	"context"
	"encoding/json"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// PreCompact is the preCompact hook event.
type PreCompact struct {
	Envelope
	// Trigger is the compaction trigger.
	Trigger string `json:"trigger"`
}

// EventName returns the canonical hook event name.
func (PreCompact) EventName() string { return EventPreCompact }

// PreCompactOutput is the response for preCompact events.
type PreCompactOutput struct {
	// UserMessage is shown to the user.
	UserMessage string
}

func (o PreCompactOutput) isZero() bool {
	return o.UserMessage == ""
}

// PreCompactResults is the hook-scoped response builder supplied to Chain handlers by registration.
type PreCompactResults interface {
	UserMessage(text string) PreCompactOutput
	isPreCompactResults()
}

type preCompactResults struct{}

func (preCompactResults) isPreCompactResults() {}

// UserMessage returns a preCompact result with a user-facing message.
func (preCompactResults) UserMessage(text string) PreCompactOutput {
	return PreCompactOutput{UserMessage: text}
}

func encodePreCompact(o PreCompactOutput) ([]byte, int, error) {
	if o.UserMessage == "" {
		return nil, 0, nil
	}
	out := map[string]any{"user_message": o.UserMessage}
	b, err := json.Marshal(out)
	return b, 0, err
}

func init() {
	registerDecoder(EventPreCompact, decodeAs[PreCompact])
}

// PreCompact registers a preCompact handler.
func (c *Chain) PreCompact(fn func(context.Context, PreCompactHook, PreCompactResults) (PreCompactOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev PreCompact) (PreCompactOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), preCompactResults{})
	})
	return &Chain{}
}
