package copilot

import (
	"context"
	"encoding/json"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// SessionStart is the SessionStart hook event.
type SessionStart struct {
	Envelope
	// Source is the session start source.
	Source string `json:"source"`
	// InitialPromptValue is the initial prompt.
	InitialPromptValue string `json:"initial_prompt"`
}

// EventName returns the canonical hook event name.
func (SessionStart) EventName() string { return EventSessionStart }

// InitialPrompt returns the initial prompt.
func (e SessionStart) InitialPrompt() string {
	return e.InitialPromptValue
}

// SessionStartOutput is the response for SessionStart events.
// Construct via SessionStartResults builders. A nil value is a no-op.
type SessionStartOutput interface {
	run.Output
	isSessionStartOutput()
}

type sessionStartOutput struct {
	additionalContext string
}

func (sessionStartOutput) isSessionStartOutput() {}

// IsZero reports whether this hook response is empty.
func (o sessionStartOutput) IsZero() bool {
	return o.additionalContext == ""
}

// SessionStartResults is the hook-scoped response builder supplied to On* handlers by registration.
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

// Encode renders this output as Copilot stdout JSON.
func (o sessionStartOutput) Encode() ([]byte, int, error) {
	return encodeAdditionalContext(o.additionalContext)
}

func encodeAdditionalContext(context string) ([]byte, int, error) {
	if context == "" {
		return nil, 0, nil
	}
	out := map[string]any{"additional_context": context}
	b, err := json.Marshal(out)
	return b, 0, err
}

func init() {
	codec.Register(EventSessionStart, hookkit.EventDecoder[SessionStart](codec))
}

// SessionStart registers a SessionStart handler on the chain.
func (c *chain) SessionStart(fn func(context.Context, run.Hook[SessionStart], SessionStartResults) (SessionStartOutput, error)) *chain {
	if fn == nil {
		return c
	}
	c.reg.RegisterHandler(Dialect, run.Handler(func(ctx context.Context, hook run.Hook[SessionStart]) (SessionStartOutput, error) {
		return fn(ctx, hook, sessionStartResults{})
	}))
	return c
}
