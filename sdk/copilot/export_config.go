package copilot

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/config"
)

// File is the GitHub Copilot hooks.json shape.
type File = config.File

// Handler is a GitHub Copilot hook handler definition.
type Handler = config.Handler

// ParseFlatCommand decodes native handler JSON and returns the shell command when type is empty or command.
func ParseFlatCommand(raw json.RawMessage) (string, bool) {
	return config.ParseFlatCommand(raw)
}

// MarshalFlatCommand encodes a command-type handler as native hooks.json JSON.
func MarshalFlatCommand(command string) (json.RawMessage, error) {
	return config.MarshalFlatCommand(command)
}
