package copilot

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
)

// RegisterUserPrompt registers fn on the Copilot UserPromptSubmitted chain.
func RegisterUserPrompt(fn model.UserPromptHandler) {
	if fn == nil {
		return
	}
	sdkcopilot.UseHooks().UserPromptSubmitted(func(ctx context.Context, hook sdkcopilot.UserPromptSubmitted) error {
		return fn(ctx, *mapUserPromptSubmitted(hook))
	})
}

func mapUserPromptSubmitted(e sdkcopilot.UserPromptSubmitted) *model.UserPromptEvent {
	return &model.UserPromptEvent{
		Envelope: envelope(e.Envelope, e.EventName()),
		Prompt:   e.Prompt,
	}
}
