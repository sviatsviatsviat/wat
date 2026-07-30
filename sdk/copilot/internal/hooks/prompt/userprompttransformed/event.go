package userprompttransformed

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/event"
)

// Event is the UserPromptTransformed hook event.
//
// It fires after the runtime transforms a submitted prompt into model-facing
// content, just before that content is emitted and persisted to session
// history. Mutation-only: handlers may rewrite the model-facing content but
// cannot block or handle the turn.
type Event struct {
	event.Envelope
	// Prompt is the user prompt after UserPromptSubmitted hooks have run.
	Prompt string `json:"prompt"`
	// TransformedPrompt is the runtime-transformed content the model will receive.
	TransformedPrompt string `json:"transformed_prompt"`
}

// EventName returns the canonical hook event name.
func (Event) EventName() string { return event.UserPromptTransformed }

// register registers this hook event decoder on c.
func register(c *hookkit.Codec) {
	c.Register(event.UserPromptTransformed, hookkit.EventDecoder[Event](c))
}
