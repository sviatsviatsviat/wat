package agnostic

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/claude"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/copilot"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/cursor"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
)

// SubagentStartEvent is the normalized view of a SubagentStart hook invocation.
type SubagentStartEvent = model.SubagentStartEvent

// SubagentStartHandler handles observe-only SubagentStart events.
type SubagentStartHandler = model.SubagentStartHandler

// OnSubagentStart registers an observe-only handler for SubagentStart events.
func (c *hooks) OnSubagentStart(fn SubagentStartHandler) *hooks {
	if fn == nil {
		return c
	}
	return c.appendParts(
		claude.RegisterSubagentStart(fn),
		copilot.RegisterSubagentStart(fn),
		cursor.RegisterSubagentStart(fn),
	)
}
