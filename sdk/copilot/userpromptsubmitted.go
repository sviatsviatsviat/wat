package copilot

import (
	"context"
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
	registerDecoder(EventUserPromptSubmitted, decodeAs[UserPromptSubmitted])
}

// UserPromptSubmitted registers an observe-only userPromptSubmitted handler.
func (c *Chain) UserPromptSubmitted(fn func(context.Context, Hook[UserPromptSubmitted]) error) *Chain {
	registerObserveHandler(c.registerOwner(), fn)
	return c
}
