package cursor

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/run"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

// RegisterUserPrompt registers fn on the Cursor BeforeSubmitPrompt chain.
func RegisterUserPrompt(r *run.Registry, fn model.UserPromptHandler) {
	if fn == nil {
		return
	}
	sdkcursor.UseHooks(r).BeforeSubmitPrompt(func(ctx context.Context, hook run.Hook[sdkcursor.BeforeSubmitPrompt], _ sdkcursor.BeforeSubmitPromptResults) (sdkcursor.BeforeSubmitPromptOutput, error) {
		return nil, fn(ctx, model.NewUserPromptHook(hook.Invocation(), mapBeforeSubmitPrompt(hook.Event)))
	})
}

func mapBeforeSubmitPrompt(e sdkcursor.BeforeSubmitPrompt) *model.UserPromptEvent {
	return &model.UserPromptEvent{
		Envelope: envelope(e.Envelope, e.EventName()),
		Prompt:   e.Prompt,
	}
}
