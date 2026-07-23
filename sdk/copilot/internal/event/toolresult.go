package event

import "encoding/json"

// ToolResult is a Copilot tool result object.
type ToolResult struct {
	// ResultType is the result type.
	ResultType string `json:"result_type"`
	// TextResultForLLM is the LLM-facing text.
	TextResultForLLM string `json:"text_result_for_llm"`
}

// Text returns the textual result.
func (r ToolResult) Text() string {
	return r.TextResultForLLM
}

// ErrorDetail is a structured Copilot error object.
type ErrorDetail struct {
	// Message is the error message.
	Message string `json:"message"`
	// Name is the error category name.
	Name string `json:"name"`
	// Stack is the error stack trace when provided.
	Stack string `json:"stack"`
}

// MarshalToolResult encodes a ToolResult as JSON for ResultRaw helpers.
func MarshalToolResult(r ToolResult) json.RawMessage {
	out := map[string]string{}
	if r.ResultType != "" {
		out["result_type"] = r.ResultType
	}
	if r.TextResultForLLM != "" {
		out["text_result_for_llm"] = r.TextResultForLLM
	}
	if len(out) == 0 {
		return nil
	}
	b, err := json.Marshal(out)
	if err != nil {
		return nil
	}
	return b
}
