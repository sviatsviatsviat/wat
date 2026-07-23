package agnostic

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/claude"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/copilot"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/cursor"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
)

// SessionStartEvent is the normalized view of a SessionStart hook invocation.
type SessionStartEvent = model.SessionStartEvent

// SessionStartResult is the portable hook response for SessionStart events.
type SessionStartResult = model.SessionStartResult

// SessionStartResults is the hook-scoped response builder supplied to OnSessionStart handlers by registration.
type SessionStartResults = model.SessionStartResults

// SessionStartHandler handles portable SessionStart events.
type SessionStartHandler = model.SessionStartHandler

// OnSessionStart registers a handler for SessionStart events across all agents.
func (c *hooks) OnSessionStart(fn SessionStartHandler) *hooks {
	if fn == nil {
		return c
	}
	return c.appendParts(
		claude.RegisterSessionStart(fn),
		copilot.RegisterSessionStart(fn),
		cursor.RegisterSessionStart(fn),
	)
}
