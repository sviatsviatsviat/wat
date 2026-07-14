package claude

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// UserPromptSubmit is the UserPromptSubmit hook event.
type UserPromptSubmit struct {
	Envelope
	// Prompt is the submitted user prompt text.
	Prompt string `json:"prompt"`
}

// EventName returns the hook event name.
func (UserPromptSubmit) EventName() string { return EventUserPromptSubmit }

func init() {
	registerDecoder(EventUserPromptSubmit, decodeAs[UserPromptSubmit])
}

// UserPromptSubmitOutput is the response for UserPromptSubmit events.
type UserPromptSubmitOutput struct {
	Common
	// Block rejects the submitted prompt when true.
	Block bool
	// Reason is the block reason.
	Reason string
	// AdditionalContext injects model context.
	AdditionalContext string
	// SessionTitle sets the session title.
	SessionTitle string
	// SuppressOriginalPrompt suppresses the original prompt when true.
	SuppressOriginalPrompt bool
}

func (o UserPromptSubmitOutput) isZero() bool {
	return o.Common.isZero() && !o.Block && o.Reason == "" &&
		o.AdditionalContext == "" && o.SessionTitle == "" && !o.SuppressOriginalPrompt
}

// UserPromptSubmitResults is the hook-scoped response builder supplied to Chain handlers by registration.
type UserPromptSubmitResults interface {
	// Block returns a block result with an agent-facing reason.
	Block(reason string) UserPromptSubmitOutput
	isUserPromptSubmitResults()
}

type userPromptSubmitResults struct{}

func (userPromptSubmitResults) isUserPromptSubmitResults() {}

// Block returns a block result with an agent-facing reason.
func (userPromptSubmitResults) Block(reason string) UserPromptSubmitOutput {
	return UserPromptSubmitOutput{Block: true, Reason: reason}
}

// UserPromptSubmit registers a UserPromptSubmit handler.
func (c *Chain) UserPromptSubmit(fn func(context.Context, UserPromptSubmitHook, UserPromptSubmitResults) (UserPromptSubmitOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev UserPromptSubmit) (UserPromptSubmitOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), userPromptSubmitResults{})
	})
	return &Chain{}
}
