package cursor

import (
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/session/sessionend"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/session/sessionstart"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/session/workspaceopen"
)

// SessionStart is the sessionStart hook event, including optional
// ComposerMode ("agent", "ask", or "edit").
//
// Cursor runs this hook as fire-and-forget and does not enforce continue or
// user_message. It is not available for cloud agents. See the sessionstart
// Event godoc for host semantics.
type SessionStart = sessionstart.Event

// SessionStartOutput is the response for sessionStart events (env and
// additional_context). continue and user_message are not exposed because the
// host does not enforce them.
type SessionStartOutput = sessionstart.Output

// SessionStartResults is the hook-scoped response builder for SessionStart.
type SessionStartResults = sessionstart.Results

// SessionEnd is the sessionEnd hook event, including duration_ms,
// final_status, and optional error_message.
//
// Cursor runs this hook as fire-and-forget and ignores the response body. It
// is not available for cloud agents. See the sessionend Event godoc for host
// semantics.
type SessionEnd = sessionend.Event

// WorkspaceOpen is the workspaceOpen hook event.
// WorkspaceOpen is an app lifecycle hook for the Cursor desktop app and CLI.
// WorkspaceOpen does not run in cloud agents.
type WorkspaceOpen = workspaceopen.Event

// WorkspaceOpenOutput is the response for workspaceOpen events.
type WorkspaceOpenOutput = workspaceopen.Output

// WorkspaceOpenResults is the hook-scoped response builder for WorkspaceOpen.
type WorkspaceOpenResults = workspaceopen.Results
