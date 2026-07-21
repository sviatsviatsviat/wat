package copilot

import (
	"context"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// AgentStop is the Stop (agentStop) hook event.
// When the host scopes the stop to a subagent, AgentName and/or AgentDisplayName
// are set; hook authors should check IsSubagent (or those fields) rather than
// expecting a separate SubagentStop wire name.
type AgentStop struct {
	Envelope
	// AgentName is the subagent name when the stop is subagent-scoped.
	AgentName string `json:"agent_name"`
	// AgentDisplayName is the subagent display name when present.
	AgentDisplayName string `json:"agent_display_name"`
	// StopReason is the stop reason.
	StopReason string `json:"stop_reason"`
}

// EventName returns the canonical hook event name.
func (AgentStop) EventName() string { return EventAgentStop }

// IsSubagent reports whether this Stop payload is scoped to a subagent.
func (e AgentStop) IsSubagent() bool {
	return e.AgentName != "" || e.AgentDisplayName != ""
}

// Name returns the agent name when present.
func (e AgentStop) Name() string {
	return e.AgentName
}

// DisplayName returns the agent display name when present.
func (e AgentStop) DisplayName() string {
	return e.AgentDisplayName
}

// Reason returns the stop reason.
func (e AgentStop) Reason() string {
	return e.StopReason
}

func init() {
	codec.Register(EventAgentStop, hookkit.EventDecoder[AgentStop](codec))
}

// AgentStop registers a AgentStop handler on the chain.
func (c *chain) AgentStop(fn func(context.Context, run.Hook[AgentStop], StopResults) (StopOutput, error)) *chain {
	if fn == nil {
		return c
	}
	c.reg.RegisterHandler(Dialect, run.Handler(func(ctx context.Context, hook run.Hook[AgentStop]) (StopOutput, error) {
		return fn(ctx, hook, stopResults{})
	}))
	return c
}
