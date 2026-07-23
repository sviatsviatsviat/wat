package agnostic

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/claude"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/copilot"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/cursor"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
)

// StopEvent is the normalized view of Stop and SubagentStop hook invocations.
type StopEvent = model.StopEvent

// StopResult is the portable hook response for Stop and SubagentStop events.
type StopResult = model.StopResult

// StopResults is the hook-scoped response builder supplied to OnStop and OnSubagentStop handlers by registration.
type StopResults = model.StopResults

// StopHandler handles portable Stop and SubagentStop events.
type StopHandler = model.StopHandler

// OnStop registers a handler for Stop events across all agents.
// On Copilot, this receives agent-scoped Stop payloads only; subagent-scoped
// Stop (agent_name / agent_display_name set) goes to OnSubagentStop.
func (c *hooks) OnStop(fn StopHandler) *hooks {
	if fn == nil {
		return c
	}
	return c.appendParts(
		claude.RegisterStop(fn),
		copilot.RegisterStop(fn),
		cursor.RegisterStop(fn),
	)
}

// OnSubagentStop registers a handler for SubagentStop events across all agents.
// On Copilot, this also receives wire Stop payloads scoped to a subagent
// (agent_name / agent_display_name set).
func (c *hooks) OnSubagentStop(fn StopHandler) *hooks {
	if fn == nil {
		return c
	}
	return c.appendParts(
		claude.RegisterSubagentStop(fn),
		copilot.RegisterSubagentStop(fn),
		cursor.RegisterSubagentStop(fn),
	)
}
