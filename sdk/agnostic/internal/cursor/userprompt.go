package cursor

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

// RegisterUserPrompt registers fn on the Cursor BeforeSubmitPrompt chain.
func RegisterUserPrompt(fn model.UserPromptHandler) {
	if fn == nil {
		return
	}
	new(sdkcursor.Chain).BeforeSubmitPrompt(func(ctx context.Context, hook sdkcursor.Hook[sdkcursor.BeforeSubmitPrompt], _ sdkcursor.BeforeSubmitPromptResults) (sdkcursor.BeforeSubmitPromptOutput, error) {
		return nil, fn(ctx, model.NewUserPromptHook(hook.Invocation(), mapBeforeSubmitPrompt(hook.Event, hook.Raw())))
	})
}

func mapBeforeSubmitPrompt(e sdkcursor.BeforeSubmitPrompt, raw []byte) *model.UserPromptEvent {
	return &model.UserPromptEvent{
		Envelope: envelope(e, raw),
		Prompt:   e.Prompt,
	}
}
