package agnostic

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/claude"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/copilot"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/cursor"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
)

// SessionEndEvent is the normalized view of a SessionEnd hook invocation.
type SessionEndEvent = model.SessionEndEvent

// SessionEndHandler handles observe-only SessionEnd events.
type SessionEndHandler = model.SessionEndHandler

// OnSessionEnd registers an observe-only handler for SessionEnd events.
func (c *hooks) OnSessionEnd(fn SessionEndHandler) *hooks {
	if fn == nil {
		return c
	}
	return c.appendParts(
		claude.RegisterSessionEnd(fn),
		copilot.RegisterSessionEnd(fn),
		cursor.RegisterSessionEnd(fn),
	)
}
