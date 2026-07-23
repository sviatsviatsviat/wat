package beforetabfileread

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/event"
)

// Event is the beforeTabFileRead hook event.
type Event struct {
	event.Envelope
	// FilePath is the file path being read.
	FilePath string `json:"file_path"`
	// Content is the file content.
	Content string `json:"content"`
}

// EventName returns the canonical hook event name.
func (Event) EventName() string { return event.BeforeTabFileRead }

// Register registers this hook event decoder on c.
func Register(c *hookkit.Codec) {
	c.Register(event.BeforeTabFileRead, hookkit.EventDecoder[Event](c))
}
