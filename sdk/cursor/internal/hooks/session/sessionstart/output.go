package sessionstart

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
)

// Output is the response for sessionStart events.
// Construct via Results builders and With* methods. A nil value is a no-op.
//
// Meaningful fields are env and additional_context. Cursor's schema also
// accepts continue and user_message, but the host does not enforce them for
// sessionStart, so this SDK does not expose builders for those fields.
type Output interface {
	hookkit.Output
	isOutput()
	// WithEnv sets environment variables for the session.
	WithEnv(env map[string]string) Output
	// WithAdditionalContext injects model context.
	WithAdditionalContext(text string) Output
}

type output struct {
	env               map[string]string
	additionalContext string
}

func (output) isOutput() {}

// IsZero reports whether this hook response is empty.
func (o output) IsZero() bool {
	return len(o.env) == 0 && o.additionalContext == ""
}

// WithEnv sets environment variables for the session.
func (o output) WithEnv(env map[string]string) Output {
	o.env = env
	return o
}

// WithAdditionalContext injects model context.
func (o output) WithAdditionalContext(text string) Output {
	o.additionalContext = text
	return o
}

// Encode renders this output as Cursor stdout JSON.
func (o output) Encode() ([]byte, int, error) {
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

// Merge combines other into this sessionStart output.
func (o output) Merge(other hookkit.Output) (hookkit.Output, []string, error) {
	b, ok := other.(output)
	if !ok {
		return nil, nil, hookkit.ErrMergeType(o, other)
	}
	var warnings []string
	env, w := hookkit.TakeLastMap("env", o.env, b.env)
	if w != "" {
		warnings = append(warnings, w)
	}
	return output{
		env:               env,
		additionalContext: hookkit.JoinContextStrings(o.additionalContext, b.additionalContext),
	}, warnings, nil
}

// Stop reports whether remaining handlers should be skipped.
func (o output) Stop() bool {
	return false
}
