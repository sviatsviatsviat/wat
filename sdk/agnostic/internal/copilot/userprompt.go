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
	sdkcopilot.OnUserPromptSubmitted(func(ctx context.Context, hook sdkcopilot.Hook[sdkcopilot.UserPromptSubmitted]) error {
		return fn(ctx, model.NewUserPromptHook(hook.Invocation(), mapUserPromptSubmitted(hook.Event, hook.Raw())))
	})
}

func mapUserPromptSubmitted(e sdkcopilot.UserPromptSubmitted, raw []byte) *model.UserPromptEvent {
	return &model.UserPromptEvent{
		Envelope: envelope(e, raw),
		Prompt:   e.Prompt,
	}
}
