package copilot

import (
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/tools"
)

// Input is the tool input payload on a Copilot hook event.
type Input = tools.Input

// BashInput is the input schema for the bash tool.
type BashInput = tools.BashInput

// CreateInput is the typed input for the create tool.
type CreateInput = tools.CreateInput

// EditInput is the typed input for the edit tool.
type EditInput = tools.EditInput

// ViewInput is the typed input for the view tool.
type ViewInput = tools.ViewInput

// WebFetchInput is the typed input for the webfetch tool.
type WebFetchInput = tools.WebFetchInput

const (
	// ToolBash is the bash tool name.
	ToolBash = tools.ToolBash
	// ToolPowerShell is the PowerShell tool name.
	ToolPowerShell = tools.ToolPowerShell
	// ToolShell is the shell tool name.
	ToolShell = tools.ToolShell
	// ToolCreate is the create tool name.
	ToolCreate = tools.ToolCreate
	// ToolEdit is the edit tool name.
	ToolEdit = tools.ToolEdit
	// ToolView is the view tool name.
	ToolView = tools.ToolView
	// ToolWebFetch is the webfetch tool name.
	ToolWebFetch = tools.ToolWebFetch
)
