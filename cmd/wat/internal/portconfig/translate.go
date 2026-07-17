package portconfig

import (
	"encoding/json"
	"regexp"
	"strings"
	"unicode"

	portclaude "github.com/sviatsviatsviat/wat/cmd/wat/internal/portconfig/claude"
	portcopilot "github.com/sviatsviatsviat/wat/cmd/wat/internal/portconfig/copilot"
	portcursor "github.com/sviatsviatsviat/wat/cmd/wat/internal/portconfig/cursor"
	"github.com/sviatsviatsviat/wat/cmd/wat/internal/portconfig/model"
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/agnostic"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/tools"
	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

const (
	claudeDefaultTimeoutSec  = 600
	copilotDefaultTimeoutSec = 30
)

var copilotAnchoredPattern = regexp.MustCompile(`^\^\(\?:(.+)\)\$$`)

var toolNamesFor = map[string]map[string]string{
	sdkclaude.Dialect: {
		tools.ToolBash:      "Bash",
		tools.ToolEdit:      "Edit",
		tools.ToolWrite:     "Write",
		tools.ToolRead:      "Read",
		tools.ToolGlob:      "Glob",
		tools.ToolGrep:      "Grep",
		tools.ToolTask:      "Agent",
		tools.ToolWebFetch:  "WebFetch",
		tools.ToolWebSearch: "WebSearch",
	},
	sdkcopilot.Dialect: {
		tools.ToolBash:     "bash",
		tools.ToolEdit:     "edit",
		tools.ToolWrite:    "create",
		tools.ToolRead:     "view",
		tools.ToolGlob:     "glob",
		tools.ToolGrep:     "grep",
		tools.ToolTask:     "task",
		tools.ToolWebFetch: "web_fetch",
	},
	sdkcursor.Dialect: {
		tools.ToolBash:   "Shell",
		tools.ToolEdit:   "Write",
		tools.ToolWrite:  "Write",
		tools.ToolRead:   "Read",
		tools.ToolGrep:   "Grep",
		tools.ToolTask:   "Task",
		tools.ToolDelete: "Delete",
	},
}

var knownHandlerKeys = map[string]map[string]bool{
	sdkclaude.Dialect: {
		"type": true, "command": true, "args": true, "prompt": true, "url": true, "timeout": true, "if": true,
	},
	sdkcopilot.Dialect: {
		"type": true, "bash": true, "powershell": true, "command": true, "url": true, "prompt": true,
		"matcher": true, "timeoutSec": true, "timeout": true, "cwd": true, "env": true,
	},
	sdkcursor.Dialect: {
		"command": true, "type": true, "prompt": true, "matcher": true, "timeout": true,
		"loop_limit": true, "failClosed": true,
	},
}

// Translate converts a native hook config from one agent dialect to another.
// Warnings describe lossy or unsupported mappings; unmappable hooks are omitted
// from output with an explicit warning rather than silently dropped.
func Translate(data []byte, from, to string) ([]byte, []Warning, error) {
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

func prepareForTarget(cfg *Config, from, to string) []Warning {
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

func adaptEntry(e *Entry, kind agnostic.Kind, from, to string, timeoutWarned *bool) ([]Warning, bool) {
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
	if from == sdkcursor.Dialect && to != sdkcursor.Dialect && portcursor.IsDedicatedEvent(e.NativeEvent) {
		warns = append(warns, model.Warnf("%s: Cursor dedicated event maps to generic %s on %s; review matcher and payload semantics",
			e.NativeEvent, eventForKind(to, kind), to))
	}
	if from == sdkclaude.Dialect {
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

func applyExplicitTimeout(e *Entry, from, to string, timeoutWarned *bool) []Warning {
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

func eventForKind(d string, kind agnostic.Kind) string {
	switch d {
	case sdkclaude.Dialect:
		return portclaude.EventForKind[kind]
	case sdkcopilot.Dialect:
		return portcopilot.EventForKind[kind]
	case sdkcursor.Dialect:
		return portcursor.EventForKind[kind]
	default:
		return ""
	}
}

func handlerSupportedOnTarget(e Entry, kind agnostic.Kind, to string) ([]Warning, bool) {
	switch to {
	case sdkcursor.Dialect:
		if e.Type == "http" {
			event := e.NativeEvent
			if event == "" {
				event = eventForKind(to, kind)
			}
			return []Warning{model.Warnf("%s: Cursor has no http hooks; not ported", event)}, false
		}
	case sdkcopilot.Dialect:
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

func droppedRawFieldWarnings(from string, raw json.RawMessage) []Warning {
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

func defaultTimeoutFor(d string) int {
	switch d {
	case sdkclaude.Dialect:
		return claudeDefaultTimeoutSec
	case sdkcopilot.Dialect:
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

func translateMatcher(matcher string, kind agnostic.Kind, from, to string) (string, []Warning) {
	if matcher == "" || matcher == "*" || !kindHasToolMatcher(kind) {
		return matcher, nil
	}
	var warns []Warning
	if from == sdkcopilot.Dialect {
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
	if to == sdkcopilot.Dialect {
		if isSimpleAlternation(matcher) {
			matcher = "^(?:" + matcher + ")$"
		} else if !copilotAnchoredPattern.MatchString(matcher) {
			warns = append(warns, model.Warnf("matcher %q: complex regex kept verbatim for Copilot; verify anchored form", matcher))
		}
	}
	return matcher, warns
}

func translateToolTokens(matcher string, to string) (string, []Warning) {
	target := toolNamesFor[to]
	var warns []Warning
	sep := "|"
	parts := splitMatcherTokens(matcher)
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		tok := strings.TrimSpace(part)
		canon, mcp := hookkit.NormalizeToolName(tok)
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
