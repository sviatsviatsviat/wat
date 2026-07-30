package agentstop

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
)

// Output is the response for AgentStop events.
// Construct via Results builders. A nil value is a no-op.
// SubagentStop uses its own Output type, which also supports modifiedResponse.
type Output interface {
	hookkit.Output
	isOutput()
}

type output struct {
	reason string
}

func (output) isOutput() {}

// IsZero reports whether this hook response is empty.
func (o output) IsZero() bool {
	return o.reason == ""
}

// Encode renders this output as Copilot stdout JSON.
func (o output) Encode() ([]byte, int, error) {
	if o.reason == "" {
		return nil, 0, nil
	}
	out := map[string]any{
		"decision": "block",
		"reason":   o.reason,
	}
	b, err := json.Marshal(out)
	return b, 0, err
}

// Merge combines other into the receiver. other must be an output.
func (o output) Merge(other hookkit.Output) (hookkit.Output, []string, error) {
	b, ok := other.(output)
	if !ok {
		return nil, nil, hookkit.ErrMergeType(o, other)
	}
	reason, warn := hookkit.TakeLastString("reason", o.reason, b.reason)
	var warnings []string
	if warn != "" {
		warnings = append(warnings, warn)
	}
	return output{reason: reason}, warnings, nil
}

// Stop reports whether remaining handlers should be skipped.
func (o output) Stop() bool {
	return o.reason != ""
}
