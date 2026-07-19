package copilot

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/run"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
)

// RegisterUserPrompt registers fn on the Copilot UserPromptSubmitted chain.
func RegisterUserPrompt(registry *run.Registry, fn model.UserPromptHandler) {
	if fn == nil {
		return
	}
	sdkcopilot.UseHooks(registry).UserPromptSubmitted(func(ctx context.Context, hook run.Hook[sdkcopilot.UserPromptSubmitted]) error {
		return fn(ctx, model.NewUserPromptHook(hook.Invocation(), mapUserPromptSubmitted(hook.Event)))
	})
}

func mapUserPromptSubmitted(e sdkcopilot.UserPromptSubmitted) *model.UserPromptEvent {
	return &model.UserPromptEvent{
		Envelope: envelope(e.Envelope, e.EventName()),
		Prompt:   e.Prompt,
	}
}
