package paths

import (
	"path/filepath"

	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

// AgentConfig is one agent dialect and its well-known config path.
type AgentConfig struct {
	// Agent is the dialect name (claude, copilot, or cursor).
	Agent string
	// Path is the absolute or project-relative config file path.
	Path string
}

// ConfigPath returns the well-known hook config file path for dialect under
// projectRoot. Unknown dialects return "".
func ConfigPath(dialect, projectRoot string) string {
	switch dialect {
	case sdkclaude.Dialect:
		return filepath.Join(projectRoot, ".claude", "settings.json")
	case sdkcopilot.Dialect:
		return filepath.Join(projectRoot, ".github", "hooks", "wat.json")
	case sdkcursor.Dialect:
		return filepath.Join(projectRoot, ".cursor", "hooks.json")
	default:
		return ""
	}
}

// All returns config paths for every supported agent dialect under projectRoot.
func All(projectRoot string) []AgentConfig {
	return []AgentConfig{
		{Agent: sdkclaude.Dialect, Path: ConfigPath(sdkclaude.Dialect, projectRoot)},
		{Agent: sdkcopilot.Dialect, Path: ConfigPath(sdkcopilot.Dialect, projectRoot)},
		{Agent: sdkcursor.Dialect, Path: ConfigPath(sdkcursor.Dialect, projectRoot)},
	}
}
