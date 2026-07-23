package elicitation

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
)

// Event is the Elicitation hook event.
type Event struct {
	event.Envelope
	// ServerName is the MCP server name.
	ServerName string `json:"server_name"`
	// Message is the elicitation message.
	Message string `json:"message"`
	// Schema is the requested input schema JSON.
	Schema json.RawMessage `json:"requested_schema"`
}

// EventName returns the hook event name.
func (Event) EventName() string { return event.Elicitation }

// Register registers this hook event decoder on c.
func Register(c *hookkit.Codec) {
	c.Register(event.Elicitation, hookkit.EventDecoder[Event](c))
}
