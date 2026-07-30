package cwdchanged

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/runtime"
)

// Output is the response for CwdChanged events.
// Construct via Results builders and With* methods.
// A nil value is a no-op.
type Output interface {
	hookkit.Output
	isOutput()
	// WithWatchPaths replaces the dynamic FileChanged watch list.
	WithWatchPaths(paths []string) Output
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
	watchPaths []string
	env        map[string]string
}

func (output) isOutput() {}

// IsZero reports whether this hook response is empty.
func (o output) IsZero() bool {
	return o.Common.IsZero() && o.watchPaths == nil && len(o.env) == 0
}

// WithWatchPaths replaces the dynamic FileChanged watch list.
// An empty slice clears the dynamic watch list.
func (o output) WithWatchPaths(paths []string) Output {
	if paths == nil {
		o.watchPaths = []string{}
	} else {
		o.watchPaths = append([]string(nil), paths...)
	}
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
	if o.watchPaths != nil {
		hso["watchPaths"] = o.watchPaths
	}
}

// Encode renders this output as Claude Code stdout JSON.
// WithEnv appends export lines to CLAUDE_ENV_FILE when that env var is set.
func (o output) Encode() ([]byte, int, error) {
	if err := runtime.WriteEnvFile(o.env, nil, nil); err != nil {
		return nil, event.SuccessExit, err
	}
	return event.MarshalHookOutput(event.CwdChanged, o.encodeInto)
}

// Merge combines other into this CwdChanged output.
func (o output) Merge(other hookkit.Output) (hookkit.Output, []string, error) {
	b, ok := other.(output)
	if !ok {
		return nil, nil, hookkit.ErrMergeType(o, other)
	}
	mergedCommon, warnings := o.Common.Merge(b.Common)
	watchPaths, w := hookkit.TakeLastSlice("watchPaths", o.watchPaths, b.watchPaths)
	if w != "" {
		warnings = append(warnings, w)
	}
	env, w := hookkit.TakeLastMap("env", o.env, b.env)
	if w != "" {
		warnings = append(warnings, w)
	}
	return output{
		Common:     mergedCommon,
		watchPaths: watchPaths,
		env:        env,
	}, warnings, nil
}

// Stop reports whether remaining handlers should be skipped.
func (o output) Stop() bool {
	return o.Common.Stop()
}
