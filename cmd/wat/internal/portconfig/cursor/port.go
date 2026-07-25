package cursor

import (
	"encoding/json"

	hostcursor "github.com/sviatsviatsviat/wat/cmd/wat/internal/hostconfig/cursor"
	"github.com/sviatsviatsviat/wat/cmd/wat/internal/portconfig/flat"
	"github.com/sviatsviatsviat/wat/cmd/wat/internal/portconfig/model"
	"github.com/sviatsviatsviat/wat/cmd/wat/internal/hookconfig"
)

// Parse reads Cursor hooks JSON into a normalized configuration.
func Parse(data []byte) (model.Config, []model.Warning, error) {
	return flat.Parse(data, flat.ParseOptions{
		Dialect:        "cursor hooks",
		KindForEvent:   kindForEvent,
		SkipKind:       func(kind model.Kind) bool { return kind == model.KindOther },
		HandlerToEntry: cursorHandlerToEntry,
	})
}

func cursorHandlerToEntry(event string, kind model.Kind, handlerRaw json.RawMessage) (model.Entry, json.RawMessage, []model.Warning, bool) {
	h, warns, ok := model.ParseHandlerJSON[hostcursor.Handler](event, handlerRaw)
	if !ok {
		return model.Entry{}, model.CloneRaw(handlerRaw), warns, false
	}
	e := model.Entry{
		Kind:        kind,
		NativeEvent: event,
		Matcher:     h.Matcher,
		TimeoutSec:  h.TimeoutSeconds(),
		Raw:         model.CloneRaw(handlerRaw),
	}
	switch h.Type {
	case hostcursor.HandlerTypeCommand, "":
		e.Type = hostcursor.HandlerTypeCommand
		e.Command = h.Command
		return e, nil, warns, true
	case hostcursor.HandlerTypePrompt:
		e.Type = hostcursor.HandlerTypePrompt
		e.Prompt = h.Prompt
		return e, nil, warns, true
	default:
		warns = append(warns, model.Warnf("%s: unknown handler type %q preserved in Extras", event, h.Type))
		return model.Entry{}, model.CloneRaw(handlerRaw), warns, false
	}
}

// Emit renders cfg as Cursor hooks JSON.
func Emit(cfg model.Config) ([]byte, []model.Warning, error) {
	return flat.Emit(cfg, flat.EmitOptions{
		Agent:           "Cursor",
		KindForEventMap: kindForEventMap,
		EventForKind:    EventForKind,
		AllowEntry:      cursorAllowEntry,
		EncodeHandler:   cursorHandlerRaw,
	})
}

func cursorAllowEntry(e model.Entry, _ model.Kind, event string) ([]model.Warning, bool) {
	switch e.Type {
	case hostcursor.HandlerTypeCommand, "", hostcursor.HandlerTypePrompt:
		return nil, true
	case "http":
		return []model.Warning{model.Warnf("%s: Cursor has no http hooks; dropped", e.NativeEvent)}, false
	default:
		if event == "" {
			event = e.NativeEvent
		}
		return []model.Warning{model.Warnf("%s: unsupported handler type %q; dropped", event, e.Type)}, false
	}
}

func cursorHandlerRaw(e model.Entry) (json.RawMessage, error) {
	if len(e.Raw) == 0 {
		h := hostcursor.Handler{
			Command: e.Command,
			Prompt:  e.Prompt,
			Matcher: e.Matcher,
			Timeout: e.TimeoutSec,
		}
		if e.Type != "" && e.Type != hostcursor.HandlerTypeCommand {
			h.Type = e.Type
		}
		return hookconfig.MarshalHandler(h)
	}
	var m map[string]any
	if err := json.Unmarshal(e.Raw, &m); err != nil {
		return nil, err
	}
	overlayCursorHandlerFields(m, e)
	return json.Marshal(m)
}

func overlayCursorHandlerFields(m map[string]any, e model.Entry) {
	if e.Matcher != "" {
		m["matcher"] = e.Matcher
	} else {
		delete(m, "matcher")
	}
	if e.TimeoutSec != 0 {
		m["timeout"] = e.TimeoutSec
	} else {
		delete(m, "timeout")
	}

	handlerType := e.Type
	if handlerType == "" {
		handlerType = hostcursor.HandlerTypeCommand
	}
	delete(m, "command")
	delete(m, "prompt")
	delete(m, "type")

	switch handlerType {
	case hostcursor.HandlerTypeCommand:
		if e.Command != "" {
			m["command"] = e.Command
		}
	case hostcursor.HandlerTypePrompt:
		m["type"] = hostcursor.HandlerTypePrompt
		if e.Prompt != "" {
			m["prompt"] = e.Prompt
		}
	default:
		if handlerType != "" {
			m["type"] = handlerType
		}
	}
}
