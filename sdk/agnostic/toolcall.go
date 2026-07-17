package agnostic

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/tools"
)

// newToolCall builds a ToolCall from native tool metadata and extracts a shell
// command when the tool is a shell execution.
func newToolCall(native string, input json.RawMessage, id string) *ToolCall {
	name, mcp := hookkit.NormalizeToolName(native)
	tc := &ToolCall{
		Name:   name,
		Native: native,
		ID:     id,
		MCP:    mcp,
	}
	if name != "" || native != "" || len(input) > 0 {
		tc.Input = tools.NewInput(name, native, input)
	}
	if name == tools.ToolBash {
		tc.Shell = hookkit.ExtractShellCommand(input)
	}
	return tc
}
