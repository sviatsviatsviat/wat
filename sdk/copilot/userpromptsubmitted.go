package copilot

import (
	"context"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// UserPromptSubmitted is the userPromptSubmitted hook event.
type UserPromptSubmitted struct {
	Envelope
	// Prompt is the submitted user prompt text.
	Prompt string `json:"prompt"`
}

// EventName returns the canonical hook event name.
func (UserPromptSubmitted) EventName() string { return EventUserPromptSubmitted }

func init() {
	codec.Register(EventUserPromptSubmitted, hookkit.EventDecoder[UserPromptSubmitted](codec))
}

// OnUserPromptSubmitted registers an observe-only userPromptSubmitted handler.
func OnUserPromptSubmitted(fn func(context.Context, run.Hook[UserPromptSubmitted]) error) *chain {
	return (&chain{}).UserPromptSubmitted(fn)
}

// UserPromptSubmitted registers another UserPromptSubmitted handler on the chain.
func (c *chain) UserPromptSubmitted(fn func(context.Context, run.Hook[UserPromptSubmitted]) error) *chain {
	registerObserveHandler(fn)
	return c
}
