package portconfig

import (
	"encoding/json"
	"regexp"
	"strings"
	"unicode"

	"github.com/sviatsviatsviat/wat/cmd/wat/internal/portconfig/model"
	"github.com/sviatsviatsviat/wat/sdk/agnostic"
	agclaude "github.com/sviatsviatsviat/wat/sdk/agnostic/claude"
	agcopilot "github.com/sviatsviatsviat/wat/sdk/agnostic/copilot"
	agcursor "github.com/sviatsviatsviat/wat/sdk/agnostic/cursor"
)

const (
	claudeDefaultTimeoutSec  = 600
	copilotDefaultTimeoutSec = 30
)

var copilotAnchoredPattern = regexp.MustCompile(`^\^\(\?:(.+)\)\$$`)

var toolNamesFor = map[agnostic.Dialect]map[string]string{
	agnostic.Claude: {
		agnostic.ToolBash:      "Bash",
		agnostic.ToolEdit:      "Edit",
		agnostic.ToolWrite:     "Write",
		agnostic.ToolRead:      "Read",
		agnostic.ToolGlob:      "Glob",
		agnostic.ToolGrep:      "Grep",
		agnostic.ToolTask:      "Agent",
		agnostic.ToolWebFetch:  "WebFetch",
		agnostic.ToolWebSearch: "WebSearch",
	},
	agnostic.Copilot: {
		agnostic.ToolBash:     "bash",
		agnostic.ToolEdit:     "edit",
		agnostic.ToolWrite:    "create",
		agnostic.ToolRead:     "view",
		agnostic.ToolGlob:     "glob",
		agnostic.ToolGrep:     "grep",
		agnostic.ToolTask:     "task",
		agnostic.ToolWebFetch: "web_fetch",
	},
	agnostic.Cursor: {
		agnostic.ToolBash:   "Shell",
		agnostic.ToolEdit:   "Write",
		agnostic.ToolWrite:  "Write",
		agnostic.ToolRead:   "Read",
		agnostic.ToolGrep:   "Grep",
		agnostic.ToolTask:   "Task",
		agnostic.ToolDelete: "Delete",
	},
}

var knownHandlerKeys = map[agnostic.Dialect]map[string]bool{
	agnostic.Claude: {
		"type": true, "command": true, "args": true, "prompt": true, "url": true, "timeout": true, "if": true,
	},
	agnostic.Copilot: {
		"type": true, "bash": true, "powershell": true, "command": true, "url": true, "prompt": true,
		"matcher": true, "timeoutSec": true, "timeout": true, "cwd": true, "env": true,
	},
	agnostic.Cursor: {
		"command": true, "type": true, "prompt": true, "matcher": true, "timeout": true,
		"loop_limit": true, "failClosed": true,
	},
}

// Translate converts a native hook config from one agent dialect to another.
// Warnings describe lossy or unsupported mappings; unmappable hooks are omitted
// from output with an explicit warning rather than silently dropped.
func Translate(data []byte, from, to agnostic.Dialect) ([]byte, []Warning, error) {
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

func prepareForTarget(cfg *Config, from, to agnostic.Dialect) []Warning {
	if from == to {
		return nil
	}
	var warns []Warning
	timeoutWarned := false
	filtered := make(map[agnostic.Kind][]Entry, len(cfg.Hooks))
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
		warns = append(warns, model.Warnf("%s: native entry not portable to %s; not ported", extra.Event, to))
	}
	cfg.Extras = nil
	return warns
}

