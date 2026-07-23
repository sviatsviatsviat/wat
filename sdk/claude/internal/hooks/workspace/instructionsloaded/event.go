package instructionsloaded

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
)

// Event is the InstructionsLoaded hook event.
type Event struct {
	event.Envelope
	// FilePath is the loaded instruction file path.
	FilePath string `json:"file_path"`
	// MemoryType is the memory type (User, Project, Local, Managed).
	MemoryType string `json:"memory_type"`
	// LoadReason is why the file was loaded.
	LoadReason string `json:"load_reason"`
	// Globs are glob patterns when applicable.
	Globs []string `json:"globs,omitempty"`
	// TriggerFilePath is the file that triggered loading.
	TriggerFilePath string `json:"trigger_file_path,omitempty"`
	// ParentFilePath is the parent file path when nested.
	ParentFilePath string `json:"parent_file_path,omitempty"`
}

// EventName returns the hook event name.
func (Event) EventName() string { return event.InstructionsLoaded }

// Register registers this hook event decoder on c.
func Register(c *hookkit.Codec) {
	c.Register(event.InstructionsLoaded, hookkit.EventDecoder[Event](c))
}
