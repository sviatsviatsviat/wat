package cursor

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// PostToolOutput is the response for post-tool events.
// Construct via PostToolResults builders and With* methods. A nil value is a no-op.
type PostToolOutput interface {
	run.Output
	isPostToolOutput()
	// WithUpdatedMCPOutput replaces MCP tool output when set.
	WithUpdatedMCPOutput(output any) PostToolOutput
	// WithAdditionalContext injects model context.
	WithAdditionalContext(text string) PostToolOutput
}

type postToolOutput struct {
	updatedMCPOutput  any
	additionalContext string
}

func (postToolOutput) isPostToolOutput() {}

// IsZero reports whether this hook response is empty.
func (o postToolOutput) IsZero() bool {
	return o.updatedMCPOutput == nil && o.additionalContext == ""
}

// WithUpdatedMCPOutput replaces MCP tool output when set.
func (o postToolOutput) WithUpdatedMCPOutput(output any) PostToolOutput {
	o.updatedMCPOutput = output
	return o
}

// WithAdditionalContext injects model context.
func (o postToolOutput) WithAdditionalContext(text string) PostToolOutput {
	o.additionalContext = text
	return o
}

// PostToolResults is the hook-scoped response builder supplied to On* handlers by registration.
type PostToolResults interface {
	// Context returns a context-injection-only PostTool result.
	Context(text string) PostToolOutput
	// Noop returns an empty response (silent stdout). Prefer nil from handlers when not chaining With*.
	Noop() PostToolOutput
	isPostToolResults()
}

type postToolResults struct{}

func (postToolResults) isPostToolResults() {}

// Context returns a context-injection-only PostTool result.
func (postToolResults) Context(text string) PostToolOutput {
	return postToolOutput{additionalContext: text}
}

// Noop returns an empty response (silent stdout).
func (postToolResults) Noop() PostToolOutput {
	return postToolOutput{}
}

// Encode renders this output as Cursor stdout JSON.
func (o postToolOutput) Encode() ([]byte, int, error) {
	out := map[string]any{}
	if o.updatedMCPOutput != nil {
		out["updated_mcp_tool_output"] = o.updatedMCPOutput
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

// Merge combines other into this post-tool output.
func (o postToolOutput) Merge(other run.Output) (run.Output, []string, error) {
	b, ok := other.(postToolOutput)
	if !ok {
		return nil, nil, hookkit.ErrMergeType(o, other)
	}
	var warnings []string
	updatedMCPOutput := o.updatedMCPOutput
	if b.updatedMCPOutput != nil {
		if o.updatedMCPOutput != nil {
			warnings = append(warnings, hookkit.OverwriteWarning("updatedMCPOutput"))
		}
		updatedMCPOutput = b.updatedMCPOutput
	}
	return postToolOutput{
		updatedMCPOutput:  updatedMCPOutput,
		additionalContext: hookkit.JoinContextStrings(o.additionalContext, b.additionalContext),
	}, warnings, nil
}

// Stop reports whether remaining handlers should be skipped.
func (o postToolOutput) Stop() bool {
	return false
}
