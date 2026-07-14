package claude

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// SessionStart is the SessionStart hook event.
type SessionStart struct {
	Envelope
	// Source is the session start source (startup, resume, clear, compact).
	Source string `json:"source"`
	// Model is the model name when provided.
	Model string `json:"model,omitempty"`
	// SessionTitle is the session title when provided.
	SessionTitle string `json:"session_title,omitempty"`
}

// EventName returns the hook event name.
func (SessionStart) EventName() string { return EventSessionStart }

func init() {
	registerDecoder(EventSessionStart, decodeAs[SessionStart])
}

// SessionStartOutput is the response for SessionStart events.
type SessionStartOutput struct {
	Common
	// AdditionalContext injects model context.
	AdditionalContext string
	// InitialUserMessage sets the initial user message.
	InitialUserMessage string
	// SessionTitle sets the session title.
	SessionTitle string
	// WatchPaths registers filesystem watch paths.
	WatchPaths []string
	// ReloadSkills reloads skills when true.
	ReloadSkills bool
	// Env carries session environment variables written to CLAUDE_ENV_FILE.
	Env map[string]string
}

func (o SessionStartOutput) isZero() bool {
	return o.Common.isZero() && o.AdditionalContext == "" && o.InitialUserMessage == "" &&
		o.SessionTitle == "" && len(o.WatchPaths) == 0 && !o.ReloadSkills && len(o.Env) == 0
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

// SessionStart registers a SessionStart handler.
func (c *Chain) SessionStart(fn func(context.Context, SessionStartHook, SessionStartResults) (SessionStartOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev SessionStart) (SessionStartOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), sessionStartResults{})
	})
	return &Chain{}
}
