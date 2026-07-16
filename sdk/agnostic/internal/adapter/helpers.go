// Package adapter provides shared helpers for agent dialect codecs.
package adapter

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
)

// CloneRaw returns a defensive copy of raw JSON.
func CloneRaw(raw json.RawMessage) json.RawMessage {
	return hookkit.CloneRaw(raw)
}

// NewToolCall builds a ToolCall from native tool metadata and extracts a shell
// command when the tool is a shell execution.
func NewToolCall(native string, input json.RawMessage, id string) *model.ToolCall {
	name, mcp := model.NormalizeToolName(native)
	tc := &model.ToolCall{
		Name:   name,
		Native: native,
		ID:     id,
		MCP:    mcp,
	}
	if name != "" || native != "" || len(input) > 0 {
		tc.Input = model.NewToolInput(name, native, input)
	}
	if name == model.ToolBash {
		tc.Shell = hookkit.ExtractShellCommand(input)
	}
	return tc
}

// RawToText extracts a best-effort textual form of a tool_response value.
func RawToText(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	var s string
	if json.Unmarshal(raw, &s) == nil {
		return s
	}
	return string(raw)
}

// InvertKindEvent builds an event-name to kind map from a kind to event map.
func InvertKindEvent(m map[model.Kind]string) map[string]model.Kind {
	out := make(map[string]model.Kind, len(m))
	for kind, event := range m {
		out[event] = kind
	}
	return out
}

// MappedDecodeError rewrites an SDK decode error for agnostic while preserving Unwrap.
type MappedDecodeError struct {
	Msg string
	Err error
}

func (e *MappedDecodeError) Error() string { return e.Msg }

func (e *MappedDecodeError) Unwrap() error { return e.Err }

// MapDecodeErrorMessage rewrites SDK decode errors with an agent label.
func MapDecodeErrorMessage(err error, agent, sdk string) error {
	prefix := sdk + ": "
	msg := err.Error()
	if strings.HasPrefix(msg, prefix) {
		return &MappedDecodeError{
			Msg: fmt.Sprintf("%s: %s", agent, msg[len(prefix):]),
			Err: err,
		}
	}
	return fmt.Errorf("%s: %w", agent, err)
}

// RemapDecodeError replaces the error message while preserving Unwrap.
func RemapDecodeError(err error, msg string) error {
	return &MappedDecodeError{Msg: msg, Err: err}
}
