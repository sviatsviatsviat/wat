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
// WorkspaceOpen is an app lifecycle hook for the Cursor desktop app and CLI.
// WorkspaceOpen does not run in cloud agents.
type WorkspaceOpen = workspaceopen.Event

// WorkspaceOpenOutput is the response for workspaceOpen events.
type WorkspaceOpenOutput = workspaceopen.Output

// WorkspaceOpenResults is the hook-scoped response builder for WorkspaceOpen.
type WorkspaceOpenResults = workspaceopen.Results
