package claude

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// Stop is the Stop hook event.
type Stop struct {
	Envelope
	// StopHookActive is true when a stop hook already forced continuation.
	StopHookActive bool `json:"stop_hook_active"`
	// LastAssistantMessage is the final assistant text of the turn.
	LastAssistantMessage string `json:"last_assistant_message"`
}

// EventName returns the hook event name.
func (Stop) EventName() string { return EventStop }

func init() {
	registerDecoder(EventStop, decodeAs[Stop])
}

// StopOutput is the response for Stop and SubagentStop events.
type StopOutput struct {
	Common
	// Block keeps the agent working when true.
	Block bool
	// Reason is fed back to Claude as the next instruction.
	Reason string
	// AdditionalContext is non-error feedback that continues the conversation.
	AdditionalContext string
}

func (o StopOutput) isZero() bool {
	return o.Common.isZero() && !o.Block && o.Reason == "" && o.AdditionalContext == ""
}

// StopResults is the hook-scoped response builder supplied to Chain handlers by registration.
type StopResults interface {
	// FollowUp blocks completion and feeds reason back to Claude.
	FollowUp(reason string) StopOutput
	isStopResults()
}

type stopResults struct{}

func (stopResults) isStopResults() {}

// FollowUp blocks completion and feeds reason back to Claude.
func (stopResults) FollowUp(reason string) StopOutput {
	return StopOutput{Block: true, Reason: reason}
}

// Stop registers a Stop handler.
func (c *Chain) Stop(fn func(context.Context, StopHook, StopResults) (StopOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev Stop) (StopOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), stopResults{})
	})
	return &Chain{}
}
