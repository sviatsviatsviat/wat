package agnostic

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/claude"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/copilot"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/cursor"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
)

// PostToolEvent is the normalized view of a PostTool hook invocation.
type PostToolEvent = model.PostToolEvent

// PostToolResult is the portable hook response for PostTool events.
type PostToolResult = model.PostToolResult

// PostToolResults is the hook-scoped response builder supplied to PostToolHandler by registration.
type PostToolResults = model.PostToolResults

// PostToolHandler handles portable PostTool events.
type PostToolHandler = model.PostToolHandler

// OnPostTool registers a handler for PostTool events across all agents.
// On Cursor this expands to postToolUse, afterMCPExecution, and afterFileEdit.
// afterMCPExecution and afterFileEdit are observe-only (no host-consumed
// output). afterShellExecution is native observe-only and is not included.
func (c *hooks) OnPostTool(fn PostToolHandler) *hooks {
	if fn == nil {
		return c
	}
	return c.appendParts(
		claude.RegisterPostTool(fn),
		copilot.RegisterPostTool(fn),
		cursor.RegisterPostTool(fn),
	)
}
