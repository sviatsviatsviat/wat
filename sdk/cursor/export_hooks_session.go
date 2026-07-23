package cursor

import (
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/session/sessionend"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/session/sessionstart"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/session/workspaceopen"
)

// SessionStart is the sessionStart hook event.
type SessionStart = sessionstart.Event

// SessionStartOutput is the response for sessionStart events.
type SessionStartOutput = sessionstart.Output

// SessionStartResults is the hook-scoped response builder for SessionStart.
type SessionStartResults = sessionstart.Results

// SessionEnd is the sessionEnd hook event.
type SessionEnd = sessionend.Event

// WorkspaceOpen is the workspaceOpen hook event.
type WorkspaceOpen = workspaceopen.Event
