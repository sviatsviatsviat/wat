package claude

import (
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/stop/stopevent"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/stop/stopfailure"
)

// Stop is the Stop hook event.
type Stop = stopevent.Event

// StopOutput is the response for Stop and SubagentStop events.
type StopOutput = stopevent.Output

// StopResults is the hook-scoped response builder for Stop and SubagentStop.
type StopResults = stopevent.Results

// StopFailure is the StopFailure hook event.
type StopFailure = stopfailure.Event
