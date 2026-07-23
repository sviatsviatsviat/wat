package beforereadfile

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/event"
)

// Event is the beforeReadFile hook event.
type Event struct {
	event.Envelope
	// FilePath is the file path being read.
	FilePath string `json:"file_path"`
	// Content is the file content.
	Content string `json:"content"`
	// Attachments are additional file attachments.
	Attachments []event.Attachment `json:"attachments"`
}

// EventName returns the canonical hook event name.
func (Event) EventName() string { return event.BeforeReadFile }

// register registers this hook event decoder on c.
func register(c *hookkit.Codec) {
	c.Register(event.BeforeReadFile, hookkit.EventDecoder[Event](c))
}
