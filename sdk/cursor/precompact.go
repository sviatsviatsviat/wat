package cursor

import (
	"context"
	"encoding/json"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

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
	run.Output
	isPreCompactOutput()
}

type preCompactOutput struct {
	userMessage string
}

func (preCompactOutput) isPreCompactOutput() {}

// IsZero reports whether this hook response is empty.
func (o preCompactOutput) IsZero() bool {
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

// Encode renders this output as Cursor stdout JSON.
func (o preCompactOutput) Encode() ([]byte, int, error) {
	if o.userMessage == "" {
		return nil, 0, nil
	}
	out := map[string]any{"user_message": o.userMessage}
	b, err := json.Marshal(out)
	return b, 0, err
}

// Merge combines other into this preCompact output.
func (o preCompactOutput) Merge(other run.Output) (run.Output, []string, error) {
	b, ok := other.(preCompactOutput)
	if !ok {
		return nil, nil, hookkit.ErrMergeType(o, other)
	}
	var warnings []string
	userMessage, w := hookkit.TakeLastString("userMessage", o.userMessage, b.userMessage)
	if w != "" {
		warnings = append(warnings, w)
	}
	return preCompactOutput{userMessage: userMessage}, warnings, nil
}

// Stop reports whether remaining handlers should be skipped.
func (o preCompactOutput) Stop() bool {
	return false
}

func init() {
	codec.Register(EventPreCompact, hookkit.EventDecoder[PreCompact](codec))
}

// PreCompact registers a PreCompact handler on the chain.
func (c *chain) PreCompact(fn func(context.Context, run.Hook[PreCompact], PreCompactResults) (PreCompactOutput, error)) *chain {
	if fn == nil {
		return c
	}
	c.reg.RegisterHandler(Dialect, run.Handler(func(ctx context.Context, hook run.Hook[PreCompact]) (PreCompactOutput, error) {
		return fn(ctx, hook, preCompactResults{})
	}))
	return c
}
