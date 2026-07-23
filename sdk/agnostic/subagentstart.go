package agnostic

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/claude"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/copilot"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/cursor"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
)

// SubagentStartEvent is the normalized view of a SubagentStart hook invocation.
type SubagentStartEvent = model.SubagentStartEvent

// SubagentStartHook is the handler context for portable SubagentStart events.
type SubagentStartHook = model.SubagentStartHook

// SubagentStartHandler handles observe-only SubagentStart events.
type SubagentStartHandler = model.SubagentStartHandler

// OnSubagentStart registers an observe-only handler for SubagentStart events.
func (c *chain) OnSubagentStart(fn SubagentStartHandler) *chain {
	if fn == nil {
		return c
	}
	claude.RegisterSubagentStart(fn)
	copilot.RegisterSubagentStart(fn)
	cursor.RegisterSubagentStart(fn)
	return c
}
