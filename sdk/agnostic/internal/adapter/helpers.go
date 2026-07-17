// Package adapter provides shared helpers for agent dialect codecs.
package adapter

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
)

// NewToolCall builds a ToolCall from native tool metadata and extracts a shell
// command when the tool is a shell execution.
func NewToolCall(native string, input json.RawMessage, id string) *model.ToolCall {
	name, mcp := model.NormalizeToolName(native)
	tc := &model.ToolCall{
		Name:   name,
		Native: native,
		ID:     id,
		MCP:    mcp,
	}
	if name != "" || native != "" || len(input) > 0 {
		tc.Input = model.NewToolInput(name, native, input)
	}
	if name == model.ToolBash {
		tc.Shell = hookkit.ExtractShellCommand(input)
	}
	return tc
}
