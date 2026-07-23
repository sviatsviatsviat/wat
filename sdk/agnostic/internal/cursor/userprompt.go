package cursor

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/run"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

// RegisterUserPrompt registers fn on the Cursor BeforeSubmitPrompt chain.
func RegisterUserPrompt(fn model.UserPromptHandler) run.Hooks {
	if fn == nil {
		return nil
	}
	return sdkcursor.UseHooks().BeforeSubmitPrompt(func(ctx context.Context, hook sdkcursor.BeforeSubmitPrompt, _ sdkcursor.BeforeSubmitPromptResults) (sdkcursor.BeforeSubmitPromptOutput, error) {
		return nil, fn(ctx, *mapBeforeSubmitPrompt(hook))
	})
}

func mapBeforeSubmitPrompt(e sdkcursor.BeforeSubmitPrompt) *model.UserPromptEvent {
	return &model.UserPromptEvent{
		Envelope: envelope(e.Envelope, e.EventName()),
		Prompt:   e.Prompt,
	}
}
