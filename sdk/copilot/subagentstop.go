package copilot

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// SubagentStop is the subagentStop hook event.
type SubagentStop struct {
	Envelope
	// AgentName is the agent name (VS Code).
	AgentName string `json:"agent_name"`
	// AgentNameCamel is the agent name (camelCase).
	AgentNameCamel string `json:"agentName"`
	// AgentDisplayName is the display name (VS Code).
	AgentDisplayName string `json:"agent_display_name"`
	// AgentDisplayNameCamel is the display name (camelCase).
	AgentDisplayNameCamel string `json:"agentDisplayName"`
	// StopReason is the stop reason (VS Code).
	StopReason string `json:"stop_reason"`
	// StopReasonCamel is the stop reason (camelCase).
	StopReasonCamel string `json:"stopReason"`
}

// EventName returns the canonical hook event name.
func (SubagentStop) EventName() string { return EventSubagentStop }

// Name returns the agent name from either wire format.
func (e SubagentStop) Name() string {
	if e.AgentName != "" {
		return e.AgentName
	}
	return e.AgentNameCamel
}

// DisplayName returns the agent display name from either wire format.
func (e SubagentStop) DisplayName() string {
	if e.AgentDisplayName != "" {
		return e.AgentDisplayName
	}
	return e.AgentDisplayNameCamel
}

// Reason returns the stop reason from either wire format.
func (e SubagentStop) Reason() string {
	if e.StopReason != "" {
		return e.StopReason
	}
	return e.StopReasonCamel
}

func init() {
	registerDecoder(EventSubagentStop, decodeAs[SubagentStop])
}

// SubagentStop registers a SubagentStop handler.
func (c *Chain) SubagentStop(fn func(context.Context, SubagentStopHook, StopResults) (StopOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev SubagentStop) (StopOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), stopResults{})
	})
	return &Chain{}
}
