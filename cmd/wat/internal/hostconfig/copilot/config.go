package copilot

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/internal/hookconfig"
)

// File is the GitHub Copilot hooks.json shape.
type File struct {
	// Version is the hooks file schema version.
	Version int `json:"version"`
	// DisableAllHooks disables every hook in this file when true.
	DisableAllHooks bool `json:"disableAllHooks,omitempty"`
	// Hooks maps event names to handler definitions.
	Hooks map[string][]json.RawMessage `json:"hooks"`
}

// HooksMap returns hook entries keyed by event name.
func (f File) HooksMap() map[string][]json.RawMessage {
	return f.Hooks
}

// ParseFlatCommand decodes native handler JSON and returns the shell command when type is empty or command.
func ParseFlatCommand(raw json.RawMessage) (string, bool) {
	return hookconfig.ParseFlatCommand(raw)
}

// MarshalFlatCommand encodes a command-type handler as native hooks.json JSON.
func MarshalFlatCommand(command string) (json.RawMessage, error) {
	return hookconfig.MarshalFlatCommand(command)
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
	// Cwd is the working directory for command-type handlers.
	Cwd string `json:"cwd,omitempty"`
	// Env holds additional environment variables for command-type handlers.
	Env map[string]string `json:"env,omitempty"`
	// Headers are HTTP headers for http-type handlers.
	Headers map[string]string `json:"headers,omitempty"`
	// AllowedEnvVars lists env vars that may be expanded inside headers.
	AllowedEnvVars []string `json:"allowedEnvVars,omitempty"`
	// Matcher is the tool or event matcher string.
	Matcher string `json:"matcher,omitempty"`
	// TimeoutSec is the handler timeout in seconds.
	TimeoutSec int `json:"timeoutSec,omitempty"`
	// Timeout is an alternate timeout field in seconds.
	Timeout int `json:"timeout,omitempty"`
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
