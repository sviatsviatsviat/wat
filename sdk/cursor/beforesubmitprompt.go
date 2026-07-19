package cursor

import (
	"context"
	"encoding/json"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

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
// Construct via BeforeSubmitPromptResults builders and With* methods. A nil value is a no-op.
type BeforeSubmitPromptOutput interface {
	Output
	isBeforeSubmitPromptOutput()
	// WithContinue sets whether prompt submission should continue.
	WithContinue(v bool) BeforeSubmitPromptOutput
	// WithUserMessage sets a user-facing message when blocking.
	WithUserMessage(msg string) BeforeSubmitPromptOutput
}

type beforeSubmitPromptOutput struct {
	cont        *bool
	userMessage string
}

func (beforeSubmitPromptOutput) isCursorOutput() {}

func (beforeSubmitPromptOutput) isBeforeSubmitPromptOutput() {}

func (o beforeSubmitPromptOutput) isZero() bool {
	return o.cont == nil && o.userMessage == ""
}

// WithContinue sets whether prompt submission should continue.
func (o beforeSubmitPromptOutput) WithContinue(v bool) BeforeSubmitPromptOutput {
	o.cont = &v
	return o
}

// WithUserMessage sets a user-facing message when blocking.
func (o beforeSubmitPromptOutput) WithUserMessage(msg string) BeforeSubmitPromptOutput {
	o.userMessage = msg
	return o
}

// BeforeSubmitPromptResults is the hook-scoped response builder supplied to On* handlers by registration.
type BeforeSubmitPromptResults interface {
	// Block blocks prompt submission with a user-facing message.
	Block(userMessage string) BeforeSubmitPromptOutput
	// Noop returns an empty response (silent stdout). Prefer nil from handlers when not chaining With*.
	Noop() BeforeSubmitPromptOutput
	isBeforeSubmitPromptResults()
}

type beforeSubmitPromptResults struct{}

func (beforeSubmitPromptResults) isBeforeSubmitPromptResults() {}

// Block blocks prompt submission with a user-facing message.
func (beforeSubmitPromptResults) Block(userMessage string) BeforeSubmitPromptOutput {
	cont := false
	return beforeSubmitPromptOutput{cont: &cont, userMessage: userMessage}
}

// Noop returns an empty response (silent stdout).
func (beforeSubmitPromptResults) Noop() BeforeSubmitPromptOutput {
	return beforeSubmitPromptOutput{}
}

func (beforeSubmitPromptOutput) allowedEvents() []string {
	return []string{EventBeforeSubmitPrompt}
}

func (o beforeSubmitPromptOutput) encode(eventName string) ([]byte, int, error) {
	_ = eventName
	out := map[string]any{}
	if o.cont != nil {
		out["continue"] = *o.cont
	}
	if o.userMessage != "" {
		out["user_message"] = o.userMessage
	}
	if len(out) == 0 {
		return nil, 0, nil
	}
	b, err := json.Marshal(out)
	return b, 0, err
}

func init() {
	codec.Register(EventBeforeSubmitPrompt, hookkit.EventDecoder[BeforeSubmitPrompt](codec))
}

// OnBeforeSubmitPrompt registers a beforeSubmitPrompt handler.
func OnBeforeSubmitPrompt(fn func(context.Context, run.Hook[BeforeSubmitPrompt], BeforeSubmitPromptResults) (BeforeSubmitPromptOutput, error)) *chain {
	return (&chain{}).BeforeSubmitPrompt(fn)
}

// BeforeSubmitPrompt registers another BeforeSubmitPrompt handler on the chain.
func (c *chain) BeforeSubmitPrompt(fn func(context.Context, run.Hook[BeforeSubmitPrompt], BeforeSubmitPromptResults) (BeforeSubmitPromptOutput, error)) *chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev BeforeSubmitPrompt) (BeforeSubmitPromptOutput, error) {
		return fn(ctx, run.NewHook(run.InvocationFrom(ctx), ev), beforeSubmitPromptResults{})
	})
	return c
}
