package claude

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
)

// RegisterUserPrompt registers fn on the Claude UserPromptSubmit chain.
func RegisterUserPrompt(fn model.UserPromptHandler) {
	if fn == nil {
		return
	}
	new(sdkclaude.Chain).UserPromptSubmit(func(ctx context.Context, hook sdkclaude.Hook[sdkclaude.UserPromptSubmit], _ sdkclaude.UserPromptSubmitResults) (sdkclaude.UserPromptSubmitOutput, error) {
		return nil, fn(ctx, model.NewUserPromptHook(hook.Invocation(), mapUserPromptSubmit(hook.Event, hook.Raw())))
	})
}

func mapUserPromptSubmit(e sdkclaude.UserPromptSubmit, raw []byte) *model.UserPromptEvent {
	return &model.UserPromptEvent{
		Envelope: envelope(e, raw),
		Prompt:   e.Prompt,
	}
}
