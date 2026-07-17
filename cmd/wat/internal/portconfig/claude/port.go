package claude

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sviatsviatsviat/wat/cmd/wat/internal/portconfig/model"
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/agnostic"
	agclaude "github.com/sviatsviatsviat/wat/sdk/agnostic/claude"
	"github.com/sviatsviatsviat/wat/sdk/claude"
)

// Parse reads Claude Code settings JSON into a normalized configuration.
func Parse(data []byte) (model.Config, []model.Warning, error) {
	var f claude.Settings
	if err := json.Unmarshal(data, &f); err != nil {
		return model.Config{}, nil, fmt.Errorf("portconfig: parse claude settings: %w", err)
	}
	cfg := model.Config{}
	var warns []model.Warning
	for event, groups := range f.Hooks {
		kind, known := agclaude.KindForEvent(event)
		if !known {
			for _, g := range groups {
				raw, err := json.Marshal(g)
				if err != nil {
					return model.Config{}, warns, fmt.Errorf("portconfig: marshal claude extra: %w", err)
				}
				model.AppendExtra(&cfg, event, raw)
			}
			continue
		}
		for _, g := range groups {
			for _, handlerRaw := range g.Hooks {
				entry, extraRaw, w, ok := claudeHandlerToEntry(event, kind, g.Matcher, g.If, handlerRaw)
				warns = append(warns, w...)
				if !ok {
					if extraRaw != nil {
						model.AppendExtra(&cfg, event, extraRaw)
					}
					continue
				}
				model.AppendEntry(&cfg, kind, entry)
			}
		}
	}
	return cfg, warns, nil
}

func claudeHandlerToEntry(event string, kind agnostic.Kind, matcher string, groupIf json.RawMessage, handlerRaw json.RawMessage) (model.Entry, json.RawMessage, []model.Warning, bool) {
	h, warns, ok := model.ParseHandlerJSON[claude.Handler](event, handlerRaw)
	if !ok {
		return model.Entry{}, model.CloneRaw(handlerRaw), warns, false
	}
	e := model.Entry{
		Kind:          kind,
		NativeEvent:   event,
		Matcher:       matcher,
		TimeoutSec:    h.Timeout,
		ClaudeGroupIf: model.CloneRaw(groupIf),
		Raw:           model.CloneRaw(handlerRaw),
	}
	switch h.Type {
	case "command", "":
		e.Type = "command"
		e.Command = h.Command
		if len(h.Args) > 0 {
			e.Command = h.Command + " " + strings.Join(h.Args, " ")
			warns = append(warns, model.Warnf("%s: exec-form args flattened into a shell command string", event))
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
		groupRaw, err := json.Marshal(claude.MatcherGroup{
			Matcher: matcher,
			If:      groupIf,
			Hooks:   []json.RawMessage{handlerRaw},
		})
		if err != nil {
			warns = append(warns, model.Warnf("%s: handler type %q could not be preserved: %v", event, h.Type, err))
			return model.Entry{}, nil, warns, false
		}
		warns = append(warns, model.Warnf("%s: handler type %q is Claude-only; preserved in Extras", event, h.Type))
		return model.Entry{}, groupRaw, warns, false
	}
}

// Emit renders cfg as Claude Code settings JSON.
func Emit(cfg model.Config) ([]byte, []model.Warning, error) {
	f := claude.Settings{Hooks: map[string][]claude.MatcherGroup{}}
	var warns []model.Warning
	for kind, entries := range cfg.Hooks {
		for _, e := range entries {
			event := model.EventNameForEmit(e, agclaude.KindForEventMap, agclaude.EventForKind)
			if event == "" {
				warns = append(warns, model.Warnf("kind %q has no Claude Code event name; dropped", kind))
				continue
			}
			handlerRaw, err := claudeHandlerRaw(e)
			if err != nil {
				warns = append(warns, model.Warnf("%s: could not encode handler: %v", event, err))
				continue
			}
			f.Hooks[event] = append(f.Hooks[event], claude.MatcherGroup{
				Matcher: e.Matcher,
				If:      model.CloneRaw(e.ClaudeGroupIf),
				Hooks:   []json.RawMessage{handlerRaw},
			})
		}
	}
	for _, extra := range cfg.Extras {
		var g claude.MatcherGroup
		if err := json.Unmarshal(extra.Raw, &g); err != nil {
			warns = append(warns, model.Warnf("%s: could not restore extra entry: %v", extra.Event, err))
			continue
		}
		f.Hooks[extra.Event] = append(f.Hooks[extra.Event], g)
	}
	out, err := json.MarshalIndent(f, "", "  ")
	return out, warns, err
}

func claudeHandlerRaw(e model.Entry) (json.RawMessage, error) {
	if len(e.Raw) == 0 {
		return hookkit.MarshalHandler(claude.Handler{
			Type:    e.Type,
			Command: e.Command,
			Prompt:  e.Prompt,
			URL:     e.URL,
			Timeout: e.TimeoutSec,
		})
	}
	var m map[string]any
	if err := json.Unmarshal(e.Raw, &m); err != nil {
		return nil, err
	}
	overlayClaudeHandlerFields(m, e)
	return json.Marshal(m)
}

func overlayClaudeHandlerFields(m map[string]any, e model.Entry) {
	if e.Type != "" {
		m["type"] = e.Type
	}
	if e.Command != "" {
		m["command"] = e.Command
	}
	if e.Prompt != "" {
		m["prompt"] = e.Prompt
	}
	if e.URL != "" {
		m["url"] = e.URL
	}
	if e.TimeoutSec != 0 {
		m["timeout"] = e.TimeoutSec
	}
}
