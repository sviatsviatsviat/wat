package elicitationresult

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
)

// Event is the ElicitationResult hook event.
type Event struct {
	event.Envelope
	// MCPServerName is the MCP server that requested input.
	MCPServerName string `json:"mcp_server_name"`
	// Action is the user action (accept, decline, cancel).
	Action string `json:"action"`
	// Mode is the elicitation mode when provided ("form" or "url").
	Mode string `json:"mode"`
	// ElicitationID is the elicitation request identifier when provided.
	ElicitationID string `json:"elicitation_id"`
	// Content is the response content JSON when provided.
	Content json.RawMessage `json:"content"`
}

// EventName returns the hook event name.
func (Event) EventName() string { return event.ElicitationResult }

// register registers this hook event decoder on c.
func register(c *hookkit.Codec) {
	c.Register(event.ElicitationResult, hookkit.EventDecoder[Event](c))
}
