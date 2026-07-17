package cursor

import (
	"context"

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
	registerDecoder(EventSubagentStop, decodeAs[SubagentStop])
}

// SubagentStop registers a subagentStop handler.
func (c *Chain) SubagentStop(fn func(context.Context, Hook[SubagentStop], StopResults) (StopOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev SubagentStop) (StopOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), stopResults{})
	})
	return c
}
