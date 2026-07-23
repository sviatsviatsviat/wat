package sessionstart

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// Output is the response for SessionStart events.
// Construct via Results builders. A nil value is a no-op.
type Output interface {
	run.Output
	isOutput()
}

type output struct {
	additionalContext string
}

func (output) isOutput() {}

// IsZero reports whether this hook response is empty.
func (o output) IsZero() bool {
	return o.additionalContext == ""
}

// Encode renders this output as Copilot stdout JSON.
func (o output) Encode() ([]byte, int, error) {
	return event.EncodeAdditionalContext(o.additionalContext)
}

// Merge combines other into the receiver. other must be an output.
func (o output) Merge(other run.Output) (run.Output, []string, error) {
	b, ok := other.(output)
	if !ok {
		return nil, nil, hookkit.ErrMergeType(o, other)
	}
	return output{
		additionalContext: hookkit.JoinContextStrings(o.additionalContext, b.additionalContext),
	}, nil, nil
}

// Stop reports whether remaining handlers should be skipped.
func (o output) Stop() bool {
	return false
}
