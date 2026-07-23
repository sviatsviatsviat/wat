package cursor

import (
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/tools"
)

// Input is the tool input payload on a Cursor hook event.
type Input = tools.Input

// ShellInput is the input schema for the Shell tool.
type ShellInput = tools.ShellInput

// ReadInput is the input schema for the Read tool.
type ReadInput = tools.ReadInput

// EditInput is the input schema for file edit tools.
type EditInput = tools.EditInput

const (
	// ToolShell is the shell execution tool.
	ToolShell = tools.ToolShell
	// ToolRead is the file read tool.
	ToolRead = tools.ToolRead
	// ToolWrite is the file write tool.
	ToolWrite = tools.ToolWrite
)
