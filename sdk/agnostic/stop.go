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
func OnStop(fn StopHandler) *chain {
	if fn == nil {
		return &chain{}
	}
	claude.RegisterStop(fn)
	copilot.RegisterStop(fn)
	cursor.RegisterStop(fn)
	return &chain{}
}

// OnSubagentStop registers a handler for SubagentStop events across all agents.
func OnSubagentStop(fn StopHandler) *chain {
	if fn == nil {
		return &chain{}
	}
	claude.RegisterSubagentStop(fn)
	copilot.RegisterSubagentStop(fn)
	cursor.RegisterSubagentStop(fn)
	return &chain{}
}

// OnStop registers another Stop handler on the chain.
func (c *chain) OnStop(fn StopHandler) *chain {
	return OnStop(fn)
}

// OnSubagentStop registers another SubagentStop handler on the chain.
func (c *chain) OnSubagentStop(fn StopHandler) *chain {
	return OnSubagentStop(fn)
}
