package dialect

import (
	"fmt"
	"strings"

	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

// Parse parses a dialect name from a CLI flag or config value.
// It accepts "claude", "claude-code", "claudecode", "copilot", "github-copilot",
// "gh", and "cursor". Unknown or empty strings return "".
func Parse(s string) string {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "claude", "claude-code", "claudecode":
		return sdkclaude.Dialect
	case "copilot", "github-copilot", "gh":
		return sdkcopilot.Dialect
	case "cursor":
		return sdkcursor.Dialect
	default:
		return ""
	}
}

// Validate reports whether agent is empty or a known dialect name.
func Validate(agent string) error {
	if agent == "" {
		return nil
	}
	if Parse(agent) == "" {
		return fmt.Errorf("unknown agent dialect %q (want claude, copilot, or cursor)", agent)
	}
	return nil
}
