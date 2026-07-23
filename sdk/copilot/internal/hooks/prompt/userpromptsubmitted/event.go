package userpromptsubmitted

import (
	"context"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/runtime"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// Event is the UserPromptSubmitted hook event.
type Event struct {
	event.Envelope
	// Prompt is the submitted user prompt text.
	Prompt string `json:"prompt"`
}

// EventName returns the canonical hook event name.
func (Event) EventName() string { return event.UserPromptSubmitted }

// Register registers this hook event decoder on c.
func Register(c *hookkit.Codec) {
	c.Register(event.UserPromptSubmitted, hookkit.EventDecoder[Event](c))
}

// RegisterHandler registers a UserPromptSubmitted observe handler on reg.
func RegisterHandler(reg *run.Registry, fn func(context.Context, run.Hook[Event]) error) {
	if fn == nil {
		return
	}
	reg.RegisterObserveHandler(runtime.Dialect, run.ObserveHandler(fn))
}
