package claude

import (
	"encoding/json"
)

// Elicitation is the Elicitation hook event.
type Elicitation struct {
	Envelope
	// ServerName is the MCP server name.
	ServerName string `json:"server_name"`
	// Message is the elicitation message.
	Message string `json:"message"`
	// Schema is the requested input schema JSON.
	Schema json.RawMessage `json:"requested_schema"`
}

// EventName returns the hook event name.
func (Elicitation) EventName() string { return EventElicitation }

func init() {
	registerDecoder(EventElicitation, decodeAs[Elicitation])
}

// ElicitationOutput is the response for Elicitation events.
type ElicitationOutput struct {
	Common
	// Action is the elicitation action.
	Action string
	// Content is the elicitation response content.
	Content map[string]any
}

func (o ElicitationOutput) isZero() bool {
	return o.Common.isZero() && o.Action == "" && o.Content == nil
}
