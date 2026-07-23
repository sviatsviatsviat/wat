package cursor

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/config"
)

// File is the Cursor hooks.json shape.
type File = config.File

// Handler is a Cursor hook handler definition.
type Handler = config.Handler

const (
	// HandlerTypeCommand is the command handler type.
	HandlerTypeCommand = config.HandlerTypeCommand
	// HandlerTypePrompt is the prompt handler type.
	HandlerTypePrompt = config.HandlerTypePrompt
)

// ParseFlatCommand decodes native handler JSON and returns the shell command when type is empty or command.
func ParseFlatCommand(raw json.RawMessage) (string, bool) {
	return config.ParseFlatCommand(raw)
}

// MarshalFlatCommand encodes a command-type handler as native hooks.json JSON.
func MarshalFlatCommand(command string) (json.RawMessage, error) {
	return config.MarshalFlatCommand(command)
}
