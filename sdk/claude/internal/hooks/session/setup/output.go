package setup

import (
	"maps"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/runtime"
)

// Output is the response for Setup events.
// Construct via Results builders and With* methods.
// A nil value is a no-op.
type Output interface {
	hookkit.Output
	isOutput()
	// WithAdditionalContext injects model context.
	WithAdditionalContext(text string) Output
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
	additionalContext string
	env               map[string]string
}

func (output) isOutput() {}

// IsZero reports whether this hook response is empty.
func (o output) IsZero() bool {
	return o.Common.IsZero() && o.additionalContext == "" && len(o.env) == 0
}

// WithAdditionalContext injects model context.
func (o output) WithAdditionalContext(text string) Output {
	o.additionalContext = text
	return o
}

// WithEnv sets session environment variables written to CLAUDE_ENV_FILE.
// The map is cloned so later caller mutations do not change the encoded payload.
func (o output) WithEnv(env map[string]string) Output {
	o.env = maps.Clone(env)
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
}

// Encode renders this output as Claude Code stdout JSON.
// WithEnv appends export lines to CLAUDE_ENV_FILE when that env var is set.
func (o output) Encode() ([]byte, int, error) {
	if err := runtime.WriteEnvFile(o.env, nil, nil); err != nil {
		return nil, event.SuccessExit, err
	}
	return event.MarshalHookOutput(event.Setup, o.encodeInto)
}

// Merge combines other into this Setup output.
func (o output) Merge(other hookkit.Output) (hookkit.Output, []string, error) {
	b, ok := other.(output)
	if !ok {
		return nil, nil, hookkit.ErrMergeType(o, other)
	}
	mergedCommon, warnings := o.Common.Merge(b.Common)
	env, w := hookkit.TakeLastMap("env", o.env, b.env)
	if w != "" {
		warnings = append(warnings, w)
	}
	return output{
		Common:            mergedCommon,
		additionalContext: hookkit.JoinContextStrings(o.additionalContext, b.additionalContext),
		env:               env,
	}, warnings, nil
}

// Stop reports whether remaining handlers should be skipped.
func (o output) Stop() bool {
	return o.Common.Stop()
}
