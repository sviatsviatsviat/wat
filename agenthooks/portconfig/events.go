package portconfig

import "github.com/sviatsviatsviat/wat/agenthooks"

var claudeEventForKind = map[agenthooks.Kind]string{
	agenthooks.KindSessionStart:      "SessionStart",
	agenthooks.KindSessionEnd:        "SessionEnd",
	agenthooks.KindUserPrompt:        "UserPromptSubmit",
	agenthooks.KindPreTool:           "PreToolUse",
	agenthooks.KindPostTool:          "PostToolUse",
	agenthooks.KindPostToolFailure:   "PostToolUseFailure",
	agenthooks.KindPermissionRequest: "PermissionRequest",
	agenthooks.KindSubagentStart:     "SubagentStart",
	agenthooks.KindSubagentStop:      "SubagentStop",
	agenthooks.KindStop:              "Stop",
	agenthooks.KindPreCompact:        "PreCompact",
	agenthooks.KindNotification:      "Notification",
	agenthooks.KindAgentError:        "StopFailure",
}

var copilotEventForKind = map[agenthooks.Kind]string{
	agenthooks.KindSessionStart:      "sessionStart",
	agenthooks.KindSessionEnd:        "sessionEnd",
	agenthooks.KindUserPrompt:        "userPromptSubmitted",
	agenthooks.KindPreTool:           "preToolUse",
	agenthooks.KindPostTool:          "postToolUse",
	agenthooks.KindPostToolFailure:   "postToolUseFailure",
	agenthooks.KindPermissionRequest: "permissionRequest",
	agenthooks.KindSubagentStart:     "subagentStart",
	agenthooks.KindSubagentStop:      "subagentStop",
	agenthooks.KindStop:              "agentStop",
	agenthooks.KindPreCompact:        "preCompact",
	agenthooks.KindNotification:      "notification",
	agenthooks.KindAgentError:        "errorOccurred",
}

var cursorEventForKind = map[agenthooks.Kind]string{
	agenthooks.KindSessionStart:    "sessionStart",
	agenthooks.KindSessionEnd:      "sessionEnd",
	agenthooks.KindUserPrompt:      "beforeSubmitPrompt",
	agenthooks.KindPreTool:         "preToolUse",
	agenthooks.KindPostTool:        "postToolUse",
	agenthooks.KindPostToolFailure: "postToolUseFailure",
	agenthooks.KindSubagentStart:   "subagentStart",
	agenthooks.KindSubagentStop:    "subagentStop",
	agenthooks.KindStop:            "stop",
	agenthooks.KindPreCompact:      "preCompact",
}

func invertKindEvent(m map[agenthooks.Kind]string) map[string]agenthooks.Kind {
	out := make(map[string]agenthooks.Kind, len(m))
	for k, v := range m {
		out[v] = k
	}
	return out
}

var (
	claudeKindForEvent  = invertKindEvent(claudeEventForKind)
	copilotKindForEvent = buildCopilotKindForEvent()
	cursorKindForEvent  = buildCursorKindForEvent()
)

func buildCopilotKindForEvent() map[string]agenthooks.Kind {
	m := invertKindEvent(copilotEventForKind)
	// Copilot configs may use VS Code PascalCase event keys.
	for k, v := range map[string]string{
		"SessionStart":       "sessionStart",
		"SessionEnd":         "sessionEnd",
		"UserPromptSubmit":   "userPromptSubmitted",
		"PreToolUse":         "preToolUse",
		"PostToolUse":        "postToolUse",
		"PostToolUseFailure": "postToolUseFailure",
		"PermissionRequest":  "permissionRequest",
		"SubagentStart":      "subagentStart",
		"SubagentStop":       "subagentStop",
		"Stop":               "agentStop",
		"PreCompact":         "preCompact",
		"Notification":       "notification",
		"ErrorOccurred":      "errorOccurred",
	} {
		if kind, ok := m[v]; ok {
			m[k] = kind
		}
	}
	return m
}

func buildCursorKindForEvent() map[string]agenthooks.Kind {
	m := invertKindEvent(cursorEventForKind)
	for event, kind := range map[string]agenthooks.Kind{
		"beforeShellExecution": agenthooks.KindPreTool,
		"afterShellExecution":  agenthooks.KindPostTool,
		"beforeMCPExecution":   agenthooks.KindPreTool,
		"afterMCPExecution":    agenthooks.KindPostTool,
		"beforeReadFile":       agenthooks.KindPreTool,
		"afterFileEdit":        agenthooks.KindPostTool,
		"afterAgentResponse":   agenthooks.KindOther,
		"afterAgentThought":    agenthooks.KindOther,
		"beforeTabFileRead":    agenthooks.KindOther,
		"afterTabFileEdit":     agenthooks.KindOther,
		"workspaceOpen":        agenthooks.KindOther,
	} {
		m[event] = kind
	}
	return m
}

// cursorDedicatedEvents lists Cursor hook events folded into unified kinds but
// not portable as dedicated events on other agents.
var cursorDedicatedEvents = map[string]bool{
	"beforeShellExecution": true,
	"afterShellExecution":  true,
	"beforeMCPExecution":   true,
	"afterMCPExecution":    true,
	"beforeReadFile":       true,
	"afterFileEdit":        true,
}

func isCursorDedicatedEvent(event string) bool {
	return cursorDedicatedEvents[event]
}

func eventNameForEmit(e Entry, kindForEvent map[string]agenthooks.Kind, eventForKind map[agenthooks.Kind]string) string {
	if e.NativeEvent != "" {
		if k, ok := kindForEvent[e.NativeEvent]; ok && k == e.Kind {
			return e.NativeEvent
		}
	}
	return eventForKind[e.Kind]
}
