package model

import (
	"context"
)

// UserPromptEvent is the normalized view of a UserPrompt hook invocation.
type UserPromptEvent struct {
	Envelope
	Prompt string
}

// UserPromptHandler handles observe-only UserPrompt events.
type UserPromptHandler func(ctx context.Context, event UserPromptEvent) error
