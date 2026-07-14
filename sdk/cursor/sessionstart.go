package cursor

import (
	"context"
	"encoding/json"

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
type SessionStartOutput struct {
	// Env sets environment variables for the session.
	Env map[string]string
	// AdditionalContext injects model context.
	AdditionalContext string
}

func (o SessionStartOutput) isZero() bool {
	return len(o.Env) == 0 && o.AdditionalContext == ""
}

// SessionStartResults is the hook-scoped response builder supplied to Chain handlers by registration.
type SessionStartResults interface {
	// Context returns a context-injection-only SessionStart result.
	Context(text string) SessionStartOutput
	isSessionStartResults()
}

type sessionStartResults struct{}

func (sessionStartResults) isSessionStartResults() {}

// Context returns a context-injection-only SessionStart result.
func (sessionStartResults) Context(text string) SessionStartOutput {
	return SessionStartOutput{AdditionalContext: text}
}

func encodeSessionStart(o SessionStartOutput) ([]byte, int, error) {
	out := map[string]any{}
	if len(o.Env) > 0 {
		out["env"] = o.Env
	}
	if o.AdditionalContext != "" {
		out["additional_context"] = o.AdditionalContext
	}
	if len(out) == 0 {
		return nil, 0, nil
	}
	b, err := json.Marshal(out)
	return b, 0, err
}

func init() {
	registerDecoder(EventSessionStart, decodeAs[SessionStart])
}

// SessionStart registers a sessionStart handler.
func (c *Chain) SessionStart(fn func(context.Context, SessionStartHook, SessionStartResults) (SessionStartOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev SessionStart) (SessionStartOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), sessionStartResults{})
	})
	return &Chain{}
}
