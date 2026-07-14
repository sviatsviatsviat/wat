package copilot

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
)

// File is the GitHub Copilot hooks.json shape.
type File struct {
	// Version is the hooks file schema version.
	Version int `json:"version"`
	// Hooks maps event names to handler definitions.
	Hooks map[string][]json.RawMessage `json:"hooks"`
}

// Handler is a GitHub Copilot hook handler definition.
type Handler struct {
	// Type is the handler kind: command, http, or prompt.
	Type string `json:"type,omitempty"`
	// Bash is the bash command for command-type handlers.
	Bash string `json:"bash,omitempty"`
	// PowerShell is the PowerShell command for command-type handlers.
	PowerShell string `json:"powershell,omitempty"`
	// Command is the shell command for command-type handlers.
	Command string `json:"command,omitempty"`
	// URL is the HTTP endpoint for http-type handlers.
	URL string `json:"url,omitempty"`
	// Prompt is the prompt text for prompt handlers.
	Prompt string `json:"prompt,omitempty"`
	// Matcher is the tool or event matcher string.
	Matcher string `json:"matcher,omitempty"`
	// TimeoutSec is the handler timeout in seconds.
	TimeoutSec int `json:"timeoutSec,omitempty"`
	// Timeout is an alternate timeout field in seconds.
	Timeout int `json:"timeout,omitempty"`
}

// ParseHandler decodes native handler JSON into a Handler.
func ParseHandler(raw json.RawMessage) (Handler, error) {
	return hookkit.ParseHandler[Handler](raw)
}

// MarshalHandler encodes a Handler as native handler JSON.
func MarshalHandler(h Handler) (json.RawMessage, error) {
	return hookkit.MarshalHandler(h)
}

// Handlers encodes typed handlers as native handler JSON blobs.
func Handlers(h ...Handler) ([]json.RawMessage, error) {
	return hookkit.Handlers(h...)
}

// TimeoutSeconds returns the configured timeout in seconds.
func (h Handler) TimeoutSeconds() int {
	if h.TimeoutSec != 0 {
		return h.TimeoutSec
	}
	return h.Timeout
}

// EffectiveCommand returns the command string from command, bash, or powershell.
func (h Handler) EffectiveCommand() string {
	if h.Command != "" {
		return h.Command
	}
	if h.Bash != "" {
		return h.Bash
	}
	return h.PowerShell
}
