package claude

import (
	"context"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// SubagentStop is the SubagentStop hook event.
type SubagentStop struct {
	Envelope
	// StopHookActive is true when a stop hook already forced continuation.
	StopHookActive bool `json:"stop_hook_active"`
	// LastAssistantMessage is the final assistant text of the turn.
	LastAssistantMessage string `json:"last_assistant_message"`
	// AgentTranscriptPath is the subagent transcript path when provided.
	AgentTranscriptPath string `json:"agent_transcript_path,omitempty"`
}

// EventName returns the hook event name.
func (SubagentStop) EventName() string { return EventSubagentStop }

func init() {
	codec.Register(EventSubagentStop, hookkit.EventDecoder[SubagentStop](codec))
}

// OnSubagentStop registers a SubagentStop handler.
func OnSubagentStop(fn func(context.Context, run.Hook[SubagentStop], StopResults) (StopOutput, error)) *chain {
	return (&chain{}).SubagentStop(fn)
}

// SubagentStop registers another SubagentStop handler on the chain.
func (c *chain) SubagentStop(fn func(context.Context, run.Hook[SubagentStop], StopResults) (StopOutput, error)) *chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev SubagentStop) (StopOutput, error) {
		return fn(ctx, run.NewHook(run.InvocationFrom(ctx), ev), stopResults{})
	})
	return c
}
