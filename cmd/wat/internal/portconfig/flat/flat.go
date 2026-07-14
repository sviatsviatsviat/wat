package flat

import (
	"encoding/json"
	"fmt"

	"github.com/sviatsviatsviat/wat/cmd/wat/internal/portconfig/model"
	"github.com/sviatsviatsviat/wat/sdk/agnostic"
)

// hooksFile is the shared version-1 hooks.json envelope used by Copilot and Cursor.
type hooksFile struct {
	Hooks map[string][]json.RawMessage `json:"hooks"`
}

// HandlerToEntry converts one native handler JSON blob into a normalized entry.
// When ok is false, extraRaw holds unmappable JSON to preserve in Extras when non-nil.
type HandlerToEntry func(event string, kind agnostic.Kind, handlerRaw json.RawMessage) (model.Entry, json.RawMessage, []model.Warning, bool)

// ParseOptions configures flat hooks file parsing.
type ParseOptions struct {
	// Dialect names the config source in parse errors (e.g. "copilot hooks").
	Dialect string
	// KindForEvent maps a native event name to a unified kind.
	KindForEvent func(event string) (agnostic.Kind, bool)
	// SkipKind when true stores handlers for that kind in Extras instead of normalizing.
	SkipKind func(kind agnostic.Kind) bool
	// HandlerToEntry normalizes one handler object.
	HandlerToEntry HandlerToEntry
}

// Parse reads a flat hooks.json file into a normalized configuration.
func Parse(data []byte, opts ParseOptions) (model.Config, []model.Warning, error) {
	var f hooksFile
	if err := json.Unmarshal(data, &f); err != nil {
		return model.Config{}, nil, fmt.Errorf("portconfig: parse %s: %w", opts.Dialect, err)
	}
	cfg := model.Config{}
	var warns []model.Warning
	for event, handlers := range f.Hooks {
		kind, known := opts.KindForEvent(event)
		if !known {
			for _, handlerRaw := range handlers {
				model.AppendExtra(&cfg, event, handlerRaw)
			}
			continue
		}
		if opts.SkipKind != nil && opts.SkipKind(kind) {
			for _, handlerRaw := range handlers {
				model.AppendExtra(&cfg, event, handlerRaw)
			}
			continue
		}
		for _, handlerRaw := range handlers {
			entry, extraRaw, w, ok := opts.HandlerToEntry(event, kind, handlerRaw)
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
	return cfg, warns, nil
}

// AllowEntry reports whether a normalized entry may be emitted for event.
// When false, warns describe why the entry is dropped.
type AllowEntry func(e model.Entry, kind agnostic.Kind, event string) ([]model.Warning, bool)

// EncodeHandler renders one normalized entry as native handler JSON.
type EncodeHandler func(e model.Entry) (json.RawMessage, error)

// EmitOptions configures flat hooks file emission.
type EmitOptions struct {
	// Agent names the target agent in warning messages (e.g. "Copilot", "Cursor").
	Agent string
	// KindForEventMap maps native event names to unified kinds.
	KindForEventMap map[string]agnostic.Kind
	// EventForKind maps unified kinds to native event names.
	EventForKind map[agnostic.Kind]string
	// AllowEntry gates handler types unsupported on the target agent.
	AllowEntry AllowEntry
	// EncodeHandler renders a normalized entry as native handler JSON.
	EncodeHandler EncodeHandler
}

// Emit renders cfg as a version-1 flat hooks.json file.
func Emit(cfg model.Config, opts EmitOptions) ([]byte, []model.Warning, error) {
	f := hooksFile{Hooks: map[string][]json.RawMessage{}}
	var warns []model.Warning
	for kind, entries := range cfg.Hooks {
		for _, e := range entries {
			event := model.EventNameForEmit(e, opts.KindForEventMap, opts.EventForKind)
			if event == "" {
				warns = append(warns, model.Warnf("kind %q has no %s event name; dropped", kind, opts.Agent))
				continue
			}
			if opts.AllowEntry != nil {
				if entryWarns, ok := opts.AllowEntry(e, kind, event); !ok {
					warns = append(warns, entryWarns...)
					continue
				}
			}
			handlerRaw, err := opts.EncodeHandler(e)
			if err != nil {
				warns = append(warns, model.Warnf("%s: could not encode handler: %v", event, err))
				continue
			}
			f.Hooks[event] = append(f.Hooks[event], handlerRaw)
		}
	}
	for _, extra := range cfg.Extras {
		f.Hooks[extra.Event] = append(f.Hooks[extra.Event], model.CloneRaw(extra.Raw))
	}
	out, err := json.MarshalIndent(struct {
		Version int                          `json:"version"`
		Hooks   map[string][]json.RawMessage `json:"hooks"`
	}{Version: 1, Hooks: f.Hooks}, "", "  ")
	return out, warns, err
}
