package beforesubmitprompt

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
)

// Output is the response for beforeSubmitPrompt events.
// Construct via Results builders and With* methods. A nil value is a no-op.
type Output interface {
	hookkit.Output
	isOutput()
	// WithContinue sets whether prompt submission should continue.
	WithContinue(v bool) Output
	// WithUserMessage sets a user-facing message when blocking.
	WithUserMessage(msg string) Output
}

type output struct {
	cont        *bool
	userMessage string
}

func (output) isOutput() {}

// IsZero reports whether this hook response is empty.
func (o output) IsZero() bool {
	return o.cont == nil && o.userMessage == ""
}

// WithContinue sets whether prompt submission should continue.
func (o output) WithContinue(v bool) Output {
	o.cont = &v
	return o
}

// WithUserMessage sets a user-facing message when blocking.
func (o output) WithUserMessage(msg string) Output {
	o.userMessage = msg
	return o
}

// Encode renders this output as Cursor stdout JSON.
func (o output) Encode() ([]byte, int, error) {
	out := map[string]any{}
	if o.cont != nil {
		out["continue"] = *o.cont
	}
	if o.userMessage != "" {
		out["user_message"] = o.userMessage
	}
	if len(out) == 0 {
		return nil, 0, nil
	}
	b, err := json.Marshal(out)
	return b, 0, err
}

// Merge combines other into this beforeSubmitPrompt output.
func (o output) Merge(other hookkit.Output) (hookkit.Output, []string, error) {
	b, ok := other.(output)
	if !ok {
		return nil, nil, hookkit.ErrMergeType(o, other)
	}
	var warnings []string
	cont := o.cont
	if cont != nil && !*cont {
		// continue:false is sticky
	} else if b.cont != nil {
		cont = b.cont
	}
	userMessage, w := hookkit.TakeLastString("userMessage", o.userMessage, b.userMessage)
	if w != "" {
		warnings = append(warnings, w)
	}
	return output{cont: cont, userMessage: userMessage}, warnings, nil
}

// Stop reports whether remaining handlers should be skipped.
func (o output) Stop() bool {
	return o.cont != nil && !*o.cont
}
