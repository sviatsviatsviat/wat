package subagentstop

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
)

// Output is the response for SubagentStop events.
// Construct via Results builders and With* methods. A nil value is a no-op.
//
// Copilot host rules for this event:
//   - decision "block" with reason forces another subagent turn;
//   - modifiedResponse replaces the text returned to the parent when the
//     subagent is allowed to complete;
//   - a valid block wins over modifiedResponse (rewrite is discarded);
//   - rewrites do not compose: every handler sees the original response text,
//     and the last non-empty modifiedResponse wins.
type Output interface {
	hookkit.Output
	isOutput()
	// WithModifiedResponse replaces the response returned to the parent when
	// the subagent is allowed to complete.
	WithModifiedResponse(text string) Output
}

type output struct {
	reason           string
	modifiedResponse string
}

func (output) isOutput() {}

// IsZero reports whether this hook response is empty.
func (o output) IsZero() bool {
	return o.reason == "" && o.modifiedResponse == ""
}

// WithModifiedResponse replaces the response returned to the parent when the
// subagent is allowed to complete.
func (o output) WithModifiedResponse(text string) Output {
	o.modifiedResponse = text
	return o
}

// Encode renders this output as Copilot stdout JSON.
// Field names decision, reason, and modifiedResponse match the Copilot Hooks
// reference for both camelCase and VS Code compatible configs.
func (o output) Encode() ([]byte, int, error) {
	if o.reason != "" {
		out := map[string]any{
			"decision": "block",
			"reason":   o.reason,
		}
		b, err := json.Marshal(out)
		return b, 0, err
	}
	if o.modifiedResponse == "" {
		return nil, 0, nil
	}
	out := map[string]any{
		"modifiedResponse": o.modifiedResponse,
	}
	b, err := json.Marshal(out)
	return b, 0, err
}

// Merge combines other into the receiver. other must be an output.
// Follow-up reason and modifiedResponse each take the last non-empty value.
func (o output) Merge(other hookkit.Output) (hookkit.Output, []string, error) {
	b, ok := other.(output)
	if !ok {
		return nil, nil, hookkit.ErrMergeType(o, other)
	}
	reason, warnReason := hookkit.TakeLastString("reason", o.reason, b.reason)
	modified, warnMod := hookkit.TakeLastString("modifiedResponse", o.modifiedResponse, b.modifiedResponse)
	var warnings []string
	if warnReason != "" {
		warnings = append(warnings, warnReason)
	}
	if warnMod != "" {
		warnings = append(warnings, warnMod)
	}
	return output{reason: reason, modifiedResponse: modified}, warnings, nil
}

// Stop reports whether remaining handlers should be skipped.
// A FollowUp (block) stops the chain; a response rewrite alone does not, so
// later handlers can still supply the winning modifiedResponse.
func (o output) Stop() bool {
	return o.reason != ""
}
