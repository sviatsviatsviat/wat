package agnostic

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/claude"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/copilot"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/cursor"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
)

// PreToolEvent is the normalized view of a PreTool hook invocation.
type PreToolEvent = model.PreToolEvent

// PreToolResult is the portable hook response for PreTool events.
type PreToolResult = model.PreToolResult

// PreToolResults is the hook-scoped response builder supplied to PreToolHandler by registration.
type PreToolResults = model.PreToolResults

// PreToolHandler handles portable PreTool events.
type PreToolHandler = model.PreToolHandler

// OnPreTool registers a handler for PreTool events across all agents.
func (c *chain) OnPreTool(fn PreToolHandler) *chain {
	if fn == nil {
		return c
	}
	claude.RegisterPreTool(fn)
	copilot.RegisterPreTool(fn)
	cursor.RegisterPreTool(fn)
	cursor.RegisterBeforeReadFile(fn)
	return c
}
