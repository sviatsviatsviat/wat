package pretooluse

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/event"
)

// Output is the response for PreToolUse events.
// Construct via Results builders and With* methods. A nil value is a no-op.
type Output interface {
	hookkit.Output
	isOutput()
	// WithModifiedArgs replaces tool arguments when set.
	WithModifiedArgs(args map[string]any) Output
}

type output struct {
	decision     event.PermissionDecision
	reason       string
	modifiedArgs map[string]any
}

func (output) isOutput() {}

// IsZero reports whether this hook response is empty.
func (o output) IsZero() bool {
	return o.decision == "" && o.reason == "" && o.modifiedArgs == nil
}

// WithModifiedArgs replaces tool arguments when set.
func (o output) WithModifiedArgs(args map[string]any) Output {
	o.modifiedArgs = args
	return o
}

// Encode renders this output as Copilot stdout JSON.
func (o output) Encode() ([]byte, int, error) {
	out := map[string]any{}
	if o.decision != "" {
		out["permission_decision"] = string(o.decision)
		if o.reason != "" {
			out["permission_decision_reason"] = o.reason
		}
	}
	if o.modifiedArgs != nil {
		out["modified_args"] = o.modifiedArgs
	}
	if len(out) == 0 {
		return nil, 0, nil
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
	decision, reason := hookkit.MergeRankedString(
		string(o.decision), o.reason,
		string(b.decision), b.reason,
		hookkit.PermissionRankString,
	)
	modifiedArgs, warn := hookkit.TakeLastMap("modified_args", o.modifiedArgs, b.modifiedArgs)
	var warnings []string
	if warn != "" {
		warnings = append(warnings, warn)
	}
	return output{
		decision:     event.PermissionDecision(decision),
		reason:       reason,
		modifiedArgs: modifiedArgs,
	}, warnings, nil
}

// Stop reports whether remaining handlers should be skipped.
func (o output) Stop() bool {
	return o.decision == event.DecisionDeny
}
