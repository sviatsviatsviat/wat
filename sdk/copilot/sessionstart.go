package copilot

import (
	"context"
	"encoding/json"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// SessionStart is the sessionStart hook event.
type SessionStart struct {
	Envelope
	// Source is the session start source.
	Source string `json:"source"`
	// InitialPrompt is the initial prompt (camelCase).
	InitialPromptCamel string `json:"initialPrompt"`
	// InitialPromptSnake is the initial prompt (VS Code).
	InitialPromptSnake string `json:"initial_prompt"`
}

// EventName returns the canonical hook event name.
func (SessionStart) EventName() string { return EventSessionStart }

// InitialPrompt returns the initial prompt from either wire format.
func (e SessionStart) InitialPrompt() string {
	if e.InitialPromptSnake != "" {
		return e.InitialPromptSnake
	}
	return e.InitialPromptCamel
}

// SessionStartOutput is the response for sessionStart events.
// Construct via SessionStartResults builders. A nil value is a no-op.
type SessionStartOutput interface {
	isSessionStartOutput()
}

type sessionStartOutput struct {
	additionalContext string
}

func (sessionStartOutput) isSessionStartOutput() {}

func (o sessionStartOutput) isZero() bool {
	return o.additionalContext == ""
}

// SessionStartResults is the hook-scoped response builder supplied to Chain handlers by registration.
type SessionStartResults interface {
	// Context returns a context-injection-only SessionStart result.
	Context(text string) SessionStartOutput
	// Noop returns an empty response (silent stdout). Prefer nil from handlers when not chaining With*.
	Noop() SessionStartOutput
	isSessionStartResults()
}

type sessionStartResults struct{}

func (sessionStartResults) isSessionStartResults() {}

// Context returns a context-injection-only SessionStart result.
func (sessionStartResults) Context(text string) SessionStartOutput {
	return sessionStartOutput{additionalContext: text}
}

// Noop returns an empty response (silent stdout).
func (sessionStartResults) Noop() SessionStartOutput {
	return sessionStartOutput{}
}

func (sessionStartOutput) allowedEvents() []string {
	return []string{EventSessionStart}
}

func (o sessionStartOutput) encode() ([]byte, int, error) {
	return encodeAdditionalContext(o.additionalContext)
}

func encodeAdditionalContext(context string) ([]byte, int, error) {
	if context == "" {
		return nil, 0, nil
	}
	out := map[string]any{"additionalContext": context}
	b, err := json.Marshal(out)
	return b, 0, err
}

func init() {
	registerDecoder(EventSessionStart, decodeAs[SessionStart])
}

// SessionStart registers a SessionStart handler.
func (c *Chain) SessionStart(fn func(context.Context, Hook[SessionStart], SessionStartResults) (SessionStartOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev SessionStart) (SessionStartOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), sessionStartResults{})
	})
	return &Chain{}
}
