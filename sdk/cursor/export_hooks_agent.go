package cursor

import (
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/agent/afteragentresponse"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/agent/afteragentthought"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/agent/stopevent"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/agent/subagentstart"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/agent/subagentstop"
)

// Stop is the stop hook event (status, loop_count, and shared envelope fields).
//
// Cursor caps FollowUp auto-submits per script with hooks.json loop_limit
// (default 5; null means unlimited). Use Stop.LoopCount to stay within that
// budget before returning FollowUp.
type Stop = stopevent.Event

// StopOutput is the response for stop and subagentStop events.
// Non-empty FollowUp encodes followup_message with exit 0.
type StopOutput = stopevent.Output

// StopResults is the hook-scoped response builder for Stop and SubagentStop.
// FollowUp auto-submits a user message; for subagentStop Cursor only consumes
// it when the input status is "completed". See Stop for loop_limit semantics.
type StopResults = stopevent.Results

// SubagentStart is the subagentStart hook event.
type SubagentStart = subagentstart.Event

// SubagentStartResults is the hook-scoped response builder for SubagentStart.
// Deny writes user_message and exits 0 (Cursor's subagentStart schema).
type SubagentStartResults = subagentstart.Results

// SubagentStop is the subagentStop hook event, including Cursor telemetry
// fields (description, duration_ms, message_count, tool_call_count,
// modified_files) documented in Cursor Hooks.
type SubagentStop = subagentstop.Event

// AfterAgentResponse is the afterAgentResponse hook event.
type AfterAgentResponse = afteragentresponse.Event

// AfterAgentThought is the afterAgentThought hook event (observe-only).
// DurationMs carries the optional thinking-block duration from the wire.
// Cursor's hooks.json matcher for this event is the fixed value AgentThought.
type AfterAgentThought = afteragentthought.Event
