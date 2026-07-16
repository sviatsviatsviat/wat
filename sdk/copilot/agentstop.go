package copilot

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// AgentStop is the agentStop hook event.
type AgentStop struct {
	Envelope
	// StopReason is the stop reason (VS Code).
	StopReason string `json:"stop_reason"`
	// StopReasonCamel is the stop reason (camelCase).
	StopReasonCamel string `json:"stopReason"`
}

// EventName returns the canonical hook event name.
func (AgentStop) EventName() string { return EventAgentStop }

// Reason returns the stop reason from either wire format.
func (e AgentStop) Reason() string {
	if e.StopReason != "" {
		return e.StopReason
	}
	return e.StopReasonCamel
}

func init() {
	registerDecoder(EventAgentStop, decodeAs[AgentStop])
}

// AgentStop registers an AgentStop handler.
func (c *Chain) AgentStop(fn func(context.Context, Hook[AgentStop], StopResults) (StopOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev AgentStop) (StopOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), stopResults{})
	})
	return &Chain{}
}
