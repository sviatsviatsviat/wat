package claude

import "encoding/json"

// PostToolBatchCall is one tool invocation entry in a PostToolBatch event.
type PostToolBatchCall struct {
	// ToolName is the tool name.
	ToolName string `json:"tool_name"`
	// ToolInput is the native tool input JSON.
	ToolInput json.RawMessage `json:"tool_input"`
	// ToolUseID is the tool use identifier.
	ToolUseID string `json:"tool_use_id"`
	// ToolResponse is the tool response JSON (string or content-block array).
	ToolResponse json.RawMessage `json:"tool_response"`
}

// PostToolBatch is the PostToolBatch hook event.
type PostToolBatch struct {
	Envelope
	// ToolCalls holds per-call metadata for the resolved parallel batch.
	ToolCalls []PostToolBatchCall `json:"tool_calls"`
}

// EventName returns the hook event name.
func (PostToolBatch) EventName() string { return EventPostToolBatch }

func init() {
	registerDecoder(EventPostToolBatch, decodeAs[PostToolBatch])
}
