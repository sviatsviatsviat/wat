package claude

import (
	"context"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

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
	codec.Register(EventUserPromptSubmit, hookkit.EventDecoder[UserPromptSubmit](codec))
}

// UserPromptSubmitOutput is the response for UserPromptSubmit events.
// Construct via UserPromptSubmitResults builders and With* methods.
// A nil value is a no-op.
type UserPromptSubmitOutput interface {
	run.Output
	isUserPromptSubmitOutput()
	// WithAdditionalContext injects model context.
	WithAdditionalContext(text string) UserPromptSubmitOutput
	// WithSessionTitle sets the session title.
	WithSessionTitle(title string) UserPromptSubmitOutput
	// WithSuppressOriginalPrompt suppresses the original prompt when true.
	WithSuppressOriginalPrompt(v bool) UserPromptSubmitOutput
	// WithContinue sets whether Claude should continue the session.
	WithContinue(v bool) UserPromptSubmitOutput
	// WithStopReason explains why the session was stopped.
	WithStopReason(reason string) UserPromptSubmitOutput
	// WithSuppressOutput suppresses hook stdout when true.
	WithSuppressOutput(v bool) UserPromptSubmitOutput
	// WithSystemMessage sets a user-visible system message.
	WithSystemMessage(msg string) UserPromptSubmitOutput
	// WithTerminalSequence sets an OSC terminal sequence (allowlisted).
	WithTerminalSequence(seq string) UserPromptSubmitOutput
}

type userPromptSubmitOutput struct {
	common
	block                  bool
	reason                 string
	additionalContext      string
	sessionTitle           string
	suppressOriginalPrompt bool
}

func (userPromptSubmitOutput) isUserPromptSubmitOutput() {}

// IsZero reports whether this hook response is empty.
func (o userPromptSubmitOutput) IsZero() bool {
	return o.common.IsZero() && !o.block && o.reason == "" &&
		o.additionalContext == "" && o.sessionTitle == "" && !o.suppressOriginalPrompt
}

// WithAdditionalContext injects model context.
func (o userPromptSubmitOutput) WithAdditionalContext(text string) UserPromptSubmitOutput {
	o.additionalContext = text
	return o
}

// WithSessionTitle sets the session title.
func (o userPromptSubmitOutput) WithSessionTitle(title string) UserPromptSubmitOutput {
	o.sessionTitle = title
	return o
}

// WithSuppressOriginalPrompt suppresses the original prompt when true.
func (o userPromptSubmitOutput) WithSuppressOriginalPrompt(v bool) UserPromptSubmitOutput {
	o.suppressOriginalPrompt = v
	return o
}

// WithContinue sets whether Claude should continue the session.
func (o userPromptSubmitOutput) WithContinue(v bool) UserPromptSubmitOutput {
	o.common = o.common.WithContinue(v)
	return o
}

// WithStopReason explains why the session was stopped.
func (o userPromptSubmitOutput) WithStopReason(reason string) UserPromptSubmitOutput {
	o.common = o.common.WithStopReason(reason)
	return o
}

// WithSuppressOutput suppresses hook stdout when true.
func (o userPromptSubmitOutput) WithSuppressOutput(v bool) UserPromptSubmitOutput {
	o.common = o.common.WithSuppressOutput(v)
	return o
}

// WithSystemMessage sets a user-visible system message.
func (o userPromptSubmitOutput) WithSystemMessage(msg string) UserPromptSubmitOutput {
	o.common = o.common.WithSystemMessage(msg)
	return o
}

// WithTerminalSequence sets an OSC terminal sequence (allowlisted).
func (o userPromptSubmitOutput) WithTerminalSequence(seq string) UserPromptSubmitOutput {
	o.common = o.common.WithTerminalSequence(seq)
	return o
}

// UserPromptSubmitResults is the hook-scoped response builder supplied to On* handlers by registration.
type UserPromptSubmitResults interface {
	// Context returns a non-blocking context-injection result.
	Context(text string) UserPromptSubmitOutput
	// Block returns a block result with an agent-facing reason.
	Block(reason string) UserPromptSubmitOutput
	// Noop returns an empty response (silent stdout). Prefer nil from handlers when not chaining With*.
	Noop() UserPromptSubmitOutput
	isUserPromptSubmitResults()
}

type userPromptSubmitResults struct{}

func (userPromptSubmitResults) isUserPromptSubmitResults() {}

// Context returns a non-blocking context-injection result.
func (userPromptSubmitResults) Context(text string) UserPromptSubmitOutput {
	return userPromptSubmitOutput{additionalContext: text}
}

// Block returns a block result with an agent-facing reason.
func (userPromptSubmitResults) Block(reason string) UserPromptSubmitOutput {
	return userPromptSubmitOutput{block: true, reason: reason}
}

// Noop returns an empty response (silent stdout).
func (userPromptSubmitResults) Noop() UserPromptSubmitOutput {
	return userPromptSubmitOutput{}
}

func (o userPromptSubmitOutput) encodeInto(top, hso map[string]any) {
	applyCommon(top, o.common)
	if o.block {
		top["decision"] = "block"
		if o.reason != "" {
			top["reason"] = o.reason
		}
	}
	if o.additionalContext != "" {
		hso["additionalContext"] = o.additionalContext
	}
	if o.sessionTitle != "" {
		hso["sessionTitle"] = o.sessionTitle
	}
	if o.suppressOriginalPrompt {
		hso["suppressOriginalPrompt"] = true
	}
}

// Encode renders this output as Claude Code stdout JSON.
func (o userPromptSubmitOutput) Encode() ([]byte, int, error) {
	return marshalHookOutput(EventUserPromptSubmit, o.encodeInto)
}

// UserPromptSubmit registers a UserPromptSubmit handler on the chain.
func (c *chain) UserPromptSubmit(fn func(context.Context, run.Hook[UserPromptSubmit], UserPromptSubmitResults) (UserPromptSubmitOutput, error)) *chain {
	if fn == nil {
		return c
	}
	c.reg.RegisterHandler(Dialect, run.Handler(func(ctx context.Context, hook run.Hook[UserPromptSubmit]) (UserPromptSubmitOutput, error) {
		return fn(ctx, hook, userPromptSubmitResults{})
	}))
	return c
}
