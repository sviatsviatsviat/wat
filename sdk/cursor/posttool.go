package cursor

import (
	"encoding/json"
)

// PostToolOutput is the response for post-tool events.
// Construct via PostToolResults builders and With* methods. A nil value is a no-op.
type PostToolOutput interface {
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

func (o postToolOutput) isZero() bool {
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

// PostToolResults is the hook-scoped response builder supplied to Chain handlers by registration.
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

func (postToolOutput) allowedEvents() []string {
	return []string{
		EventPostToolUse,
		EventPostToolUseFailure,
		EventAfterMCPExecution,
		EventAfterShellExecution,
		EventAfterFileEdit,
	}
}

func (o postToolOutput) encode(eventName string) ([]byte, int, error) {
	_ = eventName
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
