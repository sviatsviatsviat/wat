package agnostic

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/claude"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/copilot"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/cursor"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
)

// StopEvent is the normalized view of Stop and SubagentStop hook invocations.
type StopEvent = model.StopEvent

// StopHook is the handler context for portable Stop and SubagentStop events.
type StopHook = model.StopHook

// StopResult is the portable hook response for Stop and SubagentStop events.
type StopResult = model.StopResult

// StopResults is the hook-scoped response builder supplied to OnStop and OnSubagentStop handlers by registration.
type StopResults = model.StopResults

// StopHandler handles portable Stop and SubagentStop events.
type StopHandler = model.StopHandler

// OnStop registers a handler for Stop events across all agents.
func OnStop(fn StopHandler) *Chain {
	if fn == nil {
		return &Chain{}
	}
	claude.RegisterStop(fn)
	copilot.RegisterStop(fn)
	cursor.RegisterStop(fn)
	return &Chain{}
}

// OnSubagentStop registers a handler for SubagentStop events across all agents.
func OnSubagentStop(fn StopHandler) *Chain {
	if fn == nil {
		return &Chain{}
	}
	claude.RegisterSubagentStop(fn)
	copilot.RegisterSubagentStop(fn)
	cursor.RegisterSubagentStop(fn)
	return &Chain{}
}

// OnStop registers another Stop handler on the chain.
func (c *Chain) OnStop(fn StopHandler) *Chain {
	return OnStop(fn)
}

// OnSubagentStop registers another SubagentStop handler on the chain.
func (c *Chain) OnSubagentStop(fn StopHandler) *Chain {
	return OnSubagentStop(fn)
}
