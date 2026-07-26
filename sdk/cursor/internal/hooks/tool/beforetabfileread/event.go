package beforetabfileread

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/event"
)

// Event is the beforeTabFileRead hook event.
//
// Tab (inline completions) only; Cursor does not run this hook in cloud agents.
// Input is file_path and content (no attachments). Output is permission
// allow|deny only — see [Results].
type Event struct {
	event.Envelope
	// FilePath is the file path being read.
	FilePath string `json:"file_path"`
	// Content is the file content.
	Content string `json:"content"`
}

// EventName returns the canonical hook event name.
func (Event) EventName() string { return event.BeforeTabFileRead }

// register registers this hook event decoder on c.
func register(c *hookkit.Codec) {
	c.Register(event.BeforeTabFileRead, hookkit.EventDecoder[Event](c))
}
