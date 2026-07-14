package portconfig

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/agnostic"
)

// Config is a normalized hook configuration independent of any single agent's
// config file shape.
type Config struct {
	// Hooks holds mappable hook registrations keyed by unified event kind.
	Hooks map[agnostic.Kind][]Entry
	// Extras preserves unmappable native entries for round-trip safety.
	Extras []NativeEntry
}

// Entry is one normalized hook registration.
type Entry struct {
	// Kind is the unified event category.
	Kind agnostic.Kind
	// NativeEvent is the original config event key (e.g. "beforeShellExecution").
	NativeEvent string
	// Matcher is the native matcher string when the dialect supports one.
	Matcher string
	// Type is the handler kind: command, prompt, or http.
	Type string
	// Command is the shell command for command-type handlers.
	Command string
	// Prompt is the prompt text for prompt-type handlers.
	Prompt string
	// URL is the endpoint for http-type handlers.
	URL string
	// TimeoutSec is the handler timeout in seconds when set.
	TimeoutSec int
	// ClaudeGroupIf holds a Claude-only group-level permission rule when present.
	ClaudeGroupIf json.RawMessage
	// Raw is the untouched native handler JSON preserved for round-trip emit.
	Raw json.RawMessage
}

// NativeEntry preserves an unmappable native hook block for round-trip emit.
type NativeEntry struct {
	// Event is the native event name (e.g. "MessageDisplay").
	Event string
	// Raw is the untouched native JSON (a Claude group or flat handler object).
	Raw json.RawMessage
}

// Warning describes information lost or approximated during parse or emit.
type Warning string

func appendEntry(cfg *Config, kind agnostic.Kind, e Entry) {
	if cfg.Hooks == nil {
		cfg.Hooks = make(map[agnostic.Kind][]Entry)
	}
	cfg.Hooks[kind] = append(cfg.Hooks[kind], e)
}

func appendExtra(cfg *Config, event string, raw json.RawMessage) {
	cfg.Extras = append(cfg.Extras, NativeEntry{Event: event, Raw: cloneRaw(raw)})
}

func cloneRaw(raw json.RawMessage) json.RawMessage {
	return hookkit.CloneRaw(raw)
}

func parseHandlerJSON[T any](event string, handlerRaw json.RawMessage) (T, []Warning, bool) {
	h, err := hookkit.ParseHandler[T](handlerRaw)
	if err != nil {
		var zero T
		return zero, []Warning{warnf("%s: invalid handler JSON: %v", event, err)}, false
	}
	return h, nil, true
}
