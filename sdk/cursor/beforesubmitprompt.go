package cursor

import (
	"context"
	"encoding/json"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// BeforeSubmitPrompt is the beforeSubmitPrompt hook event.
type BeforeSubmitPrompt struct {
	Envelope
	// Prompt is the user prompt about to be submitted.
	Prompt string `json:"prompt"`
	// Attachments are context attachments associated with the prompt.
	Attachments []Attachment `json:"attachments"`
}

// EventName returns the canonical hook event name.
func (BeforeSubmitPrompt) EventName() string { return EventBeforeSubmitPrompt }

// BeforeSubmitPromptOutput is the response for beforeSubmitPrompt events.
type BeforeSubmitPromptOutput struct {
	// Continue is false to block prompt submission.
	Continue *bool
	// UserMessage is shown to the user when blocking.
	UserMessage string
}

func (o BeforeSubmitPromptOutput) isZero() bool {
	return o.Continue == nil && o.UserMessage == ""
}

// BeforeSubmitPromptResults is the hook-scoped response builder supplied to Chain handlers by registration.
type BeforeSubmitPromptResults interface {
	// Block blocks prompt submission with a user-facing message.
	Block(userMessage string) BeforeSubmitPromptOutput
	isBeforeSubmitPromptResults()
}

type beforeSubmitPromptResults struct{}

func (beforeSubmitPromptResults) isBeforeSubmitPromptResults() {}

// Block blocks prompt submission with a user-facing message.
func (beforeSubmitPromptResults) Block(userMessage string) BeforeSubmitPromptOutput {
	cont := false
	return BeforeSubmitPromptOutput{Continue: &cont, UserMessage: userMessage}
}

func encodeBeforeSubmitPrompt(o BeforeSubmitPromptOutput) ([]byte, int, error) {
	out := map[string]any{}
	if o.Continue != nil {
		out["continue"] = *o.Continue
	}
	if o.UserMessage != "" {
		out["user_message"] = o.UserMessage
	}
	if len(out) == 0 {
		return nil, 0, nil
	}
	b, err := json.Marshal(out)
	return b, 0, err
}

func init() {
	registerDecoder(EventBeforeSubmitPrompt, decodeAs[BeforeSubmitPrompt])
}

// BeforeSubmitPrompt registers a beforeSubmitPrompt handler.
func (c *Chain) BeforeSubmitPrompt(fn func(context.Context, BeforeSubmitPromptHook, BeforeSubmitPromptResults) (BeforeSubmitPromptOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev BeforeSubmitPrompt) (BeforeSubmitPromptOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), beforeSubmitPromptResults{})
	})
	return &Chain{}
}
