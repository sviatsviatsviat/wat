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
func OnSubagentStart(fn SubagentStartHandler) *Chain {
	if fn == nil {
		return &Chain{}
	}
	claude.RegisterSubagentStart(fn)
	copilot.RegisterSubagentStart(fn)
	cursor.RegisterSubagentStart(fn)
	return &Chain{}
}

// OnSubagentStart registers another observe-only SubagentStart handler on the chain.
func (c *Chain) OnSubagentStart(fn SubagentStartHandler) *Chain {
	return OnSubagentStart(fn)
}
