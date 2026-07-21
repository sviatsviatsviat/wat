package cursor

import (
	"context"
	"encoding/json"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// Stop is the stop hook event.
type Stop struct {
	Envelope
	// Status is the stop status.
	Status string `json:"status"`
	// LoopCount is the stop-loop iteration count.
	LoopCount int `json:"loop_count"`
}

// EventName returns the canonical hook event name.
func (Stop) EventName() string { return EventStop }

// StopOutput is the response for stop and subagentStop events.
// Construct via StopResults builders. A nil value is a no-op.
type StopOutput interface {
	run.Output
	isStopOutput()
}

type stopOutput struct {
	followUpMessage string
}

func (stopOutput) isStopOutput() {}

// IsZero reports whether this hook response is empty.
func (o stopOutput) IsZero() bool {
	return o.followUpMessage == ""
}

// StopResults is the hook-scoped response builder supplied to On* handlers by registration.
type StopResults interface {
	// FollowUp blocks completion and feeds a follow-up instruction to the agent.
	FollowUp(text string) StopOutput
	// Noop returns an empty response (silent stdout). Prefer nil from handlers when not chaining.
	Noop() StopOutput
	isStopResults()
}

type stopResults struct{}

func (stopResults) isStopResults() {}

// FollowUp blocks completion and feeds a follow-up instruction to the agent.
func (stopResults) FollowUp(text string) StopOutput {
	return stopOutput{followUpMessage: text}
}

// Noop returns an empty response (silent stdout).
func (stopResults) Noop() StopOutput {
	return stopOutput{}
}

// Encode renders this output as Cursor stdout JSON.
func (o stopOutput) Encode() ([]byte, int, error) {
	if o.followUpMessage == "" {
		return nil, 0, nil
	}
	out := map[string]any{"followup_message": o.followUpMessage}
	b, err := json.Marshal(out)
	return b, 0, err
}

func init() {
	codec.Register(EventStop, hookkit.EventDecoder[Stop](codec))
}

// Stop registers a Stop handler on the chain.
func (c *chain) Stop(fn func(context.Context, run.Hook[Stop], StopResults) (StopOutput, error)) *chain {
	if fn == nil {
		return c
	}
	registerHandler(c.reg, func(ctx context.Context, ev Stop) (StopOutput, error) {
		return fn(ctx, run.NewHook(run.InvocationFrom(ctx), ev), stopResults{})
	})
	return c
}
