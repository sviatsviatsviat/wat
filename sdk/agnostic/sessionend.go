package agnostic

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/claude"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/copilot"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/cursor"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
)

// SessionEndEvent is the normalized view of a SessionEnd hook invocation.
type SessionEndEvent = model.SessionEndEvent

// SessionEndHook is the handler context for portable SessionEnd events.
type SessionEndHook = model.SessionEndHook

// SessionEndHandler handles observe-only SessionEnd events.
type SessionEndHandler = model.SessionEndHandler

// OnSessionEnd registers an observe-only handler for SessionEnd events.
func OnSessionEnd(fn SessionEndHandler) *chain {
	if fn == nil {
		return &chain{}
	}
	claude.RegisterSessionEnd(fn)
	copilot.RegisterSessionEnd(fn)
	cursor.RegisterSessionEnd(fn)
	return &chain{}
}

// OnSessionEnd registers another observe-only SessionEnd handler on the chain.
func (c *chain) OnSessionEnd(fn SessionEndHandler) *chain {
	return OnSessionEnd(fn)
}
