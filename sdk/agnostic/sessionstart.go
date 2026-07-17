package agnostic

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/claude"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/copilot"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/cursor"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
)

// SessionStartEvent is the normalized view of a SessionStart hook invocation.
type SessionStartEvent = model.SessionStartEvent

// SessionStartHook is the handler context for portable SessionStart events.
type SessionStartHook = model.SessionStartHook

// SessionStartResult is the portable hook response for SessionStart events.
type SessionStartResult = model.SessionStartResult

// SessionStartResults is the hook-scoped response builder supplied to OnSessionStart handlers by registration.
type SessionStartResults = model.SessionStartResults

// SessionStartHandler handles portable SessionStart events.
type SessionStartHandler = model.SessionStartHandler

// OnSessionStart registers a handler for SessionStart events across all agents.
func OnSessionStart(fn SessionStartHandler) *Chain {
	if fn == nil {
		return &Chain{}
	}
	claude.RegisterSessionStart(fn)
	copilot.RegisterSessionStart(fn)
	cursor.RegisterSessionStart(fn)
	return &Chain{}
}

// OnSessionStart registers another SessionStart handler on the chain.
func (c *Chain) OnSessionStart(fn SessionStartHandler) *Chain {
	return OnSessionStart(fn)
}
