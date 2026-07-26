package model

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/tools"
)

// ToolCall describes the tool invocation a pre/post tool event refers to.
type ToolCall struct {
	// Name is the normalized tool name (bash, edit, write, read, glob, grep, task, web_fetch, ...).
	Name string
	// Native is the exact original tool name (Bash vs bash vs Shell).
	Native string
	// Input is the typed tool input for this call.
	Input tools.Input
	// ID is the tool_use_id or tool call id when available.
	ID string
	// Shell is the command string for shell execution tools when available.
	Shell string
	// MCP is true when the call targets an MCP server tool.
	MCP bool
}

// ToolResult describes the outcome of a completed or failed tool call.
type ToolResult struct {
	// Text is the textual result as seen by the model, when available.
	Text string
	// Raw is the native result payload JSON.
	Raw json.RawMessage
	// Error is the failure message for post-tool failure events.
	Error string
	// FailureType is the failure category (error, timeout, permission_denied).
	FailureType string
	// DurationMs is the tool execution duration in milliseconds when available.
	DurationMs int64
	// IsInterrupt is true when the failure was caused by a user interrupt
	// or cancellation, when the native payload supplies it.
	IsInterrupt bool
}

// NewToolCall builds a ToolCall from native tool metadata and extracts a shell
// command when the tool is a shell execution.
func NewToolCall(native string, input json.RawMessage, id string) *ToolCall {
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
