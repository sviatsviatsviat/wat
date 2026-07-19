package copilot

import (
	"context"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// SubagentStop is the SubagentStop hook event when the host emits
// hook_event_name "SubagentStop". VS Code-style Stop payloads that include
// agent_name decode as AgentStop instead; use AgentStop.IsSubagent there.
type SubagentStop struct {
	Envelope
	// AgentName is the agent name.
	AgentName string `json:"agent_name"`
	// AgentDisplayName is the display name.
	AgentDisplayName string `json:"agent_display_name"`
	// StopReason is the stop reason.
	StopReason string `json:"stop_reason"`
}

// EventName returns the canonical hook event name.
func (SubagentStop) EventName() string { return EventSubagentStop }

// Name returns the agent name.
func (e SubagentStop) Name() string {
	return e.AgentName
}

// DisplayName returns the agent display name.
func (e SubagentStop) DisplayName() string {
	return e.AgentDisplayName
}

// Reason returns the stop reason.
func (e SubagentStop) Reason() string {
	return e.StopReason
}

func init() {
	codec.Register(EventSubagentStop, hookkit.EventDecoder[SubagentStop](codec))
}

// SubagentStop registers a SubagentStop handler on the chain.
func (c *chain) SubagentStop(fn func(context.Context, run.Hook[SubagentStop], StopResults) (StopOutput, error)) *chain {
	if fn == nil {
		return c
	}
	registerHandler(c.reg, func(ctx context.Context, ev SubagentStop) (StopOutput, error) {
		return fn(ctx, run.NewHook(run.InvocationFrom(ctx), ev), stopResults{})
	})
	return c
}
