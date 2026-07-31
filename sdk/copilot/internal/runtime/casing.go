package runtime

import (
	"fmt"
	"strings"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
)

// camelEventHints maps Copilot CLI camelCase event ids to wat PascalCase wire
// names. Used only for clearer reject messages — not for decode aliases.
var camelEventHints = map[string]string{
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

// RegisterCasingRejects registers camelCase event names as explicit decode
// failures so mistaken CLI samples fail with a PascalCase hint instead of a
// bare "unknown hook event" message.
func RegisterCasingRejects(c *hookkit.Codec) {
	if c == nil {
		return
	}
	for camel, pascal := range camelEventHints {
		camel, pascal := camel, pascal
		c.Register(camel, func([]byte) (hookkit.Event, error) {
			return nil, fmt.Errorf("unknown hook event %s (use PascalCase %q, not camelCase %q)", camel, pascal, camel)
		})
	}
}

// AdviseMissingHandlers rejects empty-handler dispatch when the event name is a
// known Copilot CLI camelCase (or case-folded) id. Exact PascalCase unknowns stay
// successful no-ops so unhandled hosted events remain silent.
func AdviseMissingHandlers(eventName string) error {
	if pascal, ok := camelEventHints[eventName]; ok {
		return fmt.Errorf("no handlers for %q; wat expects PascalCase event %q (not camelCase). See docs/agents/copilot.md", eventName, pascal)
	}
	if pascal := caseFoldPascal(eventName); pascal != "" {
		return fmt.Errorf("no handlers for %q; wat expects PascalCase event %q. See docs/agents/copilot.md", eventName, pascal)
	}
	return nil
}

func caseFoldPascal(eventName string) string {
	if eventName == "" {
		return ""
	}
	for camel, pascal := range camelEventHints {
		if eventName == pascal {
			return ""
		}
		if strings.EqualFold(eventName, camel) || strings.EqualFold(eventName, pascal) {
			return pascal
		}
	}
	return ""
}
