package cursor

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// Encode renders a typed output struct as Cursor stdout JSON and returns the
// process exit code.
func Encode(eventName string, out any) ([]byte, int, error) {
	out = normalizeOutput(out)
	if out == nil || isZeroOutput(out) {
		return nil, 0, nil
	}
	if err := validateEncodePair(eventName, out); err != nil {
		return nil, 0, err
	}

	switch o := out.(type) {
	case PermissionOutput:
		return encodePermission(eventName, o)
	case BeforeSubmitPromptOutput:
		return encodeBeforeSubmitPrompt(o)
	case PostToolOutput:
		return encodePostTool(o)
	case StopOutput:
		return encodeStop(o)
	case SessionStartOutput:
		return encodeSessionStart(o)
	case PreCompactOutput:
		return encodePreCompact(o)
	default:
		return nil, 0, fmt.Errorf("cursor: encode: unsupported output type %T", out)
	}
}

func normalizeOutput(out any) any {
	if out == nil {
		return nil
	}
	v := reflect.ValueOf(out)
	if v.Kind() != reflect.Pointer {
		return out
	}
	if v.IsNil() {
		return nil
	}
	return v.Elem().Interface()
}

func isZeroOutput(out any) bool {
	if z, ok := out.(interface{ isZero() bool }); ok {
		return z.isZero()
	}
	return reflect.ValueOf(out).IsZero()
}

func validateEncodePair(eventName string, out any) error {
	if eventName == "" {
		return nil
	}
	allowed, ok := allowedEventsForOutput(out)
	if !ok {
		return fmt.Errorf("cursor: encode: unsupported output type %T", out)
	}
	for _, name := range allowed {
		if eventName == name {
			return nil
		}
	}
	return fmt.Errorf("cursor: encode: event %q incompatible with output type %T", eventName, out)
}

func allowedEventsForOutput(out any) ([]string, bool) {
	switch out.(type) {
	case PermissionOutput:
		return []string{
			EventPreToolUse,
			EventBeforeShellExecution,
			EventBeforeMCPExecution,
			EventBeforeReadFile,
			EventSubagentStart,
			EventBeforeTabFileRead,
		}, true
	case BeforeSubmitPromptOutput:
		return []string{EventBeforeSubmitPrompt}, true
	case PostToolOutput:
		return []string{
			EventPostToolUse,
			EventAfterMCPExecution,
			EventAfterShellExecution,
			EventAfterFileEdit,
		}, true
	case StopOutput:
		return []string{EventStop, EventSubagentStop}, true
	case SessionStartOutput:
		return []string{EventSessionStart}, true
	case PreCompactOutput:
		return []string{EventPreCompact}, true
	default:
		return nil, false
	}
}

func encodePermission(eventName string, o PermissionOutput) ([]byte, int, error) {
	out := map[string]any{}
	if o.Decision != "" {
		out["permission"] = string(o.Decision)
	}
	if o.UserMessage != "" {
		out["user_message"] = o.UserMessage
	}
	if o.AgentMessage != "" {
		out["agent_message"] = o.AgentMessage
	}
	if o.UpdatedInput != nil && (eventName == "" || eventName == EventPreToolUse) {
		out["updated_input"] = o.UpdatedInput
	}
	if len(out) == 0 {
		return nil, 0, nil
	}
	b, err := json.Marshal(out)
	if err != nil {
		return nil, 0, err
	}
	exitCode := 0
	if o.Decision == DecisionDeny {
		exitCode = PermissionDenyExit
	}
	return b, exitCode, nil
}

func encodeBeforeSubmitPrompt(o BeforeSubmitPromptOutput) ([]byte, int, error) {
	out := map[string]any{}
	if o.Continue != nil {
		out["continue"] = *o.Continue
	}
	if o.UserMessage != "" {
		out["user_message"] = o.UserMessage
	}
	if len(out) == 0 {
		return nil, 0, nil
	}
	b, err := json.Marshal(out)
	return b, 0, err
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

func encodeStop(o StopOutput) ([]byte, int, error) {
	if o.FollowUpMessage == "" {
		return nil, 0, nil
	}
	out := map[string]any{"followup_message": o.FollowUpMessage}
	b, err := json.Marshal(out)
	return b, 0, err
}

func encodeSessionStart(o SessionStartOutput) ([]byte, int, error) {
	out := map[string]any{}
	if len(o.Env) > 0 {
		out["env"] = o.Env
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

func encodePreCompact(o PreCompactOutput) ([]byte, int, error) {
	if o.UserMessage == "" {
		return nil, 0, nil
	}
	out := map[string]any{"user_message": o.UserMessage}
	b, err := json.Marshal(out)
	return b, 0, err
}
