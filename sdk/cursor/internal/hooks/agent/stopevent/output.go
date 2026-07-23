package stopevent

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// Output is the response for stop and subagentStop events.
// Construct via Results builders. A nil value is a no-op.
type Output interface {
	run.Output
	isOutput()
}

type output struct {
	followUpMessage string
}

func (output) isOutput() {}

// IsZero reports whether this hook response is empty.
func (o output) IsZero() bool {
	return o.followUpMessage == ""
}

// Encode renders this output as Cursor stdout JSON.
func (o output) Encode() ([]byte, int, error) {
	if o.followUpMessage == "" {
		return nil, 0, nil
	}
	out := map[string]any{"followup_message": o.followUpMessage}
	b, err := json.Marshal(out)
	return b, 0, err
}

// Merge combines other into this stop output.
func (o output) Merge(other run.Output) (run.Output, []string, error) {
	b, ok := other.(output)
	if !ok {
		return nil, nil, hookkit.ErrMergeType(o, other)
	}
	var warnings []string
	followUpMessage, w := hookkit.TakeLastString("followUpMessage", o.followUpMessage, b.followUpMessage)
	if w != "" {
		warnings = append(warnings, w)
	}
	return output{followUpMessage: followUpMessage}, warnings, nil
}

// Stop reports whether remaining handlers should be skipped.
func (o output) Stop() bool {
	return false
}
