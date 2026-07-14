package copilot

import (
	"encoding/json"
	"fmt"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
)

// Encode renders a typed output struct as Copilot flat camelCase stdout JSON and
// returns the process exit code.
func Encode(eventName string, out any) ([]byte, int, error) {
	out = hookkit.NormalizeOutput(out)
	if out == nil || isZeroOutput(out) {
		return nil, 0, nil
	}
	if err := validateEncodePair(eventName, out); err != nil {
		return nil, 0, err
	}

	switch o := out.(type) {
	case PreToolOutput:
		return encodePreTool(o)
	case PostToolOutput:
		return encodePostTool(o)
	case StopOutput:
		return encodeStop(o)
	case PermissionRequestOutput:
		return encodePermissionRequest(o)
	case PostToolFailureOutput:
		return encodePostToolFailure(o)
	case SessionStartOutput:
		return encodeAdditionalContext(o.AdditionalContext)
	case SubagentStartOutput:
		return encodeAdditionalContext(o.AdditionalContext)
	case NotificationOutput:
		return encodeAdditionalContext(o.AdditionalContext)
	default:
		return nil, 0, fmt.Errorf("copilot: encode: unsupported output type %T", out)
	}
}

func isZeroOutput(out any) bool {
	if z, ok := out.(interface{ isZero() bool }); ok {
		return z.isZero()
	}
	return hookkit.IsZeroOutput(out)
}

func validateEncodePair(eventName string, out any) error {
	allowed, ok := allowedEventsForOutput(out)
	if !ok {
		return fmt.Errorf("copilot: encode: unsupported output type %T", out)
	}
	return hookkit.ValidateEncodePair("copilot", eventName, out, allowed, func(name string) (string, bool) {
		canonical, known := CanonicalEventName(name)
		if !known {
			return name, true
		}
		return canonical, true
	})
}

func allowedEventsForOutput(out any) ([]string, bool) {
	switch out.(type) {
	case PreToolOutput:
		return []string{EventPreToolUse}, true
	case PostToolOutput:
		return []string{EventPostToolUse}, true
	case StopOutput:
		return []string{EventAgentStop, EventSubagentStop}, true
	case PermissionRequestOutput:
		return []string{EventPermissionRequest}, true
	case PostToolFailureOutput:
		return []string{EventPostToolUseFailure}, true
	case SessionStartOutput:
		return []string{EventSessionStart}, true
	case SubagentStartOutput:
		return []string{EventSubagentStart}, true
	case NotificationOutput:
		return []string{EventNotification}, true
	default:
		return nil, false
	}
}

func encodePreTool(o PreToolOutput) ([]byte, int, error) {
	out := map[string]any{}
	if o.Decision != "" {
		out["permissionDecision"] = string(o.Decision)
		if o.Reason != "" {
			out["permissionDecisionReason"] = o.Reason
		}
	}
	if o.ModifiedArgs != nil {
		out["modifiedArgs"] = o.ModifiedArgs
	}
	if len(out) == 0 {
		return nil, 0, nil
	}
	b, err := json.Marshal(out)
	return b, 0, err
}

func encodePostTool(o PostToolOutput) ([]byte, int, error) {
	out := map[string]any{}
	if o.ModifiedResult != "" {
		out["modifiedResult"] = map[string]any{
			"resultType":       "success",
			"textResultForLlm": o.ModifiedResult,
		}
	}
	if o.AdditionalContext != "" {
		out["additionalContext"] = o.AdditionalContext
	}
	if len(out) == 0 {
		return nil, 0, nil
	}
	b, err := json.Marshal(out)
	return b, 0, err
}

func encodeStop(o StopOutput) ([]byte, int, error) {
	if o.Reason == "" {
		return nil, 0, nil
	}
	out := map[string]any{
		"decision": "block",
		"reason":   o.Reason,
	}
	b, err := json.Marshal(out)
	return b, 0, err
}

func encodePermissionRequest(o PermissionRequestOutput) ([]byte, int, error) {
	if o.Behavior == "" && o.Message == "" && !o.Interrupt {
		return nil, 0, nil
	}
	out := map[string]any{}
	if o.Behavior != "" {
		out["behavior"] = o.Behavior
	}
	if o.Message != "" {
		out["message"] = o.Message
	}
	if o.Interrupt {
		out["interrupt"] = true
	}
	b, err := json.Marshal(out)
	if err != nil {
		return nil, 0, err
	}
	exitCode := 0
	if o.Behavior == "deny" && !o.SuppressWarnExit {
		exitCode = WarnExit
	}
	return b, exitCode, err
}

func encodePostToolFailure(o PostToolFailureOutput) ([]byte, int, error) {
	if o.Context == "" {
		return nil, 0, nil
	}
	return []byte(o.Context), WarnExit, nil
}

func encodeAdditionalContext(context string) ([]byte, int, error) {
	if context == "" {
		return nil, 0, nil
	}
	out := map[string]any{"additionalContext": context}
	b, err := json.Marshal(out)
	return b, 0, err
}
