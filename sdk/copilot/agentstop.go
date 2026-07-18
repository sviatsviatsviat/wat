package copilot

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// AgentStop is the Stop (agentStop) hook event.
type AgentStop struct {
	Envelope
	// StopReason is the stop reason.
	StopReason string `json:"stop_reason"`
}

// EventName returns the canonical hook event name.
func (AgentStop) EventName() string { return EventAgentStop }

// Reason returns the stop reason.
func (e AgentStop) Reason() string {
	return e.StopReason
}

func init() {
	registerDecoder(EventAgentStop, decodeAs[AgentStop])
}

// OnAgentStop registers an AgentStop handler.
func OnAgentStop(fn func(context.Context, Hook[AgentStop], StopResults) (StopOutput, error)) *chain {
	return (&chain{}).AgentStop(fn)
}

// AgentStop registers another AgentStop handler on the chain.
func (c *chain) AgentStop(fn func(context.Context, Hook[AgentStop], StopResults) (StopOutput, error)) *chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev AgentStop) (StopOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), stopResults{})
	})
	return c
}
