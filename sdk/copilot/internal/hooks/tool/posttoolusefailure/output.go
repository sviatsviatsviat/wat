package posttoolusefailure

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/event"
)

// Output is the response for PostToolUseFailure events.
// Construct via Results builders. A nil value is a no-op.
type Output interface {
	hookkit.Output
	isOutput()
}

type output struct {
	context string
}

func (output) isOutput() {}

// IsZero reports whether this hook response is empty.
func (o output) IsZero() bool {
	return o.context == ""
}

// Encode renders this output as Copilot stdout JSON.
func (o output) Encode() ([]byte, int, error) {
	if o.context == "" {
		return nil, 0, nil
	}
	return []byte(o.context), event.WarnExit, nil
}

// Merge combines other into the receiver. other must be an output.
func (o output) Merge(other hookkit.Output) (hookkit.Output, []string, error) {
	b, ok := other.(output)
	if !ok {
		return nil, nil, hookkit.ErrMergeType(o, other)
	}
	return output{
		context: hookkit.JoinContextStrings(o.context, b.context),
	}, nil, nil
}

// Stop reports whether remaining handlers should be skipped.
func (o output) Stop() bool {
	return false
}