func adaptEntry(e *Entry, kind agnostic.Kind, from, to agnostic.Dialect, timeoutWarned *bool) ([]Warning, bool) {
	if eventForKind(to, kind) == "" {
		event := e.NativeEvent
		if event == "" {
			event = string(kind)
		}
		return []Warning{model.Warnf("%s: no %s equivalent; not ported", event, to)}, false
	}
	if w, ok := handlerSupportedOnTarget(*e, kind, to); !ok {
		return w, false
	}
	var warns []Warning
	if from == agnostic.Cursor && to != agnostic.Cursor && agcursor.IsDedicatedEvent(e.NativeEvent) {
		warns = append(warns, model.Warnf("%s: Cursor dedicated event maps to generic %s on %s; review matcher and payload semantics",
			e.NativeEvent, eventForKind(to, kind), to))
	}
	if from == agnostic.Claude {
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

func applyExplicitTimeout(e *Entry, from, to agnostic.Dialect, timeoutWarned *bool) []Warning {
	fromDefault := defaultTimeoutFor(from)
	toDefault := defaultTimeoutFor(to)
	if fromDefault == 0 || toDefault == 0 || fromDefault == toDefault {
		return nil
	}
	var warns []Warning
	if !*timeoutWarned {
		warns = append(warns, model.Warnf("unset timeout: %s default %ds vs %s default %ds; emitting explicit %ds from source",
			from, fromDefault, to, toDefault, fromDefault))
		*timeoutWarned = true
	}
	e.TimeoutSec = fromDefault
	return warns
}

func eventForKind(d agnostic.Dialect, kind agnostic.Kind) string {
	switch d {
	case agnostic.Claude:
		return agclaude.EventForKind[kind]
	case agnostic.Copilot:
		return agcopilot.EventForKind[kind]
	case agnostic.Cursor:
		return agcursor.EventForKind[kind]
	default:
		return ""
	}
}

func handlerSupportedOnTarget(e Entry, kind agnostic.Kind, to agnostic.Dialect) ([]Warning, bool) {
	switch to {
	case agnostic.Cursor:
		if e.Type == "http" {
			event := e.NativeEvent
			if event == "" {
				event = eventForKind(to, kind)
			}
			return []Warning{model.Warnf("%s: Cursor has no http hooks; not ported", event)}, false
		}
	case agnostic.Copilot:
		switch e.Type {
		case "http", "command", "", "prompt":
			if e.Type == "prompt" && kind != agnostic.KindSessionStart {
				event := e.NativeEvent
				if event == "" {
					event = eventForKind(to, kind)
				}
				return []Warning{model.Warnf("%s: Copilot supports prompt hooks only on sessionStart; not ported", event)}, false
			}
		default:
			event := e.NativeEvent
			if event == "" {
				event = eventForKind(to, kind)
			}
			return []Warning{model.Warnf("%s: unsupported handler type %q on Copilot; not ported", event, e.Type)}, false
		}
	}
	return nil, true
}

func claudeIfWarnings(groupIf, handlerRaw json.RawMessage) []Warning {
	if len(groupIf) > 0 || hasJSONField(handlerRaw, "if") {
		return []Warning{model.Warnf("Claude if permission rule has no target equivalent")}
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

func droppedRawFieldWarnings(from agnostic.Dialect, raw json.RawMessage) []Warning {
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
	return []Warning{model.Warnf("handler fields not portable to other agents: %s", strings.Join(dropped, ", "))}
}

func defaultTimeoutFor(d agnostic.Dialect) int {
	switch d {
	case agnostic.Claude:
		return claudeDefaultTimeoutSec
	case agnostic.Copilot:
		return copilotDefaultTimeoutSec
	default:
		return 0
	}
}

func kindHasToolMatcher(k agnostic.Kind) bool {
	switch k {
	case agnostic.KindPreTool, agnostic.KindPostTool, agnostic.KindPostToolFailure, agnostic.Kind("PermissionRequest"):
		return true
	default:
		return false
	}
}

func translateMatcher(matcher string, kind agnostic.Kind, from, to agnostic.Dialect) (string, []Warning) {
	if matcher == "" || matcher == "*" || !kindHasToolMatcher(kind) {
		return matcher, nil
	}
	var warns []Warning
	if from == agnostic.Copilot {
		if m := copilotAnchoredPattern.FindStringSubmatch(matcher); len(m) == 2 {
			original := matcher
			matcher = m[1]
			warns = append(warns, model.Warnf("matcher %q: Copilot anchored regex un-anchored for %s", original, to))
		}
	}
	if !isSimpleAlternation(matcher) && !isSingleSimpleToken(matcher) {
		warns = append(warns, model.Warnf("matcher %q: complex regex kept verbatim for %s", matcher, to))
		return matcher, warns
	}
	matcher, tokenWarns := translateToolTokens(matcher, to)
	warns = append(warns, tokenWarns...)
	if to == agnostic.Copilot {
		if isSimpleAlternation(matcher) {
			matcher = "^(?:" + matcher + ")$"
		} else if !copilotAnchoredPattern.MatchString(matcher) {
			warns = append(warns, model.Warnf("matcher %q: complex regex kept verbatim for Copilot; verify anchored form", matcher))
		}
	}
	return matcher, warns
}

func translateToolTokens(matcher string, to agnostic.Dialect) (string, []Warning) {
	target := toolNamesFor[to]
	var warns []Warning
	sep := "|"
	parts := splitMatcherTokens(matcher)
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		tok := strings.TrimSpace(part)
		canon, mcp := agnostic.NormalizeToolName(tok)
		if mcp {
			out = append(out, tok)
			warns = append(warns, model.Warnf("matcher %q: MCP tool pattern kept verbatim; verify %s naming", tok, to))
			continue
		}
		if native, ok := target[canon]; ok {
			out = append(out, native)
			continue
		}
		out = append(out, tok)
		warns = append(warns, model.Warnf("matcher token %q has no %s equivalent; kept verbatim", tok, to))
	}
	return strings.Join(out, sep), warns
}

func splitMatcherTokens(matcher string) []string {
	return strings.FieldsFunc(matcher, func(r rune) bool {
		return r == '|' || r == ','
	})
}

func isSingleSimpleToken(matcher string) bool {
	parts := splitMatcherTokens(matcher)
	if len(parts) != 1 {
		return false
	}
	tok := strings.TrimSpace(parts[0])
	if tok == "" {
		return false
	}
	for _, r := range tok {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' && r != '-' {
			return false
		}
	}
	return true
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
