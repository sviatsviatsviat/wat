package userprompttransformed

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
)

// Output is the response for UserPromptTransformed events.
// Construct via Results builders and With* methods. A nil value is a no-op.
type Output interface {
	hookkit.Output
	isOutput()
	// WithModifiedTransformedPrompt replaces the model-facing transformed prompt when set.
	WithModifiedTransformedPrompt(text string) Output
}

type output struct {
	modifiedTransformedPrompt string
}

func (output) isOutput() {}

// IsZero reports whether this hook response is empty.
func (o output) IsZero() bool {
	return o.modifiedTransformedPrompt == ""
}

// WithModifiedTransformedPrompt replaces the model-facing transformed prompt when set.
func (o output) WithModifiedTransformedPrompt(text string) Output {
	o.modifiedTransformedPrompt = text
	return o
}

// Encode renders this output as Copilot stdout JSON.
func (o output) Encode() ([]byte, int, error) {
	if o.modifiedTransformedPrompt == "" {
		return nil, 0, nil
	}
	b, err := json.Marshal(map[string]any{
		"modified_transformed_prompt": o.modifiedTransformedPrompt,
	})
	return b, 0, err
}

// Merge combines other into the receiver. other must be an output.
func (o output) Merge(other hookkit.Output) (hookkit.Output, []string, error) {
	b, ok := other.(output)
	if !ok {
		return nil, nil, hookkit.ErrMergeType(o, other)
	}
	modified, warn := hookkit.TakeLastString("modified_transformed_prompt", o.modifiedTransformedPrompt, b.modifiedTransformedPrompt)
	var warnings []string
	if warn != "" {
		warnings = append(warnings, warn)
	}
	return output{modifiedTransformedPrompt: modified}, warnings, nil
}

// Stop reports whether remaining handlers should be skipped.
func (o output) Stop() bool {
	return false
}
