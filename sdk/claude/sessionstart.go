package claude

import (
	"context"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

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
	codec.Register(EventSessionStart, hookkit.EventDecoder[SessionStart](codec))
}

// SessionStartOutput is the response for SessionStart events.
// Construct via SessionStartResults builders and With* methods.
// A nil value is a no-op.
type SessionStartOutput interface {
	run.Output
	isSessionStartOutput()
	// WithInitialUserMessage sets the initial user message.
	WithInitialUserMessage(msg string) SessionStartOutput
	// WithAdditionalContext injects model context.
	WithAdditionalContext(text string) SessionStartOutput
	// WithSessionTitle sets the session title.
	WithSessionTitle(title string) SessionStartOutput
	// WithWatchPaths registers filesystem watch paths.
	WithWatchPaths(paths []string) SessionStartOutput
	// WithReloadSkills reloads skills when true.
	WithReloadSkills(v bool) SessionStartOutput
	// WithEnv sets session environment variables written to CLAUDE_ENV_FILE.
	WithEnv(env map[string]string) SessionStartOutput
	// WithContinue sets whether Claude should continue the session.
	WithContinue(v bool) SessionStartOutput
	// WithStopReason explains why the session was stopped.
	WithStopReason(reason string) SessionStartOutput
	// WithSuppressOutput suppresses hook stdout when true.
	WithSuppressOutput(v bool) SessionStartOutput
	// WithSystemMessage sets a user-visible system message.
	WithSystemMessage(msg string) SessionStartOutput
	// WithTerminalSequence sets an OSC terminal sequence (allowlisted).
	WithTerminalSequence(seq string) SessionStartOutput
}

type sessionStartOutput struct {
	common
	additionalContext  string
	initialUserMessage string
	sessionTitle       string
	watchPaths         []string
	reloadSkills       bool
	env                map[string]string
}

func (sessionStartOutput) isSessionStartOutput() {}

// IsZero reports whether this hook response is empty.
func (o sessionStartOutput) IsZero() bool {
	return o.common.IsZero() && o.additionalContext == "" && o.initialUserMessage == "" &&
		o.sessionTitle == "" && len(o.watchPaths) == 0 && !o.reloadSkills && len(o.env) == 0
}

// WithInitialUserMessage sets the initial user message.
func (o sessionStartOutput) WithInitialUserMessage(msg string) SessionStartOutput {
	o.initialUserMessage = msg
	return o
}

// WithAdditionalContext injects model context.
func (o sessionStartOutput) WithAdditionalContext(text string) SessionStartOutput {
	o.additionalContext = text
	return o
}

// WithSessionTitle sets the session title.
func (o sessionStartOutput) WithSessionTitle(title string) SessionStartOutput {
	o.sessionTitle = title
	return o
}

// WithWatchPaths registers filesystem watch paths.
func (o sessionStartOutput) WithWatchPaths(paths []string) SessionStartOutput {
	o.watchPaths = paths
	return o
}

// WithReloadSkills reloads skills when true.
func (o sessionStartOutput) WithReloadSkills(v bool) SessionStartOutput {
	o.reloadSkills = v
	return o
}

// WithEnv sets session environment variables written to CLAUDE_ENV_FILE.
func (o sessionStartOutput) WithEnv(env map[string]string) SessionStartOutput {
	o.env = env
	return o
}

// WithContinue sets whether Claude should continue the session.
func (o sessionStartOutput) WithContinue(v bool) SessionStartOutput {
	o.common = o.common.WithContinue(v)
	return o
}

// WithStopReason explains why the session was stopped.
func (o sessionStartOutput) WithStopReason(reason string) SessionStartOutput {
	o.common = o.common.WithStopReason(reason)
	return o
}

// WithSuppressOutput suppresses hook stdout when true.
func (o sessionStartOutput) WithSuppressOutput(v bool) SessionStartOutput {
	o.common = o.common.WithSuppressOutput(v)
	return o
}

// WithSystemMessage sets a user-visible system message.
func (o sessionStartOutput) WithSystemMessage(msg string) SessionStartOutput {
	o.common = o.common.WithSystemMessage(msg)
	return o
}

// WithTerminalSequence sets an OSC terminal sequence (allowlisted).
func (o sessionStartOutput) WithTerminalSequence(seq string) SessionStartOutput {
	o.common = o.common.WithTerminalSequence(seq)
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

func (o sessionStartOutput) encodeInto(top, hso map[string]any) {
	applyCommon(top, o.common)
	if o.additionalContext != "" {
		hso["additionalContext"] = o.additionalContext
	}
	if o.sessionTitle != "" {
		hso["sessionTitle"] = o.sessionTitle
	}
	if o.initialUserMessage != "" {
		hso["initialUserMessage"] = o.initialUserMessage
	}
	if len(o.watchPaths) > 0 {
		hso["watchPaths"] = o.watchPaths
	}
	if o.reloadSkills {
		hso["reloadSkills"] = true
	}
}

// Encode renders this output as Claude Code stdout JSON.
// WithEnv appends export lines to CLAUDE_ENV_FILE when that env var is set.
func (o sessionStartOutput) Encode() ([]byte, int, error) {
	if err := writeEnvFile(o.env, nil, nil); err != nil {
		return nil, SuccessExit, err
	}
	return marshalHookOutput(EventSessionStart, o.encodeInto)
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
