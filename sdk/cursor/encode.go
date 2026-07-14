package cursor

import (
	"fmt"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
)

// Encode renders a typed output struct as Cursor stdout JSON and returns the
// process exit code.
func Encode(eventName string, out any) ([]byte, int, error) {
	out = hookkit.NormalizeOutput(out)
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
		return fmt.Errorf("cursor: encode: unsupported output type %T", out)
	}
	return hookkit.ValidateEncodePair("cursor", eventName, out, allowed, nil)
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
