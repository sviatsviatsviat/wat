package copilot

import (
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

// IsZeroOutput reports whether out is an empty hook response.
func IsZeroOutput(out any) bool { return isZeroOutput(out) }

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
