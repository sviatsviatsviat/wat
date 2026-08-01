package event

import (
	"fmt"
	"strings"

	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/tools"
)

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
// It returns an error when tool_name is missing (for example camelCase toolName
// left the snake_case field empty).
func (t *ToolFields) BindToolInput(raw []byte) error {
	if strings.TrimSpace(t.ToolName) == "" {
		return fmt.Errorf("tool_name is required (use snake_case tool_name/tool_input; camelCase toolName/toolArgs are unsupported)")
	}
	t.ToolInput = tools.NewInputFromPayload(t.ToolName, raw, "tool_input")
	return nil
}

// NativeToolName returns the tool name.
func (t ToolFields) NativeToolName() string {
	return t.ToolName
}

// Input returns tool input.
func (t ToolFields) Input() tools.Input {
	return t.ToolInput
}
