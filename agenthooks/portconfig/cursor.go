package portconfig

import (
	"encoding/json"
	"fmt"

	"github.com/sviatsviatsviat/wat/agenthooks"
)

type cursorFile struct {
	Version int                          `json:"version"`
	Hooks   map[string][]json.RawMessage `json:"hooks"`
}

type cursorHandler struct {
	Command string `json:"command,omitempty"`
	Type    string `json:"type,omitempty"`
	Prompt  string `json:"prompt,omitempty"`
	Matcher string `json:"matcher,omitempty"`
	Timeout int    `json:"timeout,omitempty"`
}

func parseCursor(data []byte) (Config, []Warning, error) {
	var f cursorFile
	if err := json.Unmarshal(data, &f); err != nil {
		return Config{}, nil, fmt.Errorf("portconfig: parse cursor hooks: %w", err)
	}
	cfg := Config{}
	var warns []Warning
	for event, handlers := range f.Hooks {
		kind, known := cursorKindForEvent[event]
		if !known {
			for _, handlerRaw := range handlers {
				appendExtra(&cfg, event, handlerRaw)
			}
			continue
		}
		if kind == agenthooks.KindOther {
			for _, handlerRaw := range handlers {
				appendExtra(&cfg, event, handlerRaw)
			}
			continue
		}
		for _, handlerRaw := range handlers {
			entry, extraRaw, w, ok := cursorHandlerToEntry(event, kind, handlerRaw)
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
	return cfg, warns, nil
}

func cursorHandlerToEntry(event string, kind agenthooks.Kind, handlerRaw json.RawMessage) (Entry, json.RawMessage, []Warning, bool) {
	var h cursorHandler
	if err := json.Unmarshal(handlerRaw, &h); err != nil {
		return Entry{}, nil, []Warning{warnf("%s: invalid handler JSON: %v", event, err)}, false
	}
	var warns []Warning
	e := Entry{
		Kind:        kind,
		NativeEvent: event,
		Matcher:     h.Matcher,
		TimeoutSec:  h.Timeout,
		Raw:         cloneRaw(handlerRaw),
	}
	switch h.Type {
	case "command", "":
		e.Type = "command"
		e.Command = h.Command
		return e, nil, warns, true
	case "prompt":
		e.Type = "prompt"
		e.Prompt = h.Prompt
		return e, nil, warns, true
	default:
		warns = append(warns, warnf("%s: unknown handler type %q preserved in Extras", event, h.Type))
		return Entry{}, cloneRaw(handlerRaw), warns, false
	}
}

func emitCursor(cfg Config) ([]byte, []Warning, error) {
	f := cursorFile{Version: 1, Hooks: map[string][]json.RawMessage{}}
	var warns []Warning
	for kind, entries := range cfg.Hooks {
		for _, e := range entries {
			event := eventNameForEmit(e, cursorKindForEvent, cursorEventForKind)
			if event == "" {
				warns = append(warns, warnf("kind %q has no Cursor event name; dropped", kind))
				continue
			}
			if e.Type == "http" {
				warns = append(warns, warnf("%s: Cursor has no http hooks; dropped", e.NativeEvent))
				continue
			}
			handlerRaw, err := cursorHandlerRaw(e)
			if err != nil {
				warns = append(warns, warnf("%s: could not encode handler: %v", event, err))
				continue
			}
			f.Hooks[event] = append(f.Hooks[event], handlerRaw)
		}
	}
	for _, extra := range cfg.Extras {
		f.Hooks[extra.Event] = append(f.Hooks[extra.Event], cloneRaw(extra.Raw))
	}
	out, err := json.MarshalIndent(f, "", "  ")
	return out, warns, err
}

func cursorHandlerRaw(e Entry) (json.RawMessage, error) {
	if len(e.Raw) > 0 {
		return cloneRaw(e.Raw), nil
	}
	h := cursorHandler{
		Command: e.Command,
		Type:    nonDefault(e.Type, "command"),
		Prompt:  e.Prompt,
		Matcher: e.Matcher,
		Timeout: e.TimeoutSec,
	}
	return json.Marshal(h)
}

func nonDefault(v, def string) string {
	if v == def {
		return ""
	}
	return v
}
