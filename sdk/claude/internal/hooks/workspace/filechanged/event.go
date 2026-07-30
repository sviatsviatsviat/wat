package filechanged

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
)

// Event is the FileChanged hook event.
type Event struct {
	event.Envelope
	// FilePath is the absolute path of the changed file.
	FilePath string `json:"file_path"`
	// Change is what happened to the file ("change", "add", or "unlink").
	Change string `json:"event"`
}

// EventName returns the hook event name.
func (Event) EventName() string { return event.FileChanged }

// register registers this hook event decoder on c.
func register(c *hookkit.Codec) {
	c.Register(event.FileChanged, hookkit.EventDecoder[Event](c))
}
