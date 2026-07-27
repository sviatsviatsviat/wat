package event

import "github.com/sviatsviatsviat/wat/sdk/copilot/internal/tools"

// ToolFields holds Copilot tool invocation wire fields shared by tool events.
//
// Embed on events that carry tool_name / tool_input, then call BindToolInput
// from the DecodeEvent after-callback.
type ToolFields struct {
	// ToolName is the tool name.
	ToolName string `json:"tool_name"`
	// ToolInput is the typed tool input.
	ToolInput tools.Input `json:"-"`
}

// BindToolInput binds ToolInput from the tool_input object field on raw JSON.
// Call from the DecodeEvent after-callback.
func (t *ToolFields) BindToolInput(raw []byte) {
	t.ToolInput = tools.NewInputFromPayload(t.ToolName, raw, "tool_input")
}

// NativeToolName returns the tool name.
func (t ToolFields) NativeToolName() string {
	return t.ToolName
}

// Input returns tool input.
func (t ToolFields) Input() tools.Input {
	return t.ToolInput
}
