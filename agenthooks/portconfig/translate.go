package portconfig

import (
	"encoding/json"
	"regexp"
	"strings"
	"unicode"

	"github.com/sviatsviatsviat/wat/agenthooks"
)

const (
	claudeDefaultTimeoutSec  = 600
	copilotDefaultTimeoutSec = 30
)

var copilotAnchoredPattern = regexp.MustCompile(`^\^\(\?:(.+)\)\$$`)

var toolNamesFor = map[agenthooks.Dialect]map[string]string{
	agenthooks.Claude: {
		agenthooks.ToolBash:      "Bash",
		agenthooks.ToolEdit:      "Edit",
		agenthooks.ToolWrite:     "Write",
		agenthooks.ToolRead:      "Read",
		agenthooks.ToolGlob:      "Glob",
		agenthooks.ToolGrep:      "Grep",
		agenthooks.ToolTask:      "Agent",
		agenthooks.ToolWebFetch:  "WebFetch",
		agenthooks.ToolWebSearch: "WebSearch",
	},
	agenthooks.Copilot: {
		agenthooks.ToolBash:     "bash",
		agenthooks.ToolEdit:     "edit",
		agenthooks.ToolWrite:    "create",
		agenthooks.ToolRead:     "view",
		agenthooks.ToolGlob:     "glob",
		agenthooks.ToolGrep:     "grep",
		agenthooks.ToolTask:     "task",
		agenthooks.ToolWebFetch: "web_fetch",
	},
	agenthooks.Cursor: {
		agenthooks.ToolBash:   "Shell",
		agenthooks.ToolEdit:   "Write",
		agenthooks.ToolWrite:  "Write",
		agenthooks.ToolRead:   "Read",
		agenthooks.ToolGrep:   "Grep",
		agenthooks.ToolTask:   "Task",
		agenthooks.ToolDelete: "Delete",
	},
}

var knownHandlerKeys = map[agenthooks.Dialect]map[string]bool{
	agenthooks.Claude: {
		"type": true, "command": true, "args": true, "prompt": true, "url": true, "timeout": true, "if": true,
	},
	agenthooks.Copilot: {
		"type": true, "bash": true, "powershell": true, "command": true, "url": true, "prompt": true,
		"matcher": true, "timeoutSec": true, "timeout": true, "cwd": true, "env": true,
	},
	agenthooks.Cursor: {
		"command": true, "type": true, "prompt": true, "matcher": true, "timeout": true,
		"loop_limit": true, "failClosed": true,
	},
}

// Translate converts a native hook config from one agent dialect to another.
// Warnings describe lossy or unsupported mappings; unmappable hooks are omitted
// from output with an explicit warning rather than silently dropped.
func Translate(data []byte, from, to agenthooks.Dialect) ([]byte, []Warning, error) {
	if from == to {
		return data, nil, nil
	}
	cfg, warns, err := Parse(data, from)
	if err != nil {
		return nil, warns, err
	}
	adaptWarns := prepareForTarget(&cfg, from, to)
	out, emitWarns, err := Emit(cfg, to)
	return out, append(append(warns, adaptWarns...), emitWarns...), err
}

func prepareForTarget(cfg *Config, from, to agenthooks.Dialect) []Warning {
	if from == to {
		return nil
	}
	var warns []Warning
	timeoutWarned := false
	filtered := make(map[agenthooks.Kind][]Entry, len(cfg.Hooks))
	for kind, entries := range cfg.Hooks {
		for _, e := range entries {
			entryWarns, keep := adaptEntry(&e, kind, from, to, &timeoutWarned)
			warns = append(warns, entryWarns...)
			if keep {
				filtered[kind] = append(filtered[kind], e)
			}
		}
	}
	cfg.Hooks = filtered

	for _, extra := range cfg.Extras {
		warns = append(warns, warnf("%s: native entry not portable to %s; not ported", extra.Event, to))
	}
	cfg.Extras = nil
	return warns
}

func adaptEntry(e *Entry, kind agenthooks.Kind, from, to agenthooks.Dialect, timeoutWarned *bool) ([]Warning, bool) {
	if eventForKind(to, kind) == "" {
		event := e.NativeEvent
		if event == "" {
			event = string(kind)
		}
		return []Warning{warnf("%s: no %s equivalent; not ported", event, to)}, false
	}
	if w, ok := handlerSupportedOnTarget(*e, kind, to); !ok {
		return w, false
	}
	var warns []Warning
	if from == agenthooks.Cursor && to != agenthooks.Cursor && isCursorDedicatedEvent(e.NativeEvent) {
		warns = append(warns, warnf("%s: Cursor dedicated event maps to generic %s on %s; review matcher and payload semantics",
			e.NativeEvent, eventForKind(to, kind), to))
	}
	if from == agenthooks.Claude {
		warns = append(warns, claudeIfWarnings(e.ClaudeGroupIf, e.Raw)...)
	}
	warns = append(warns, droppedRawFieldWarnings(from, e.Raw)...)
	if e.TimeoutSec == 0 {
		warns = append(warns, applyExplicitTimeout(e, from, to, timeoutWarned)...)
	}
	matcher, matcherWarns := translateMatcher(e.Matcher, kind, from, to)
	warns = append(warns, matcherWarns...)
	e.Matcher = matcher
	e.NativeEvent = ""
	e.ClaudeGroupIf = nil
	e.Raw = nil
	return warns, true
}

func applyExplicitTimeout(e *Entry, from, to agenthooks.Dialect, timeoutWarned *bool) []Warning {
	fromDefault := defaultTimeoutFor(from)
	toDefault := defaultTimeoutFor(to)
	if fromDefault == 0 || toDefault == 0 || fromDefault == toDefault {
		return nil
	}
	var warns []Warning
	if !*timeoutWarned {
		warns = append(warns, warnf("unset timeout: %s default %ds vs %s default %ds; emitting explicit %ds from source",
			from, fromDefault, to, toDefault, fromDefault))
		*timeoutWarned = true
	}
	e.TimeoutSec = fromDefault
	return warns
}

