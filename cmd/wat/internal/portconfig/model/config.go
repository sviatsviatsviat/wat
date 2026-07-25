package model

import (
	"encoding/json"
	"fmt"

	"github.com/sviatsviatsviat/wat/cmd/wat/internal/hookconfig"
	"github.com/sviatsviatsviat/wat/internal/hookkit"
)

// Config is a normalized hook configuration independent of any single agent's
// config file shape.
type Config struct {
	// Hooks holds mappable hook registrations keyed by unified event kind.
	Hooks map[Kind][]Entry
	// Extras preserves unmappable native entries for round-trip safety.
	Extras []NativeEntry
}

// Entry is one normalized hook registration.
type Entry struct {
	// Kind is the unified event category.
	Kind Kind
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

// Warnf formats a warning message.
func Warnf(format string, args ...any) Warning {
	return Warning(fmt.Sprintf(format, args...))
}

// AppendEntry adds a hook entry to cfg.
func AppendEntry(cfg *Config, kind Kind, e Entry) {
	if cfg.Hooks == nil {
		cfg.Hooks = make(map[Kind][]Entry)
	}
	cfg.Hooks[kind] = append(cfg.Hooks[kind], e)
}

// AppendExtra adds an unmappable native entry to cfg.
func AppendExtra(cfg *Config, event string, raw json.RawMessage) {
	cfg.Extras = append(cfg.Extras, NativeEntry{Event: event, Raw: CloneRaw(raw)})
}

// CloneRaw returns a copy of raw JSON.
func CloneRaw(raw json.RawMessage) json.RawMessage {
	return hookkit.CloneBytes(raw)
}

// ParseHandlerJSON decodes handler JSON into T, returning warnings on failure.
func ParseHandlerJSON[T any](event string, handlerRaw json.RawMessage) (T, []Warning, bool) {
	h, err := hookconfig.ParseHandler[T](handlerRaw)
	if err != nil {
		var zero T
		return zero, []Warning{Warnf("%s: invalid handler JSON: %v", event, err)}, false
	}
	return h, nil, true
}
