package tools

import "strings"

// IsShellToolName reports whether name is a shell execution tool.
func IsShellToolName(name string) bool {
	switch strings.ToLower(name) {
	case ToolBash, ToolPowerShell, ToolShell:
		return true
	default:
		return false
	}
}
