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
// Construct via PreCompactResults builders. A nil value is a no-op.
type PreCompactOutput interface {
	isPreCompactOutput()
}

type preCompactOutput struct {
	userMessage string
}

func (preCompactOutput) isPreCompactOutput() {}

func (o preCompactOutput) isZero() bool {
	return o.userMessage == ""
}

// PreCompactResults is the hook-scoped response builder supplied to On* handlers by registration.
type PreCompactResults interface {
	// UserMessage returns a preCompact result with a user-facing message.
	UserMessage(text string) PreCompactOutput
	// Noop returns an empty response (silent stdout). Prefer nil from handlers when not chaining.
	Noop() PreCompactOutput
	isPreCompactResults()
}

type preCompactResults struct{}

func (preCompactResults) isPreCompactResults() {}

// UserMessage returns a preCompact result with a user-facing message.
func (preCompactResults) UserMessage(text string) PreCompactOutput {
	return preCompactOutput{userMessage: text}
}

// Noop returns an empty response (silent stdout).
func (preCompactResults) Noop() PreCompactOutput {
	return preCompactOutput{}
}

func (preCompactOutput) allowedEvents() []string {
	return []string{EventPreCompact}
}

func (o preCompactOutput) encode(eventName string) ([]byte, int, error) {
	_ = eventName
	if o.userMessage == "" {
		return nil, 0, nil
	}
	out := map[string]any{"user_message": o.userMessage}
	b, err := json.Marshal(out)
	return b, 0, err
}

func init() {
	registerDecoder(EventPreCompact, decodeAs[PreCompact])
}

// OnPreCompact registers a preCompact handler.
func OnPreCompact(fn func(context.Context, Hook[PreCompact], PreCompactResults) (PreCompactOutput, error)) *chain {
	return (&chain{}).PreCompact(fn)
}

// PreCompact registers another PreCompact handler on the chain.
func (c *chain) PreCompact(fn func(context.Context, Hook[PreCompact], PreCompactResults) (PreCompactOutput, error)) *chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev PreCompact) (PreCompactOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), preCompactResults{})
	})
	return c
}
