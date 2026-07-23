package posttooluse

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// Output is the response for PostToolUse events.
// Construct via Results builders and With* methods. A nil value is a no-op.
type Output interface {
	run.Output
	isOutput()
	// WithModifiedResult replaces the tool result text when set.
	WithModifiedResult(result string) Output
}

type output struct {
	modifiedResult    string
	additionalContext string
}

func (output) isOutput() {}

// IsZero reports whether this hook response is empty.
func (o output) IsZero() bool {
	return o.modifiedResult == "" && o.additionalContext == ""
}

// WithModifiedResult replaces the tool result text when set.
func (o output) WithModifiedResult(result string) Output {
	o.modifiedResult = result
	return o
}

// Encode renders this output as Copilot stdout JSON.
func (o output) Encode() ([]byte, int, error) {
	out := map[string]any{}
	if o.modifiedResult != "" {
		out["modified_result"] = map[string]any{
			"result_type":         "success",
			"text_result_for_llm": o.modifiedResult,
		}
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

// Merge combines other into the receiver. other must be an output.
func (o output) Merge(other run.Output) (run.Output, []string, error) {
	b, ok := other.(output)
	if !ok {
		return nil, nil, hookkit.ErrMergeType(o, other)
	}
	modifiedResult, warn := hookkit.TakeLastString("modified_result", o.modifiedResult, b.modifiedResult)
	var warnings []string
	if warn != "" {
		warnings = append(warnings, warn)
	}
	return output{
		modifiedResult:    modifiedResult,
		additionalContext: hookkit.JoinContextStrings(o.additionalContext, b.additionalContext),
	}, warnings, nil
}

// Stop reports whether remaining handlers should be skipped.
func (o output) Stop() bool {
	return false
}
