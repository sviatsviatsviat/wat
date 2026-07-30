package claude

import (
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/session/sessionend"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/session/sessionstart"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/session/setup"
)

// SessionStart is the SessionStart hook event.
type SessionStart = sessionstart.Event

// SessionStartOutput is the response for SessionStart events.
type SessionStartOutput = sessionstart.Output

// SessionStartResults is the hook-scoped response builder for SessionStart.
type SessionStartResults = sessionstart.Results

// SessionEnd is the SessionEnd hook event.
type SessionEnd = sessionend.Event

// Setup is the Setup hook event.
type Setup = setup.Event

// SetupOutput is the response for Setup events.
type SetupOutput = setup.Output

// SetupResults is the hook-scoped response builder for Setup.
type SetupResults = setup.Results
