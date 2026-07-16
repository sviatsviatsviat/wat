package tools

import "strings"

// MCPTool is a parsed Claude MCP tool reference (mcp__server__tool).
type MCPTool struct {
	// Server is the MCP server name.
	Server string
	// Tool is the MCP tool name on that server.
	Tool string
}

// AsMCPTool returns the MCP server and tool parts when this payload is an MCP tool.
func (in Input) AsMCPTool() (MCPTool, bool) {
	server, tool, ok := parseMCPName(in.Name())
	if !ok {
		return MCPTool{}, false
	}
	return MCPTool{Server: server, Tool: tool}, true
}

func parseMCPName(name string) (server, tool string, ok bool) {
	const prefix = "mcp__"
	if !strings.HasPrefix(name, prefix) {
		return "", "", false
	}
	rest := strings.TrimPrefix(name, prefix)
	parts := strings.SplitN(rest, "__", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", false
	}
	return parts[0], parts[1], true
}
