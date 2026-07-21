package cursor

import (
	"context"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// SubagentStop is the subagentStop hook event.
type SubagentStop struct {
	Envelope
	// SubagentID is the subagent identifier.
	SubagentID string `json:"subagent_id"`
	// SubagentType is the subagent type.
	SubagentType string `json:"subagent_type"`
	// Task is the subagent task description.
	Task string `json:"task"`
	// Summary is the subagent result summary.
	Summary string `json:"summary"`
	// Status is the subagent stop status.
	Status string `json:"status"`
	// LoopCount is the stop-loop iteration count.
	LoopCount int `json:"loop_count"`
	// AgentTranscriptPath is the subagent transcript path when present.
	AgentTranscriptPath *string `json:"agent_transcript_path"`
}

// EventName returns the canonical hook event name.
func (SubagentStop) EventName() string { return EventSubagentStop }

func init() {
	codec.Register(EventSubagentStop, hookkit.EventDecoder[SubagentStop](codec))
}

// SubagentStop registers a SubagentStop handler on the chain.
func (c *chain) SubagentStop(fn func(context.Context, run.Hook[SubagentStop], StopResults) (StopOutput, error)) *chain {
	if fn == nil {
		return c
	}
	c.reg.RegisterHandler(Dialect, run.Handler(func(ctx context.Context, hook run.Hook[SubagentStop]) (StopOutput, error) {
		return fn(ctx, hook, stopResults{})
	}))
	return c
}
