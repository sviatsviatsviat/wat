package cursor

import (
	"context"
	"encoding/json"

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
type StopOutput struct {
	// FollowUpMessage is sent back to the agent as the next instruction.
	FollowUpMessage string
}

func (o StopOutput) isZero() bool {
	return o.FollowUpMessage == ""
}

// StopResults is the hook-scoped response builder supplied to Chain handlers by registration.
type StopResults interface {
	// FollowUp blocks completion and feeds a follow-up instruction to the agent.
	FollowUp(text string) StopOutput
	isStopResults()
}

type stopResults struct{}

func (stopResults) isStopResults() {}

// FollowUp blocks completion and feeds a follow-up instruction to the agent.
func (stopResults) FollowUp(text string) StopOutput {
	return StopOutput{FollowUpMessage: text}
}

func encodeStop(o StopOutput) ([]byte, int, error) {
	if o.FollowUpMessage == "" {
		return nil, 0, nil
	}
	out := map[string]any{"followup_message": o.FollowUpMessage}
	b, err := json.Marshal(out)
	return b, 0, err
}

func init() {
	registerDecoder(EventStop, decodeAs[Stop])
}

// Stop registers a stop handler.
func (c *Chain) Stop(fn func(context.Context, StopHook, StopResults) (StopOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev Stop) (StopOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), stopResults{})
	})
	return &Chain{}
}
