package event

import "github.com/sviatsviatsviat/wat/internal/hookkit"

// ContextOutput is an additional_context-only response for Copilot events that
// only accept that field. Construct via ContextResult.
type ContextOutput interface {
	hookkit.Output
	isContextOutput()
}

type contextOutput struct {
	additionalContext string
}

func (contextOutput) isContextOutput() {}

// ContextResult returns a ContextOutput with optional additional context.
func ContextResult(text string) ContextOutput {
	return contextOutput{additionalContext: text}
}

// IsZero reports whether this hook response is empty.
func (o contextOutput) IsZero() bool {
	return o.additionalContext == ""
}

// Encode renders this output as Copilot stdout JSON.
func (o contextOutput) Encode() ([]byte, int, error) {
	return EncodeAdditionalContext(o.additionalContext)
}

// Merge combines other into the receiver. other must be a contextOutput.
func (o contextOutput) Merge(other hookkit.Output) (hookkit.Output, []string, error) {
	b, ok := other.(contextOutput)
	if !ok {
		return nil, nil, hookkit.ErrMergeType(o, other)
	}
	return contextOutput{
		additionalContext: hookkit.JoinContextStrings(o.additionalContext, b.additionalContext),
	}, nil, nil
}

// Stop reports whether remaining handlers should be skipped.
func (o contextOutput) Stop() bool {
	return false
}
