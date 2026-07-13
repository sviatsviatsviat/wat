package cursor

import "os"

// Environment variable names exported by Cursor for hook processes.
const (
	// EnvProjectDir is the Cursor project directory path.
	EnvProjectDir = "CURSOR_PROJECT_DIR"
	// EnvVersion is the Cursor application version.
	EnvVersion = "CURSOR_VERSION"
	// EnvClaudeProjectDir is also set by Cursor for Claude Code compatibility.
	// Prefer payload fields conversation_id and cursor_version, or an explicit
	// agent override, when distinguishing Cursor from Claude Code.
	EnvClaudeProjectDir = "CLAUDE_PROJECT_DIR"
)

// ProjectDir returns CURSOR_PROJECT_DIR from the environment.
func ProjectDir() string {
	return os.Getenv(EnvProjectDir)
}

// Version returns CURSOR_VERSION from the environment.
func Version() string {
	return os.Getenv(EnvVersion)
}

// ClaudeProjectDir returns CLAUDE_PROJECT_DIR from the environment.
//
// Cursor exports this variable for Claude Code compatibility. When both Cursor
// and Claude tooling may be present, prefer hook payload fields or an explicit
// agent selection over env-only detection.
func ClaudeProjectDir() string {
	return os.Getenv(EnvClaudeProjectDir)
}
