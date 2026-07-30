package configchange

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
)

// Event is the ConfigChange hook event.
type Event struct {
	event.Envelope
	// Source is the config source that changed.
	Source string `json:"source"`
	// FilePath is the path of the changed configuration file when provided.
	FilePath string `json:"file_path"`
}

// EventName returns the hook event name.
func (Event) EventName() string { return event.ConfigChange }

// register registers this hook event decoder on c.
func register(c *hookkit.Codec) {
	c.Register(event.ConfigChange, hookkit.EventDecoder[Event](c))
}
