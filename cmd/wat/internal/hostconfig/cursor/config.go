package cursor

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/internal/hookconfig"
)

// File is the Cursor hooks.json shape.
type File struct {
	// Version is the hooks file schema version.
	Version int `json:"version"`
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

// Handler type constants for Cursor hooks.json.
const (
	// HandlerTypeCommand is the command handler type.
	HandlerTypeCommand = "command"
	// HandlerTypePrompt is the prompt handler type.
	HandlerTypePrompt = "prompt"
)

// Handler is a Cursor hook handler definition.
type Handler struct {
	// Command is the shell command for command-type handlers.
	Command string `json:"command,omitempty"`
	// Type is the handler kind: command or prompt.
	Type string `json:"type,omitempty"`
	// Prompt is the prompt text for prompt handlers.
	Prompt string `json:"prompt,omitempty"`
	// Matcher is the tool or event matcher string.
	Matcher string `json:"matcher,omitempty"`
	// Timeout is the handler timeout in seconds.
	Timeout int `json:"timeout,omitempty"`
	// LoopLimit is the maximum stop-loop iterations when set.
	LoopLimit int `json:"loop_limit,omitempty"`
	// FailClosed makes handler errors fail-closed when set.
	FailClosed bool `json:"failClosed,omitempty"`
}

// TimeoutSeconds returns the configured timeout in seconds.
func (h Handler) TimeoutSeconds() int {
	return h.Timeout
}
