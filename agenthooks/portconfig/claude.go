package portconfig

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sviatsviatsviat/wat/agenthooks"
)

type claudeFile struct {
	Hooks map[string][]claudeGroup `json:"hooks"`
}

type claudeGroup struct {
	Matcher string            `json:"matcher,omitempty"`
	If      json.RawMessage   `json:"if,omitempty"`
	Hooks   []json.RawMessage `json:"hooks"`
}

type claudeHandler struct {
	Type    string   `json:"type"`
	Command string   `json:"command,omitempty"`
	Args    []string `json:"args,omitempty"`
	Prompt  string   `json:"prompt,omitempty"`
	URL     string   `json:"url,omitempty"`
	Timeout int      `json:"timeout,omitempty"`
}

func parseClaude(data []byte) (Config, []Warning, error) {
	var f claudeFile
	if err := json.Unmarshal(data, &f); err != nil {
		return Config{}, nil, fmt.Errorf("portconfig: parse claude settings: %w", err)
	}
	cfg := Config{}
	var warns []Warning
	for event, groups := range f.Hooks {
		kind, known := claudeKindForEvent[event]
		if !known {
			for _, g := range groups {
				raw, err := json.Marshal(g)
				if err != nil {
					return Config{}, warns, fmt.Errorf("portconfig: marshal claude extra: %w", err)
				}
				appendExtra(&cfg, event, raw)
			}
			continue
		}
		for _, g := range groups {
			for _, handlerRaw := range g.Hooks {
				entry, extraRaw, w, ok := claudeHandlerToEntry(event, kind, g.Matcher, g.If, handlerRaw)
				warns = append(warns, w...)
				if !ok {
					if extraRaw != nil {
						appendExtra(&cfg, event, extraRaw)
					}
					continue
				}
				appendEntry(&cfg, kind, entry)
			}
		}
	}
	return cfg, warns, nil
}

func claudeHandlerToEntry(event string, kind agenthooks.Kind, matcher string, groupIf json.RawMessage, handlerRaw json.RawMessage) (Entry, json.RawMessage, []Warning, bool) {
	var h claudeHandler
	if err := json.Unmarshal(handlerRaw, &h); err != nil {
		return Entry{}, nil, []Warning{warnf("%s: invalid handler JSON: %v", event, err)}, false
	}
	var warns []Warning
	e := Entry{
		Kind:          kind,
		NativeEvent:   event,
		Matcher:       matcher,
		TimeoutSec:    h.Timeout,
		ClaudeGroupIf: cloneRaw(groupIf),
		Raw:           cloneRaw(handlerRaw),
	}
	switch h.Type {
	case "command", "":
		e.Type = "command"
		e.Command = h.Command
		if len(h.Args) > 0 {
			e.Command = h.Command + " " + strings.Join(h.Args, " ")
			warns = append(warns, warnf("%s: exec-form args flattened into a shell command string", event))
		}
		return e, nil, warns, true
	case "prompt":
		e.Type = "prompt"
		e.Prompt = h.Prompt
		return e, nil, warns, true
	case "http":
		e.Type = "http"
		e.URL = h.URL
		return e, nil, warns, true
	default:
		groupRaw, err := json.Marshal(claudeGroup{Matcher: matcher, If: groupIf, Hooks: []json.RawMessage{handlerRaw}})
		if err != nil {
			warns = append(warns, warnf("%s: handler type %q could not be preserved: %v", event, h.Type, err))
			return Entry{}, nil, warns, false
		}
		warns = append(warns, warnf("%s: handler type %q is Claude-only; preserved in Extras", event, h.Type))
		return Entry{}, groupRaw, warns, false
	}
}

func emitClaude(cfg Config) ([]byte, []Warning, error) {
	f := claudeFile{Hooks: map[string][]claudeGroup{}}
	var warns []Warning
	for kind, entries := range cfg.Hooks {
		for _, e := range entries {
			event := eventNameForEmit(e, claudeKindForEvent, claudeEventForKind)
			if event == "" {
				warns = append(warns, warnf("kind %q has no Claude Code event name; dropped", kind))
				continue
			}
			handlerRaw, err := claudeHandlerRaw(e)
			if err != nil {
				warns = append(warns, warnf("%s: could not encode handler: %v", event, err))
				continue
			}
			f.Hooks[event] = append(f.Hooks[event], claudeGroup{
				Matcher: e.Matcher,
				If:      cloneRaw(e.ClaudeGroupIf),
				Hooks:   []json.RawMessage{handlerRaw},
			})
		}
	}
	for _, extra := range cfg.Extras {
		var g claudeGroup
		if err := json.Unmarshal(extra.Raw, &g); err != nil {
			warns = append(warns, warnf("%s: could not restore extra entry: %v", extra.Event, err))
			continue
		}
		f.Hooks[extra.Event] = append(f.Hooks[extra.Event], g)
	}
	out, err := json.MarshalIndent(f, "", "  ")
	return out, warns, err
}

func claudeHandlerRaw(e Entry) (json.RawMessage, error) {
	if len(e.Raw) > 0 {
		return cloneRaw(e.Raw), nil
	}
	h := claudeHandler{Type: e.Type, Command: e.Command, Prompt: e.Prompt, URL: e.URL, Timeout: e.TimeoutSec}
	return json.Marshal(h)
}
