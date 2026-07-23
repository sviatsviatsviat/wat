package userpromptexpansion

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
)

// Event is the UserPromptExpansion hook event.
type Event struct {
	event.Envelope
	// ExpansionType is the expansion kind (slash_command, mcp_prompt).
	ExpansionType string `json:"expansion_type"`
	// CommandName is the slash command name.
	CommandName string `json:"command_name"`
	// CommandArgs is the slash command arguments.
	CommandArgs string `json:"command_args"`
	// CommandSource is the command source.
	CommandSource string `json:"command_source"`
	// Prompt is the expanded prompt text.
	Prompt string `json:"prompt"`
}

// EventName returns the hook event name.
func (Event) EventName() string { return event.UserPromptExpansion }

// register registers this hook event decoder on c.
func register(c *hookkit.Codec) {
	c.Register(event.UserPromptExpansion, hookkit.EventDecoder[Event](c))
}
