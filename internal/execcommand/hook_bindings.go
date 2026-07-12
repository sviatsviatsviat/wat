package execcommand

import (
	"strconv"

	"github.com/sviatsviatsviat/wat/internal/core"
	"github.com/sviatsviatsviat/wat/internal/helpers"
)

// templateBindingsAfterFileEdit maps __KEY__ placeholders for [core.AfterFileEditHook] payloads.
type templateBindingsAfterFileEdit struct {
	hook core.AfterFileEditHook
}

func (b templateBindingsAfterFileEdit) TemplateValue(placeholderKey string) (string, bool) {
	switch placeholderKey {
	case "HOOK_EVENT_NAME":
		return b.hook.HookEventName(), true
	case "TRANSCRIPT_PATH":
		return helpers.StringFromPtr(b.hook.TranscriptPath()), true
	case "FILE_PATH":
		return b.hook.FilePath(), true
	default:
		return "", false
	}
}

// templateBindingsAfterShellExecution maps __KEY__ placeholders for [core.AfterShellExecutionHook] payloads.
type templateBindingsAfterShellExecution struct {
	hook core.AfterShellExecutionHook
}

func (b templateBindingsAfterShellExecution) TemplateValue(placeholderKey string) (string, bool) {
	switch placeholderKey {
	case "HOOK_EVENT_NAME":
		return b.hook.HookEventName(), true
	case "TRANSCRIPT_PATH":
		return helpers.StringFromPtr(b.hook.TranscriptPath()), true
	case "DURATION":
		return strconv.FormatFloat(float64(b.hook.Duration()), 'f', -1, 32), true
	case "SANDBOX":
		return strconv.FormatBool(b.hook.Sandbox()), true
	case "COMMAND":
		return b.hook.RawCommand(), true
	default:
		return "", false
	}
}
