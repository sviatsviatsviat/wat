package portconfig

import (
	"github.com/sviatsviatsviat/wat/agenthooks"
)

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

var cursorKindForEvent = buildCursorKindForEvent()

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
