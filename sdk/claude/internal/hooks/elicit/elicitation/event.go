package elicitation

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
)

// Event is the Elicitation hook event.
type Event struct {
	event.Envelope
	// MCPServerName is the MCP server requesting input.
	MCPServerName string `json:"mcp_server_name"`
	// Message is the elicitation prompt shown to the user.
	Message string `json:"message"`
	// Mode is the elicitation mode when provided ("form" or "url").
	Mode string `json:"mode"`
	// URL is the browser authentication URL for url-mode elicitation.
	URL string `json:"url"`
	// ElicitationID is the elicitation request identifier when provided.
	ElicitationID string `json:"elicitation_id"`
	// RequestedSchema is the requested form input schema JSON when provided.
	RequestedSchema json.RawMessage `json:"requested_schema"`
}

// EventName returns the hook event name.
func (Event) EventName() string { return event.Elicitation }

// register registers this hook event decoder on c.
func register(c *hookkit.Codec) {
	c.Register(event.Elicitation, hookkit.EventDecoder[Event](c))
}
