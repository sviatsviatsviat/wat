package copilot

import (
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/hooks/session/sessionend"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/hooks/session/sessionstart"
)

// SessionStart is the SessionStart hook event.
type SessionStart = sessionstart.Event

// SessionStartOutput is the response for SessionStart events.
type SessionStartOutput = sessionstart.Output

// SessionStartResults is the hook-scoped response builder for SessionStart.
type SessionStartResults = sessionstart.Results

// SessionEnd is the SessionEnd hook event.
type SessionEnd = sessionend.Event
