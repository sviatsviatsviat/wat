package hookkit

import (
	"encoding/json"
	"strings"
)

// Input is a tool name bound to raw tool-input JSON, shared by agent SDKs.
type Input struct {
	tool string
	raw  json.RawMessage
}

// NewInput returns an Input bound to the native tool name and raw JSON.
// It panics if both tool and raw are empty.
func NewInput(tool string, raw json.RawMessage) Input {
	if tool == "" && len(raw) == 0 {
		panic("hookkit: NewInput called with empty tool name and empty payload")
	}
	return Input{tool: tool, raw: CloneBytes(raw)}
}

// Name returns the native tool name bound to this input.
func (in Input) Name() string { return in.tool }

// Raw returns a copy of the native tool input JSON.
func (in Input) Raw() json.RawMessage { return CloneBytes(in.raw) }

// HasRaw reports whether this input carries a JSON payload.
func (in Input) HasRaw() bool { return len(in.raw) > 0 }

// As decodes the input as T when the bound tool name equals expected.
func As[T any](in Input, expected string) (T, bool) {
	var v T
	if in.tool != expected {
		return v, false
	}
	return decodeAs[T](in)
}

// AsFold decodes the input as T when the bound tool name matches any of expected
// (case-insensitive). Empty raw yields the zero value with ok true.
func AsFold[T any](in Input, expected ...string) (T, bool) {
	var v T
	matched := false
	for _, name := range expected {
		if strings.EqualFold(in.tool, name) {
			matched = true
			break
		}
	}
	if !matched {
		return v, false
	}
	return decodeAs[T](in)
}

func decodeAs[T any](in Input) (T, bool) {
	var v T
	if len(in.raw) == 0 {
		return v, true
	}
	if json.Unmarshal(in.raw, &v) == nil {
		return v, true
	}
	var zero T
	return zero, false
}
