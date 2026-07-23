package sessionstart

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
)

// Event is the SessionStart hook event.
type Event struct {
	event.Envelope
	// Source is the session start source (startup, resume, clear, compact).
	Source string `json:"source"`
	// Model is the model name when provided.
	Model string `json:"model,omitempty"`
	// SessionTitle is the session title when provided.
	SessionTitle string `json:"session_title,omitempty"`
}

// EventName returns the hook event name.
func (Event) EventName() string { return event.SessionStart }

// Register registers the SessionStart decoder on c.
func Register(c *hookkit.Codec) {
	c.Register(event.SessionStart, hookkit.EventDecoder[Event](c))
}
