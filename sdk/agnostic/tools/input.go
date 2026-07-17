package tools

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
)

// Input is the tool input payload on a normalized ToolCall.
type Input struct {
	name   string
	native string
	raw    json.RawMessage
}

// NewInput returns an Input bound to canonical and native tool names.
// It panics if both names and raw are empty.
func NewInput(canonical, native string, raw json.RawMessage) Input {
	if canonical == "" && native == "" && len(raw) == 0 {
		panic("tools: NewInput called with empty tool names and empty payload")
	}
	return Input{
		name:   canonical,
		native: native,
		raw:    hookkit.CloneBytes(raw),
	}
}

// Name returns the canonical tool name (bash, write, read, …).
func (in Input) Name() string { return in.name }

// Native returns the original agent tool name.
func (in Input) Native() string { return in.native }

// Raw returns a copy of the native tool input JSON.
func (in Input) Raw() json.RawMessage { return hookkit.CloneBytes(in.raw) }

// HasRaw reports whether this input carries a JSON payload.
func (in Input) HasRaw() bool { return len(in.raw) > 0 }

func as[T any](in Input, expected string) (T, bool) {
	var v T
	if in.name != expected {
		return v, false
	}
	if len(in.raw) == 0 {
		return v, true
	}
	if json.Unmarshal(in.raw, &v) == nil {
		return v, true
	}
	var zero T
	return zero, false
}
