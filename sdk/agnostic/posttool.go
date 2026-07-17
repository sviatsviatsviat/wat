package agnostic

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/claude"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/copilot"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/cursor"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
)

// PostToolEvent is the normalized view of a PostTool hook invocation.
type PostToolEvent = model.PostToolEvent

// PostToolHook is the handler context for portable PostTool events.
type PostToolHook = model.PostToolHook

// PostToolResult is the portable hook response for PostTool events.
type PostToolResult = model.PostToolResult

// PostToolResults is the hook-scoped response builder supplied to PostToolHandler by registration.
type PostToolResults = model.PostToolResults

// PostToolHandler handles portable PostTool events.
type PostToolHandler = model.PostToolHandler

// OnPostTool registers a handler for PostTool events across all agents.
func OnPostTool(fn PostToolHandler) *Chain {
	if fn == nil {
		return &Chain{}
	}
	claude.RegisterPostTool(fn)
	copilot.RegisterPostTool(fn)
	cursor.RegisterPostTool(fn)
	return &Chain{}
}

// OnPostTool registers another PostTool handler on the chain.
func (c *Chain) OnPostTool(fn PostToolHandler) *Chain {
	return OnPostTool(fn)
}
