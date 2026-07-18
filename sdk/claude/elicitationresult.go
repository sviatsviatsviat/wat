package claude

import (
	"encoding/json"
)

// ElicitationResult is the ElicitationResult hook event.
type ElicitationResult struct {
	Envelope
	// ServerName is the MCP server name.
	ServerName string `json:"server_name"`
	// Action is the user action (accept, decline, cancel).
	Action string `json:"action"`
	// Content is the response content JSON.
	Content json.RawMessage `json:"content"`
}

// EventName returns the hook event name.
func (ElicitationResult) EventName() string { return EventElicitationResult }

func init() {
	registerDecoder(EventElicitationResult, decodeAs[ElicitationResult])
}
