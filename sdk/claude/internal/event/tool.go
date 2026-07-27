package event

import "github.com/sviatsviatsviat/wat/sdk/claude/internal/tools"

// ToolFields holds Claude tool invocation wire fields shared by tool events.
//
// Embed on events that carry tool_name / tool_input / tool_use_id, then call
// BindToolInput from the DecodeEvent after-callback.
type ToolFields struct {
	// ToolName is the tool name.
	ToolName string `json:"tool_name"`
	// ToolInput is the typed tool input for ToolName.
	ToolInput tools.Input `json:"-"`
	// ToolUseID is the tool use identifier.
	ToolUseID string `json:"tool_use_id"`
}

// BindToolInput binds ToolInput from the tool_input object field on raw JSON.
// Call from the DecodeEvent after-callback.
func (t *ToolFields) BindToolInput(raw []byte) {
	t.ToolInput = tools.NewInputFromPayload(t.ToolName, raw, "tool_input")
}
