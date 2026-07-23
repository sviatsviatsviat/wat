package cursor

import (
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/agent/afteragentresponse"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/agent/afteragentthought"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/agent/stopevent"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/agent/subagentstart"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/agent/subagentstop"
)

// Stop is the stop hook event.
type Stop = stopevent.Event

// StopOutput is the response for stop and subagentStop events.
type StopOutput = stopevent.Output

// StopResults is the hook-scoped response builder for Stop and SubagentStop.
type StopResults = stopevent.Results

// SubagentStart is the subagentStart hook event.
type SubagentStart = subagentstart.Event

// SubagentStartResults is the hook-scoped response builder for SubagentStart.
type SubagentStartResults = subagentstart.Results

// SubagentStop is the subagentStop hook event.
type SubagentStop = subagentstop.Event

// AfterAgentResponse is the afterAgentResponse hook event.
type AfterAgentResponse = afteragentresponse.Event

// AfterAgentThought is the afterAgentThought hook event.
type AfterAgentThought = afteragentthought.Event
