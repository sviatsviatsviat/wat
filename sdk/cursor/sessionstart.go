package cursor

import (
	"context"
	"encoding/json"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// SessionStart is the sessionStart hook event.
type SessionStart struct {
	Envelope
	// IsBackgroundAgent reports whether this is a background agent session.
	IsBackgroundAgent bool `json:"is_background_agent"`
}

// EventName returns the canonical hook event name.
func (SessionStart) EventName() string { return EventSessionStart }

// SessionStartOutput is the response for sessionStart events.
// Construct via SessionStartResults builders and With* methods. A nil value is a no-op.
type SessionStartOutput interface {
	Output
	isSessionStartOutput()
	// WithEnv sets environment variables for the session.
	WithEnv(env map[string]string) SessionStartOutput
	// WithAdditionalContext injects model context.
	WithAdditionalContext(text string) SessionStartOutput
}

type sessionStartOutput struct {
	env               map[string]string
	additionalContext string
}

func (sessionStartOutput) isCursorOutput() {}

func (sessionStartOutput) isSessionStartOutput() {}

func (o sessionStartOutput) isZero() bool {
	return len(o.env) == 0 && o.additionalContext == ""
}

// WithEnv sets environment variables for the session.
func (o sessionStartOutput) WithEnv(env map[string]string) SessionStartOutput {
	o.env = env
	return o
}

// WithAdditionalContext injects model context.
func (o sessionStartOutput) WithAdditionalContext(text string) SessionStartOutput {
	o.additionalContext = text
	return o
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

func (sessionStartOutput) allowedEvents() []string {
	return []string{EventSessionStart}
}

func (o sessionStartOutput) encode(eventName string) ([]byte, int, error) {
	_ = eventName
	out := map[string]any{}
	if len(o.env) > 0 {
		out["env"] = o.env
	}
	if o.additionalContext != "" {
		out["additional_context"] = o.additionalContext
	}
	if len(out) == 0 {
		return nil, 0, nil
	}
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
	registerHandler(c.reg, func(ctx context.Context, ev SessionStart) (SessionStartOutput, error) {
		return fn(ctx, run.NewHook(run.InvocationFrom(ctx), ev), sessionStartResults{})
	})
	return c
}
