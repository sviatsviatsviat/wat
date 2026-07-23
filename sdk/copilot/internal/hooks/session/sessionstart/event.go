package sessionstart

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/event"
)

// Event is the SessionStart hook event.
type Event struct {
	event.Envelope
	// Source is the session start source.
	Source string `json:"source"`
	// InitialPromptValue is the initial prompt.
	InitialPromptValue string `json:"initial_prompt"`
}

// EventName returns the canonical hook event name.
func (Event) EventName() string { return event.SessionStart }

// InitialPrompt returns the initial prompt.
func (e Event) InitialPrompt() string {
	return e.InitialPromptValue
}

// Register registers this hook event decoder on c.
func Register(c *hookkit.Codec) {
	c.Register(event.SessionStart, hookkit.EventDecoder[Event](c))
}
