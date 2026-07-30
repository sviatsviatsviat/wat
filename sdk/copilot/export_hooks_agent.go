package copilot

import (
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/hooks/agent/agentstop"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/hooks/agent/subagentstart"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/hooks/agent/subagentstop"
)

// SubagentStart is the SubagentStart hook event.
type SubagentStart = subagentstart.Event

// SubagentStartOutput is the response for SubagentStart events.
type SubagentStartOutput = subagentstart.Output

// SubagentStartResults is the hook-scoped response builder for SubagentStart.
type SubagentStartResults = subagentstart.Results

// SubagentStop is the SubagentStop hook event.
type SubagentStop = subagentstop.Event

// SubagentStopOutput is the response for SubagentStop events.
type SubagentStopOutput = subagentstop.Output

// SubagentStopResults is the hook-scoped response builder for SubagentStop.
// Use FollowUp to block completion, or ModifiedResponse / WithModifiedResponse
// to rewrite the text returned to the parent when the subagent completes.
type SubagentStopResults = subagentstop.Results

// AgentStop is the Stop (agentStop) hook event.
type AgentStop = agentstop.Event

// StopOutput is the response for AgentStop events.
type StopOutput = agentstop.Output

// StopResults is the hook-scoped response builder for AgentStop.
type StopResults = agentstop.Results