func eventForKind(d agenthooks.Dialect, kind agenthooks.Kind) string {
	switch d {
	case agenthooks.Claude:
		return agenthooks.ClaudeEventForKind[kind]
	case agenthooks.Copilot:
		return copilotEventForKind[kind]
	case agenthooks.Cursor:
		return cursorEventForKind[kind]
	default:
		return ""
	}
}

func handlerSupportedOnTarget(e Entry, kind agenthooks.Kind, to agenthooks.Dialect) ([]Warning, bool) {
	switch to {
	case agenthooks.Cursor:
		if e.Type == "http" {
			event := e.NativeEvent
			if event == "" {
				event = eventForKind(to, kind)
			}
			return []Warning{warnf("%s: Cursor has no http hooks; not ported", event)}, false
		}
	case agenthooks.Copilot:
		switch e.Type {
		case "http", "command", "", "prompt":
			if e.Type == "prompt" && kind != agenthooks.KindSessionStart {
				event := e.NativeEvent
				if event == "" {
					event = eventForKind(to, kind)
				}
				return []Warning{warnf("%s: Copilot supports prompt hooks only on sessionStart; not ported", event)}, false
			}
		default:
			event := e.NativeEvent
			if event == "" {
				event = eventForKind(to, kind)
			}
			return []Warning{warnf("%s: unsupported handler type %q on Copilot; not ported", event, e.Type)}, false
		}
	}
	return nil, true
}

func claudeIfWarnings(groupIf, handlerRaw json.RawMessage) []Warning {
	if len(groupIf) > 0 || hasJSONField(handlerRaw, "if") {
		return []Warning{warnf("Claude if permission rule has no target equivalent")}
	}
	return nil
}

func hasJSONField(raw json.RawMessage, key string) bool {
	if len(raw) == 0 {
		return false
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(raw, &fields); err != nil {
		return false
	}
	_, ok := fields[key]
	return ok
}

func droppedRawFieldWarnings(from agenthooks.Dialect, raw json.RawMessage) []Warning {
	if len(raw) == 0 {
		return nil
	}
	known := knownHandlerKeys[from]
	if len(known) == 0 {
		return nil
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(raw, &fields); err != nil {
		return nil
	}
	var dropped []string
	for key := range fields {
		if !known[key] {
			dropped = append(dropped, key)
		}
	}
	if len(dropped) == 0 {
		return nil
	}
	return []Warning{warnf("handler fields not portable to other agents: %s", strings.Join(dropped, ", "))}
}

func defaultTimeoutFor(d agenthooks.Dialect) int {
	switch d {
	case agenthooks.Claude:
		return claudeDefaultTimeoutSec
	case agenthooks.Copilot:
		return copilotDefaultTimeoutSec
	default:
		return 0
	}
}

func kindHasToolMatcher(k agenthooks.Kind) bool {
	switch k {
	case agenthooks.KindPreTool, agenthooks.KindPostTool, agenthooks.KindPostToolFailure, agenthooks.KindPermissionRequest:
		return true
	default:
		return false
	}
}

func translateMatcher(matcher string, kind agenthooks.Kind, from, to agenthooks.Dialect) (string, []Warning) {
	if matcher == "" || matcher == "*" || !kindHasToolMatcher(kind) {
		return matcher, nil
	}
	var warns []Warning
	if from == agenthooks.Copilot {
		if m := copilotAnchoredPattern.FindStringSubmatch(matcher); len(m) == 2 {
			original := matcher
			matcher = m[1]
			warns = append(warns, warnf("matcher %q: Copilot anchored regex un-anchored for %s", original, to))
		}
	}
	matcher, tokenWarns := translateToolTokens(matcher, to)
	warns = append(warns, tokenWarns...)
	if to == agenthooks.Copilot {
		if isSimpleAlternation(matcher) {
			matcher = "^(?:" + matcher + ")$"
		} else if !copilotAnchoredPattern.MatchString(matcher) {
			warns = append(warns, warnf("matcher %q: complex regex kept verbatim for Copilot; verify anchored form", matcher))
		}
	}
	return matcher, warns
}

func translateToolTokens(matcher string, to agenthooks.Dialect) (string, []Warning) {
	target := toolNamesFor[to]
	var warns []Warning
	sep := "|"
	parts := splitMatcherTokens(matcher)
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		tok := strings.TrimSpace(part)
		canon, mcp := agenthooks.NormalizeToolName(tok)
		if mcp {
			out = append(out, tok)
			warns = append(warns, warnf("matcher %q: MCP tool pattern kept verbatim; verify %s naming", tok, to))
			continue
		}
		if native, ok := target[canon]; ok {
			out = append(out, native)
			continue
		}
		out = append(out, tok)
		if canon != tok {
			warns = append(warns, warnf("matcher token %q has no %s equivalent; kept verbatim", tok, to))
		}
	}
	return strings.Join(out, sep), warns
}

func splitMatcherTokens(matcher string) []string {
	return strings.FieldsFunc(matcher, func(r rune) bool {
		return r == '|' || r == ','
	})
}

func isSimpleAlternation(matcher string) bool {
	if matcher == "" {
		return false
	}
	for _, part := range splitMatcherTokens(matcher) {
		tok := strings.TrimSpace(part)
		if tok == "" {
			return false
		}
		for _, r := range tok {
			if r == '|' || r == ',' {
				return false
			}
			if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' && r != '-' {
				return false
			}
		}
	}
	return true
}
