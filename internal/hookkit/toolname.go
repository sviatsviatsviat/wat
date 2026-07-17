package hookkit

import "strings"

// Canonical tool names used by [NormalizeToolName] and agnostic ToolCall.Name.
const (
	// ToolBash is the normalized name for shell execution tools.
	ToolBash = "bash"
	// ToolEdit is the normalized name for file edit tools.
	ToolEdit = "edit"
	// ToolWrite is the normalized name for file write tools.
	ToolWrite = "write"
	// ToolRead is the normalized name for file read tools.
	ToolRead = "read"
	// ToolGlob is the normalized name for glob search tools.
	ToolGlob = "glob"
	// ToolGrep is the normalized name for grep search tools.
	ToolGrep = "grep"
	// ToolTask is the normalized name for subagent or task tools.
	ToolTask = "task"
	// ToolWebFetch is the normalized name for web fetch tools.
	ToolWebFetch = "web_fetch"
	// ToolWebSearch is the normalized name for web search tools.
	ToolWebSearch = "web_search"
	// ToolDelete is the normalized name for file delete tools.
	ToolDelete = "delete"
)

var toolAliases = map[string]string{
	// Claude Code
	"bash": ToolBash, "edit": ToolEdit, "write": ToolWrite, "read": ToolRead,
	"glob": ToolGlob, "grep": ToolGrep, "agent": ToolTask, "task": ToolTask,
	"webfetch": ToolWebFetch, "websearch": ToolWebSearch, "notebookedit": ToolEdit,
	// Copilot
	"powershell": ToolBash, "view": ToolRead, "create": ToolWrite, "web_fetch": ToolWebFetch,
	// Cursor
	"shell": ToolBash, "delete": ToolDelete,
}

// NormalizeToolName maps a native tool name onto the canonical vocabulary.
// MCP tools with a verified namespace report mcp=true and keep the native name:
//   - Claude / Copilot PascalCase: "mcp__<server>__<tool>"
//   - Cursor matcher form: "MCP:<tool>"
//
// Copilot camelCase MCP names (serverKey-toolName) are not inferred here; codecs
// set ToolCall.MCP from dialect-specific structured metadata.
func NormalizeToolName(native string) (name string, mcp bool) {
	if strings.HasPrefix(native, "mcp__") || strings.HasPrefix(native, "MCP:") {
		return native, true
	}
	if n, ok := toolAliases[strings.ToLower(native)]; ok {
		return n, false
	}
	return native, false
}
