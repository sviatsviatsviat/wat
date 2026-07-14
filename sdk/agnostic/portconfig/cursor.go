package portconfig

import (
	"encoding/json"
	"fmt"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/agnostic"
	"github.com/sviatsviatsviat/wat/sdk/cursor"
)

func parseCursor(data []byte) (Config, []Warning, error) {
	var f cursor.File
	if err := json.Unmarshal(data, &f); err != nil {
		return Config{}, nil, fmt.Errorf("portconfig: parse cursor hooks: %w", err)
	}
	cfg := Config{}
	var warns []Warning
	for event, handlers := range f.Hooks {
		kind, known := agnostic.CursorKindForEvent(event)
		if !known {
			for _, handlerRaw := range handlers {
				appendExtra(&cfg, event, handlerRaw)
			}
			continue
		}
		if kind == agnostic.KindOther {
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

func cursorHandlerToEntry(event string, kind agnostic.Kind, handlerRaw json.RawMessage) (Entry, json.RawMessage, []Warning, bool) {
	h, warns, ok := parseHandlerJSON[cursor.Handler](event, handlerRaw)
	if !ok {
		return Entry{}, nil, warns, false
	}
	e := Entry{
		Kind:        kind,
		NativeEvent: event,
		Matcher:     h.Matcher,
		TimeoutSec:  h.TimeoutSeconds(),
		Raw:         cloneRaw(handlerRaw),
	}
	switch h.Type {
	case cursor.HandlerTypeCommand, "":
		e.Type = cursor.HandlerTypeCommand
		e.Command = h.Command
		return e, nil, warns, true
	case cursor.HandlerTypePrompt:
		e.Type = cursor.HandlerTypePrompt
		e.Prompt = h.Prompt
		return e, nil, warns, true
	default:
		warns = append(warns, warnf("%s: unknown handler type %q preserved in Extras", event, h.Type))
		return Entry{}, cloneRaw(handlerRaw), warns, false
	}
}

func emitCursor(cfg Config) ([]byte, []Warning, error) {
	f := cursor.File{Version: 1, Hooks: map[string][]json.RawMessage{}}
	var warns []Warning
	for kind, entries := range cfg.Hooks {
		for _, e := range entries {
			event := eventNameForEmit(e, agnostic.CursorKindForEventMap, agnostic.CursorEventForKind)
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
	h := cursor.Handler{
		Command: e.Command,
		Prompt:  e.Prompt,
		Matcher: e.Matcher,
		Timeout: e.TimeoutSec,
	}
	if e.Type != "" && e.Type != cursor.HandlerTypeCommand {
		h.Type = e.Type
	}
	return hookkit.MarshalHandler(h)
}
