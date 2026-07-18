package agnostic

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/claude"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/copilot"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/cursor"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
)

// PreToolEvent is the normalized view of a PreTool hook invocation.
type PreToolEvent = model.PreToolEvent

// PreToolHook is the handler context for portable PreTool events.
type PreToolHook = model.PreToolHook

// PreToolResult is the portable hook response for PreTool events.
type PreToolResult = model.PreToolResult

// PreToolResults is the hook-scoped response builder supplied to PreToolHandler by registration.
type PreToolResults = model.PreToolResults

// PreToolHandler handles portable PreTool events.
type PreToolHandler = model.PreToolHandler

// OnPreTool registers a handler for PreTool events across all agents.
func OnPreTool(fn PreToolHandler) *chain {
	if fn == nil {
		return &chain{}
	}
	claude.RegisterPreTool(fn)
	copilot.RegisterPreTool(fn)
	cursor.RegisterPreTool(fn)
	cursor.RegisterBeforeReadFile(fn)
	return &chain{}
}

// OnPreTool registers another PreTool handler on the chain.
func (c *chain) OnPreTool(fn PreToolHandler) *chain {
	return OnPreTool(fn)
}
