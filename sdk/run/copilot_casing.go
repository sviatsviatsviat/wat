package run

import (
	"fmt"
	"io"
	"strings"
)

// copilotCamelEventHints maps Copilot CLI camelCase event ids to wat PascalCase
// wire names. Used only for clearer reject messages — not for decode aliases.
var copilotCamelEventHints = map[string]string{
	"sessionStart":          "SessionStart",
	"sessionEnd":            "SessionEnd",
	"userPromptSubmitted":   "UserPromptSubmit",
	"userPromptTransformed": "UserPromptTransformed",
	"preToolUse":            "PreToolUse",
	"postToolUse":           "PostToolUse",
	"postToolUseFailure":    "PostToolUseFailure",
	"permissionRequest":     "PermissionRequest",
	"subagentStart":         "SubagentStart",
	"subagentStop":          "SubagentStop",
	"agentStop":             "Stop",
	"preCompact":            "PreCompact",
	"notification":          "Notification",
	"errorOccurred":         "ErrorOccurred",
}

func copilotPascalHint(eventName string) (pascal string, ok bool) {
	pascal, ok = copilotCamelEventHints[eventName]
	return pascal, ok
}

// copilotCaseFoldPascal returns the PascalCase wire name when eventName matches a
// known Copilot event ignoring case but is not already exact (for example
// "pretooluse" or "PRETOOLUSE" → "PreToolUse").
func copilotCaseFoldPascal(eventName string) string {
	if eventName == "" {
		return ""
	}
	for camel, pascal := range copilotCamelEventHints {
		if eventName == pascal {
			return ""
		}
		if strings.EqualFold(eventName, camel) || strings.EqualFold(eventName, pascal) {
			return pascal
		}
	}
	return ""
}

func writeCopilotCasingReject(errw io.Writer, eventName, pascal string) {
	_, _ = fmt.Fprintf(errw, "run: copilot: no handlers for %q; wat expects PascalCase event %q (not camelCase). See docs/agents/copilot.md\n", eventName, pascal)
}

func augmentCopilotUnknownEvent(err error, eventName string) error {
	if err == nil {
		return nil
	}
	if !strings.Contains(err.Error(), "unknown hook event") {
		return err
	}
	if pascal, ok := copilotPascalHint(eventName); ok {
		return fmt.Errorf("%w (use PascalCase %q, not camelCase %q)", err, pascal, eventName)
	}
	if pascal := copilotCaseFoldPascal(eventName); pascal != "" {
		return fmt.Errorf("%w (use PascalCase %q)", err, pascal)
	}
	return fmt.Errorf("%w (wat Copilot dialect uses PascalCase event names; see docs/agents/copilot.md)", err)
}
