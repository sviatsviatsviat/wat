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
// On Copilot, this receives agent-scoped Stop payloads only; subagent-scoped
// Stop (agent_name / agent_display_name set) goes to OnSubagentStop.
func (c *chain) OnStop(fn StopHandler) *chain {
	if fn == nil {
		return c
	}
	claude.RegisterStop(fn)
	copilot.RegisterStop(fn)
	cursor.RegisterStop(fn)
	return c
}

// OnSubagentStop registers a handler for SubagentStop events across all agents.
// On Copilot, this also receives wire Stop payloads scoped to a subagent
// (agent_name / agent_display_name set).
func (c *chain) OnSubagentStop(fn StopHandler) *chain {
	if fn == nil {
		return c
	}
	claude.RegisterSubagentStop(fn)
	copilot.RegisterSubagentStop(fn)
	cursor.RegisterSubagentStop(fn)
	return c
}
