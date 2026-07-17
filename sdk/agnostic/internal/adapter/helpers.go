// Package adapter provides shared helpers for agent dialect codecs.
package adapter

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
)

// CloneRaw returns a defensive copy of raw JSON.
func CloneRaw(raw json.RawMessage) json.RawMessage {
	return hookkit.CloneRaw(raw)
}

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

// RawToText extracts a best-effort textual form of a tool_response value.
func RawToText(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	var s string
	if json.Unmarshal(raw, &s) == nil {
		return s
	}
	return string(raw)
}
