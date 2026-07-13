package agnostic

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/sdk/internal/hookkit"
)

func cloneRaw(raw json.RawMessage) json.RawMessage {
	return hookkit.CloneRaw(raw)
}

// newToolCall builds a ToolCall from native tool metadata and extracts a shell
// command when the tool is a shell execution.
func newToolCall(native string, input json.RawMessage, id string) *ToolCall {
	name, mcp := NormalizeToolName(native)
	tc := &ToolCall{
		Name:   name,
		Native: native,
		Input:  cloneRaw(input),
		ID:     id,
		MCP:    mcp,
	}
	if name == ToolBash {
		tc.Shell = hookkit.ExtractShellCommand(input)
	}
	return tc
}

// rawToText extracts a best-effort textual form of a tool_response value.
func rawToText(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	var s string
	if json.Unmarshal(raw, &s) == nil {
		return s
	}
	return string(raw)
}
