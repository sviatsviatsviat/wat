package portconfig

import (
	"encoding/json"
	"fmt"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/agnostic"
	"github.com/sviatsviatsviat/wat/sdk/copilot"
)

func parseCopilot(data []byte) (Config, []Warning, error) {
	var f copilot.File
	if err := json.Unmarshal(data, &f); err != nil {
		return Config{}, nil, fmt.Errorf("portconfig: parse copilot hooks: %w", err)
	}
	cfg := Config{}
	var warns []Warning
	for event, handlers := range f.Hooks {
		kind, known := agnostic.CopilotKindForEvent(event)
		if !known {
			for _, handlerRaw := range handlers {
				appendExtra(&cfg, event, handlerRaw)
			}
			continue
		}
		for _, handlerRaw := range handlers {
			entry, extraRaw, w, ok := copilotHandlerToEntry(event, kind, handlerRaw)
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

func copilotHandlerToEntry(event string, kind agnostic.Kind, handlerRaw json.RawMessage) (Entry, json.RawMessage, []Warning, bool) {
	h, warns, ok := parseHandlerJSON[copilot.Handler](event, handlerRaw)
	if !ok {
		return Entry{}, cloneRaw(handlerRaw), warns, false
	}
	e := Entry{
		Kind:        kind,
		NativeEvent: event,
		Matcher:     h.Matcher,
		TimeoutSec:  h.TimeoutSeconds(),
		Raw:         cloneRaw(handlerRaw),
	}
	switch h.Type {
	case "command", "":
		e.Type = "command"
		e.Command = h.Command
		if e.Command == "" {
			e.Command = h.Bash
		}
		if h.PowerShell != "" && h.PowerShell != e.Command {
			warns = append(warns, warnf("%s: separate powershell command dropped (only bash/command ported)", event))
		}
		return e, nil, warns, true
	case "http":
		e.Type = "http"
		e.URL = h.URL
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

func emitCopilot(cfg Config) ([]byte, []Warning, error) {
	f := copilot.File{Version: 1, Hooks: map[string][]json.RawMessage{}}
	var warns []Warning
	for kind, entries := range cfg.Hooks {
		for _, e := range entries {
			event := eventNameForEmit(e, agnostic.CopilotKindForEventMap, agnostic.CopilotEventForKind)
			if event == "" {
				warns = append(warns, warnf("kind %q has no Copilot event name; dropped", kind))
				continue
			}
			if !copilotHandlerAllowed(e.Type, event) {
				if e.Type == "prompt" {
					warns = append(warns, warnf("%s: Copilot supports prompt hooks only on sessionStart; dropped", event))
				} else {
					warns = append(warns, warnf("%s: unsupported handler type %q; dropped", event, e.Type))
				}
				continue
			}
			handlerRaw, err := copilotHandlerRaw(e)
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

func copilotHandlerAllowed(handlerType, event string) bool {
	switch handlerType {
	case "command", "":
		return true
	case "http":
		return true
	case "prompt":
		return event == copilot.EventSessionStart
	default:
		return false
	}
}

func copilotHandlerRaw(e Entry) (json.RawMessage, error) {
	if len(e.Raw) == 0 {
		h := copilot.Handler{Matcher: e.Matcher, TimeoutSec: e.TimeoutSec}
		switch e.Type {
		case "command", "":
			h.Type = "command"
			h.Command = e.Command
		case "http":
			h.Type = "http"
			h.URL = e.URL
		case "prompt":
			h.Type = "prompt"
			h.Prompt = e.Prompt
		default:
			h.Type = e.Type
		}
		return hookkit.MarshalHandler(h)
	}
	var m map[string]any
	if err := json.Unmarshal(e.Raw, &m); err != nil {
		return nil, err
	}
	overlayCopilotHandlerFields(m, e)
	return json.Marshal(m)
}

func overlayCopilotHandlerFields(m map[string]any, e Entry) {
	if e.Matcher != "" {
		m["matcher"] = e.Matcher
	}
	if e.TimeoutSec != 0 {
		m["timeoutSec"] = e.TimeoutSec
	}
	switch e.Type {
	case "command", "":
		if e.Command != "" {
			m["type"] = "command"
			m["command"] = e.Command
		}
	case "http":
		if e.URL != "" {
			m["type"] = "http"
			m["url"] = e.URL
		}
	case "prompt":
		if e.Prompt != "" {
			m["type"] = "prompt"
			m["prompt"] = e.Prompt
		}
	default:
		if e.Type != "" {
			m["type"] = e.Type
		}
	}
}
