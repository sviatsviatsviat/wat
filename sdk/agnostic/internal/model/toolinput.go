package model

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
)

// ToolInput is the tool input payload on a normalized ToolCall.
type ToolInput struct {
	name   string
	native string
	raw    json.RawMessage
}

// NewToolInput returns a ToolInput bound to canonical and native tool names.
// It panics if both names and raw are empty.
func NewToolInput(canonical, native string, raw json.RawMessage) ToolInput {
	if canonical == "" && native == "" && len(raw) == 0 {
		panic("model: NewToolInput called with empty tool names and empty payload")
	}
	return ToolInput{
		name:   canonical,
		native: native,
		raw:    hookkit.CloneBytes(raw),
	}
}

// Name returns the canonical tool name (bash, write, read, …).
func (in ToolInput) Name() string { return in.name }

// Native returns the original agent tool name.
func (in ToolInput) Native() string { return in.native }

// Raw returns a copy of the native tool input JSON.
func (in ToolInput) Raw() json.RawMessage { return hookkit.CloneBytes(in.raw) }

// HasRaw reports whether this input carries a JSON payload.
func (in ToolInput) HasRaw() bool { return len(in.raw) > 0 }

func as[T any](in ToolInput, expected string) (T, bool) {
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
