package hookkit

import "encoding/json"

const flatCommandType = "command"

// FlatCommandHandler is the command field shape shared by copilot and cursor hooks.json handlers.
type FlatCommandHandler struct {
	// Type is the handler kind; empty or command for shell commands.
	Type string `json:"type,omitempty"`
	// Command is the shell command for command-type handlers.
	Command string `json:"command,omitempty"`
}

// ParseHandler decodes native handler JSON into T.
func ParseHandler[T any](raw json.RawMessage) (T, error) {
	var h T
	if len(raw) == 0 {
		return h, nil
	}
	err := json.Unmarshal(raw, &h)
	return h, err
}

// MarshalHandler encodes h as native handler JSON.
func MarshalHandler[T any](h T) (json.RawMessage, error) {
	return json.Marshal(h)
}

// Handlers encodes typed handlers as native handler JSON blobs.
func Handlers[T any](h ...T) ([]json.RawMessage, error) {
	out := make([]json.RawMessage, 0, len(h))
	for _, handler := range h {
		raw, err := MarshalHandler(handler)
		if err != nil {
			return nil, err
		}
		out = append(out, raw)
	}
	return out, nil
}

// ParseFlatCommand decodes copilot or cursor flat handler JSON and returns the command
// when type is empty or command.
func ParseFlatCommand(raw json.RawMessage) (string, bool) {
	if len(raw) == 0 || string(raw) == "null" {
		return "", false
	}
	h, err := ParseHandler[FlatCommandHandler](raw)
	if err != nil {
		return "", false
	}
	if h.Type != "" && h.Type != flatCommandType {
		return "", false
	}
	return h.Command, true
}

// MarshalFlatCommand encodes a command-type flat handler as native JSON.
func MarshalFlatCommand(command string) (json.RawMessage, error) {
	return MarshalHandler(FlatCommandHandler{Type: flatCommandType, Command: command})
}
