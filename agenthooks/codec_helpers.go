package agenthooks

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/copilothook"
)

func cloneRaw(raw json.RawMessage) json.RawMessage {
	return copilothook.CloneRaw(raw)
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
		tc.Shell = extractShellCommand(input)
	}
	return tc
}

// extractShellCommand reads the command string from a shell tool input object.
func extractShellCommand(input json.RawMessage) string {
	if len(input) == 0 {
		return ""
	}
	var args struct {
		Command string `json:"command"`
	}
	if json.Unmarshal(input, &args) != nil {
		return ""
	}
	return args.Command
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
