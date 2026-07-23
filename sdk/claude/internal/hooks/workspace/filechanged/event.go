package filechanged

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
)

// Event is the FileChanged hook event.
type Event struct {
	event.Envelope
	// FilePath is the changed file path.
	FilePath string `json:"file_path"`
}

// EventName returns the hook event name.
func (Event) EventName() string { return event.FileChanged }

// register registers this hook event decoder on c.
func register(c *hookkit.Codec) {
	c.Register(event.FileChanged, hookkit.EventDecoder[Event](c))
}
