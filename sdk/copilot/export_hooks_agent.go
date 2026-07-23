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

// AgentStop is the Stop (agentStop) hook event.
type AgentStop = agentstop.Event

// StopOutput is the response for agentStop and subagentStop events.
type StopOutput = agentstop.Output

// StopResults is the hook-scoped response builder for agentStop and subagentStop.
type StopResults = agentstop.Results
