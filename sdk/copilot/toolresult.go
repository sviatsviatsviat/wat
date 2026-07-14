package copilot

import (
	"encoding/json"
)

// ToolResult is a Copilot tool result object in either wire format.
type ToolResult struct {
	// ResultType is the result type (camelCase).
	ResultType string `json:"resultType"`
	// TextResultForLLM is the LLM-facing text (camelCase).
	TextResultForLLM string `json:"textResultForLlm"`
	// ResultTypeSnake is the result type (VS Code).
	ResultTypeSnake string `json:"result_type"`
	// TextResultSnake is the LLM-facing text (VS Code).
	TextResultSnake string `json:"text_result_for_llm"`
}

// Text returns the textual result from either wire format.
func (r ToolResult) Text() string {
	if r.TextResultForLLM != "" {
		return r.TextResultForLLM
	}
	return r.TextResultSnake
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

// RawEvent holds an unknown hook event with the full payload preserved.
type RawEvent struct {
	Envelope
	// Raw is the untouched native JSON payload.
	Raw json.RawMessage
}

// EventName returns the canonical or received hook event name.
func (e RawEvent) EventName() string {
	if e.canonical != "" {
		return e.canonical
	}
	if e.receivedName != "" {
		return e.receivedName
	}
	return "unknown"
}
