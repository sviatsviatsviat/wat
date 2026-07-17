package model

import (
	"context"
	"encoding/json"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// UserPromptEvent is the normalized view of a UserPrompt hook invocation.
type UserPromptEvent struct {
	Envelope
	Prompt string
}

// UserPromptHook is the handler context for portable UserPrompt events.
type UserPromptHook struct {
	UserPromptEvent
	inv run.Invocation
}

// NewUserPromptHook wraps ev with serve-time invocation settings.
func NewUserPromptHook(inv run.Invocation, ev *UserPromptEvent) UserPromptHook {
	h := UserPromptHook{inv: inv}
	if ev != nil {
		h.UserPromptEvent = *ev
	}
	return h
}

// Invocation returns serve-time settings for this hook invocation.
func (h UserPromptHook) Invocation() run.Invocation { return h.inv }

// Raw returns the untouched native JSON payload.
func (h UserPromptHook) Raw() json.RawMessage { return h.UserPromptEvent.Raw }

// UserPromptHandler handles observe-only UserPrompt events.
type UserPromptHandler func(ctx context.Context, hook UserPromptHook) error
