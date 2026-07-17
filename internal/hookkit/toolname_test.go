package hookkit

import "testing"

func TestNormalizeToolName(t *testing.T) {
	tests := []struct {
		name     string
		native   string
		wantName string
		wantMCP  bool
	}{
		// Claude Code
		{name: "claude_bash", native: "Bash", wantName: ToolBash},
		{name: "claude_edit", native: "Edit", wantName: ToolEdit},
		{name: "claude_write", native: "Write", wantName: ToolWrite},
		{name: "claude_read", native: "Read", wantName: ToolRead},
		{name: "claude_agent", native: "Agent", wantName: ToolTask},
		{name: "claude_glob", native: "Glob", wantName: ToolGlob},
		{name: "claude_mcp", native: "mcp__github__create_issue", wantName: "mcp__github__create_issue", wantMCP: true},
		// Copilot
		{name: "copilot_bash", native: "bash", wantName: ToolBash},
		{name: "copilot_powershell", native: "powershell", wantName: ToolBash},
		{name: "copilot_create", native: "create", wantName: ToolWrite},
		{name: "copilot_view", native: "view", wantName: ToolRead},
		{name: "copilot_web_fetch", native: "web_fetch", wantName: ToolWebFetch},
		{name: "copilot_mcp_name_passthrough", native: "my-server-list_items", wantName: "my-server-list_items"},
		// Cursor
		{name: "cursor_shell", native: "Shell", wantName: ToolBash},
		{name: "cursor_delete", native: "delete", wantName: ToolDelete},
		{name: "cursor_mcp", native: "MCP:browser_navigate", wantName: "MCP:browser_navigate", wantMCP: true},
		// Passthrough
		{name: "unknown_tool", native: "CustomTool", wantName: "CustomTool"},
		{name: "hyphenated_non_mcp", native: "Custom-Tool", wantName: "Custom-Tool"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotName, gotMCP := NormalizeToolName(tt.native)
			if gotName != tt.wantName || gotMCP != tt.wantMCP {
				t.Fatalf("NormalizeToolName(%q) = (%q, %v), want (%q, %v)",
					tt.native, gotName, gotMCP, tt.wantName, tt.wantMCP)
			}
		})
	}
}
