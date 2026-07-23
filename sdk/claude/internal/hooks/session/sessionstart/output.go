package sessionstart

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/runtime"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// Output is the response for SessionStart events.
// Construct via Results builders and With* methods.
// A nil value is a no-op.
type Output interface {
	run.Output
	isOutput()
	// WithInitialUserMessage sets the initial user message.
	WithInitialUserMessage(msg string) Output
	// WithAdditionalContext injects model context.
	WithAdditionalContext(text string) Output
	// WithSessionTitle sets the session title.
	WithSessionTitle(title string) Output
	// WithWatchPaths registers filesystem watch paths.
	WithWatchPaths(paths []string) Output
	// WithReloadSkills reloads skills when true.
	WithReloadSkills(v bool) Output
	// WithEnv sets session environment variables written to CLAUDE_ENV_FILE.
	WithEnv(env map[string]string) Output
	// WithContinue sets whether Claude should continue the session.
	WithContinue(v bool) Output
	// WithStopReason explains why the session was stopped.
	WithStopReason(reason string) Output
	// WithSuppressOutput suppresses hook stdout when true.
	WithSuppressOutput(v bool) Output
	// WithSystemMessage sets a user-visible system message.
	WithSystemMessage(msg string) Output
	// WithTerminalSequence sets an OSC terminal sequence (allowlisted).
	WithTerminalSequence(seq string) Output
}

type output struct {
	event.Common
	additionalContext  string
	initialUserMessage string
	sessionTitle       string
	watchPaths         []string
	reloadSkills       bool
	env                map[string]string
}

func (output) isOutput() {}

// IsZero reports whether this hook response is empty.
func (o output) IsZero() bool {
	return o.Common.IsZero() && o.additionalContext == "" && o.initialUserMessage == "" &&
		o.sessionTitle == "" && len(o.watchPaths) == 0 && !o.reloadSkills && len(o.env) == 0
}

// WithInitialUserMessage sets the initial user message.
func (o output) WithInitialUserMessage(msg string) Output {
	o.initialUserMessage = msg
	return o
}

// WithAdditionalContext injects model context.
func (o output) WithAdditionalContext(text string) Output {
	o.additionalContext = text
	return o
}

// WithSessionTitle sets the session title.
func (o output) WithSessionTitle(title string) Output {
	o.sessionTitle = title
	return o
}

// WithWatchPaths registers filesystem watch paths.
func (o output) WithWatchPaths(paths []string) Output {
	o.watchPaths = paths
	return o
}

// WithReloadSkills reloads skills when true.
func (o output) WithReloadSkills(v bool) Output {
	o.reloadSkills = v
	return o
}

// WithEnv sets session environment variables written to CLAUDE_ENV_FILE.
func (o output) WithEnv(env map[string]string) Output {
	o.env = env
	return o
}

// WithContinue sets whether Claude should continue the session.
func (o output) WithContinue(v bool) Output {
	o.Common = o.Common.WithContinue(v)
	return o
}

// WithStopReason explains why the session was stopped.
func (o output) WithStopReason(reason string) Output {
	o.Common = o.Common.WithStopReason(reason)
	return o
}

// WithSuppressOutput suppresses hook stdout when true.
func (o output) WithSuppressOutput(v bool) Output {
	o.Common = o.Common.WithSuppressOutput(v)
	return o
}

// WithSystemMessage sets a user-visible system message.
func (o output) WithSystemMessage(msg string) Output {
	o.Common = o.Common.WithSystemMessage(msg)
	return o
}

// WithTerminalSequence sets an OSC terminal sequence (allowlisted).
func (o output) WithTerminalSequence(seq string) Output {
	o.Common = o.Common.WithTerminalSequence(seq)
	return o
}

func (o output) encodeInto(top, hso map[string]any) {
	event.ApplyCommon(top, o.Common)
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
func (o output) Encode() ([]byte, int, error) {
	if err := runtime.WriteEnvFile(o.env, nil, nil); err != nil {
		return nil, event.SuccessExit, err
	}
	return event.MarshalHookOutput(event.SessionStart, o.encodeInto)
}

// Merge combines other into this SessionStart output.
func (o output) Merge(other run.Output) (run.Output, []string, error) {
	b, ok := other.(output)
	if !ok {
		return nil, nil, hookkit.ErrMergeType(o, other)
	}
	mergedCommon, warnings := o.Common.Merge(b.Common)
	initialUserMessage, w := hookkit.TakeLastString("initialUserMessage", o.initialUserMessage, b.initialUserMessage)
	if w != "" {
		warnings = append(warnings, w)
	}
	sessionTitle, w := hookkit.TakeLastString("sessionTitle", o.sessionTitle, b.sessionTitle)
	if w != "" {
		warnings = append(warnings, w)
	}
	watchPaths, w := hookkit.TakeLastSlice("watchPaths", o.watchPaths, b.watchPaths)
	if w != "" {
		warnings = append(warnings, w)
	}
	env, w := hookkit.TakeLastMap("env", o.env, b.env)
	if w != "" {
		warnings = append(warnings, w)
	}
	return output{
		Common:             mergedCommon,
		additionalContext:  hookkit.JoinContextStrings(o.additionalContext, b.additionalContext),
		initialUserMessage: initialUserMessage,
		sessionTitle:       sessionTitle,
		watchPaths:         watchPaths,
		reloadSkills:       o.reloadSkills || b.reloadSkills,
		env:                env,
	}, warnings, nil
}

// Stop reports whether remaining handlers should be skipped.
func (o output) Stop() bool {
	return o.Common.Stop()
}
