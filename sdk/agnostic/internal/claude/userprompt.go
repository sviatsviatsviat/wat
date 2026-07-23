package claude

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/run"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
)

// RegisterUserPrompt registers fn on the Claude UserPromptSubmit chain.
func RegisterUserPrompt(fn model.UserPromptHandler) run.Hooks {
	if fn == nil {
		return nil
	}
	return sdkclaude.UseHooks().UserPromptSubmit(func(ctx context.Context, hook sdkclaude.UserPromptSubmit, _ sdkclaude.UserPromptSubmitResults) (sdkclaude.UserPromptSubmitOutput, error) {
		return nil, fn(ctx, *mapUserPromptSubmit(hook))
	})
}

func mapUserPromptSubmit(e sdkclaude.UserPromptSubmit) *model.UserPromptEvent {
	return &model.UserPromptEvent{
		Envelope: envelope(e.Envelope, e.EventName()),
		Prompt:   e.Prompt,
	}
}
