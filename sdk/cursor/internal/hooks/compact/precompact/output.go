package precompact

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// Output is the response for preCompact events.
// Construct via Results builders. A nil value is a no-op.
type Output interface {
	run.Output
	isOutput()
}

type output struct {
	userMessage string
}

func (output) isOutput() {}

// IsZero reports whether this hook response is empty.
func (o output) IsZero() bool {
	return o.userMessage == ""
}

// Encode renders this output as Cursor stdout JSON.
func (o output) Encode() ([]byte, int, error) {
	if o.userMessage == "" {
		return nil, 0, nil
	}
	out := map[string]any{"user_message": o.userMessage}
	b, err := json.Marshal(out)
	return b, 0, err
}

// Merge combines other into this preCompact output.
func (o output) Merge(other run.Output) (run.Output, []string, error) {
	b, ok := other.(output)
	if !ok {
		return nil, nil, hookkit.ErrMergeType(o, other)
	}
	var warnings []string
	userMessage, w := hookkit.TakeLastString("userMessage", o.userMessage, b.userMessage)
	if w != "" {
		warnings = append(warnings, w)
	}
	return output{userMessage: userMessage}, warnings, nil
}

// Stop reports whether remaining handlers should be skipped.
func (o output) Stop() bool {
	return false
}
