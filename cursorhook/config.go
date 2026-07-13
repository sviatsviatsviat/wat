package cursorhook

import "encoding/json"

// File is the Cursor hooks.json shape.
type File struct {
	// Version is the hooks file schema version.
	Version int `json:"version"`
	// Hooks maps event names to handler definitions.
	Hooks map[string][]json.RawMessage `json:"hooks"`
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

// ParseHandler decodes native handler JSON into a Handler.
func ParseHandler(raw json.RawMessage) (Handler, error) {
	var h Handler
	if len(raw) == 0 {
		return h, nil
	}
	err := json.Unmarshal(raw, &h)
	return h, err
}

// MarshalHandler encodes a Handler as native handler JSON.
func MarshalHandler(h Handler) (json.RawMessage, error) {
	return json.Marshal(h)
}

// Handlers encodes typed handlers as native handler JSON blobs.
func Handlers(h ...Handler) ([]json.RawMessage, error) {
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

// TimeoutSeconds returns the configured timeout in seconds.
func (h Handler) TimeoutSeconds() int {
	return h.Timeout
}
