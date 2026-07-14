package cursor

import (
	"encoding/json"
)

// PostToolOutput is the response for post-tool events.
type PostToolOutput struct {
	// UpdatedMCPOutput replaces MCP tool output when set.
	UpdatedMCPOutput any
	// AdditionalContext injects model context.
	AdditionalContext string
}

func (o PostToolOutput) isZero() bool {
	return o.UpdatedMCPOutput == nil && o.AdditionalContext == ""
}

// PostToolResults is the hook-scoped response builder supplied to Chain handlers by registration.
type PostToolResults interface {
	// Context returns a context-injection-only PostTool result.
	Context(text string) PostToolOutput
	isPostToolResults()
}

type postToolResults struct{}

func (postToolResults) isPostToolResults() {}

// Context returns a context-injection-only PostTool result.
func (postToolResults) Context(text string) PostToolOutput {
	return PostToolOutput{AdditionalContext: text}
}

func encodePostTool(o PostToolOutput) ([]byte, int, error) {
	out := map[string]any{}
	if o.UpdatedMCPOutput != nil {
		out["updated_mcp_tool_output"] = o.UpdatedMCPOutput
	}
	if o.AdditionalContext != "" {
		out["additional_context"] = o.AdditionalContext
	}
	if len(out) == 0 {
		return nil, 0, nil
	}
	b, err := json.Marshal(out)
	return b, 0, err
}
