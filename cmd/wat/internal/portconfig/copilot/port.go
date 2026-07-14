package copilot

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/cmd/wat/internal/portconfig/flat"
	"github.com/sviatsviatsviat/wat/cmd/wat/internal/portconfig/model"
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/agnostic"
	"github.com/sviatsviatsviat/wat/sdk/copilot"
)

// Parse reads GitHub Copilot hooks JSON into a normalized configuration.
func Parse(data []byte) (model.Config, []model.Warning, error) {
	return flat.Parse(data, flat.ParseOptions{
		Dialect:        "copilot hooks",
		KindForEvent:   agnostic.CopilotKindForEvent,
		HandlerToEntry: copilotHandlerToEntry,
	})
}

func copilotHandlerToEntry(event string, kind agnostic.Kind, handlerRaw json.RawMessage) (model.Entry, json.RawMessage, []model.Warning, bool) {
	h, warns, ok := model.ParseHandlerJSON[copilot.Handler](event, handlerRaw)
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
	case "command", "":
		e.Type = "command"
		e.Command = h.Command
		if e.Command == "" {
			e.Command = h.Bash
		}
		if h.PowerShell != "" && h.PowerShell != e.Command {
			warns = append(warns, model.Warnf("%s: separate powershell command dropped (only bash/command ported)", event))
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
		warns = append(warns, model.Warnf("%s: unknown handler type %q preserved in Extras", event, h.Type))
		return model.Entry{}, model.CloneRaw(handlerRaw), warns, false
	}
}

// Emit renders cfg as GitHub Copilot hooks JSON.
func Emit(cfg model.Config) ([]byte, []model.Warning, error) {
	return flat.Emit(cfg, flat.EmitOptions{
		Agent:           "Copilot",
		KindForEventMap: agnostic.CopilotKindForEventMap,
		EventForKind:    agnostic.CopilotEventForKind,
		AllowEntry:      copilotAllowEntry,
		EncodeHandler:   copilotHandlerRaw,
	})
}

func copilotAllowEntry(e model.Entry, kind agnostic.Kind, event string) ([]model.Warning, bool) {
	if copilotHandlerAllowed(e.Type, event) {
		return nil, true
	}
	if e.Type == "prompt" {
		return []model.Warning{model.Warnf("%s: Copilot supports prompt hooks only on sessionStart; dropped", event)}, false
	}
	return []model.Warning{model.Warnf("%s: unsupported handler type %q; dropped", event, e.Type)}, false
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

func copilotHandlerRaw(e model.Entry) (json.RawMessage, error) {
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

func overlayCopilotHandlerFields(m map[string]any, e model.Entry) {
	if e.Matcher != "" {
		m["matcher"] = e.Matcher
	} else {
		delete(m, "matcher")
	}
	if e.TimeoutSec != 0 {
		m["timeoutSec"] = e.TimeoutSec
	} else {
		delete(m, "timeoutSec")
		delete(m, "timeout")
	}

	switch e.Type {
	case "command", "":
		m["type"] = "command"
		delete(m, "url")
		delete(m, "prompt")
		if e.Command != "" {
			m["command"] = e.Command
			delete(m, "bash")
			delete(m, "powershell")
		} else {
			delete(m, "command")
		}
	case "http":
		m["type"] = "http"
		delete(m, "command")
		delete(m, "bash")
		delete(m, "powershell")
		delete(m, "prompt")
		if e.URL != "" {
			m["url"] = e.URL
		} else {
			delete(m, "url")
		}
	case "prompt":
		m["type"] = "prompt"
		delete(m, "command")
		delete(m, "bash")
		delete(m, "powershell")
		delete(m, "url")
		if e.Prompt != "" {
			m["prompt"] = e.Prompt
		} else {
			delete(m, "prompt")
		}
	default:
		if e.Type != "" {
			m["type"] = e.Type
		}
	}
}
