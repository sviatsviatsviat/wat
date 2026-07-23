package claude

import (
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/agent/subagentstart"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/agent/subagentstop"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/agent/taskcompleted"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/agent/taskcreated"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/agent/teammateidle"
)

// SubagentStart is the SubagentStart hook event.
type SubagentStart = subagentstart.Event

// SubagentStartResults is the hook-scoped response builder for SubagentStart.
type SubagentStartResults = subagentstart.Results

// SubagentStop is the SubagentStop hook event.
type SubagentStop = subagentstop.Event

// TaskCreated is the TaskCreated hook event.
type TaskCreated = taskcreated.Event

// TaskCreatedResults is the hook-scoped response builder for TaskCreated.
type TaskCreatedResults = taskcreated.Results

// TaskCompleted is the TaskCompleted hook event.
type TaskCompleted = taskcompleted.Event

// TaskCompletedResults is the hook-scoped response builder for TaskCompleted.
type TaskCompletedResults = taskcompleted.Results

// TeammateIdle is the TeammateIdle hook event.
type TeammateIdle = teammateidle.Event
