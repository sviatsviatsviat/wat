package claudehook

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
)

type decodeFn func([]byte) (Event, error)

var decoders = map[string]decodeFn{
	EventSessionStart:        decodeAs[SessionStart],
	EventSetup:               decodeAs[Setup],
	EventSessionEnd:          decodeAs[SessionEnd],
	EventUserPromptSubmit:    decodeAs[UserPromptSubmit],
	EventUserPromptExpansion: decodeAs[UserPromptExpansion],
	EventPreToolUse:          decodeAs[PreToolUse],
	EventPostToolUse:         decodeAs[PostToolUse],
	EventPostToolUseFailure:  decodeAs[PostToolUseFailure],
	EventPostToolBatch:       decodeAs[PostToolBatch],
	EventPermissionRequest:   decodeAs[PermissionRequest],
	EventPermissionDenied:    decodeAs[PermissionDenied],
	EventSubagentStart:       decodeAs[SubagentStart],
	EventSubagentStop:        decodeAs[SubagentStop],
	EventTaskCreated:         decodeAs[TaskCreated],
	EventTaskCompleted:       decodeAs[TaskCompleted],
	EventStop:                decodeAs[Stop],
	EventStopFailure:         decodeAs[StopFailure],
	EventTeammateIdle:        decodeAs[TeammateIdle],
	EventNotification:        decodeAs[Notification],
	EventMessageDisplay:      decodeAs[MessageDisplay],
	EventInstructionsLoaded:  decodeAs[InstructionsLoaded],
	EventConfigChange:        decodeAs[ConfigChange],
	EventCwdChanged:          decodeAs[CwdChanged],
	EventFileChanged:         decodeAs[FileChanged],
	EventWorktreeCreate:      decodeAs[WorktreeCreate],
	EventWorktreeRemove:      decodeAs[WorktreeRemove],
	EventPreCompact:          decodeAs[PreCompact],
	EventPostCompact:         decodeAs[PostCompact],
	EventElicitation:         decodeAs[Elicitation],
	EventElicitationResult:   decodeAs[ElicitationResult],
}

func decodeAs[T Event](raw []byte) (Event, error) {
	var ev T
	if err := json.Unmarshal(raw, &ev); err != nil {
		return nil, fmt.Errorf("claudehook: decode %T: %w", ev, err)
	}
	return attachDecodedRaw(ev, raw), nil
}

func attachDecodedRaw[T Event](ev T, raw []byte) T {
	rv := reflect.ValueOf(&ev).Elem()
	ef := rv.FieldByName("Envelope")
	if !ef.IsValid() || !ef.CanSet() {
		return ev
	}
	env, ok := ef.Interface().(Envelope)
	if !ok {
		return ev
	}
	env.decodedRaw = cloneRaw(raw)
	ef.Set(reflect.ValueOf(env))
	return ev
}

// Decode parses a Claude Code hook stdin payload into a typed Event.
func Decode(raw []byte) (Event, error) {
	if len(raw) == 0 {
		return nil, fmt.Errorf("claudehook: empty payload")
	}
	var env Envelope
	if err := json.Unmarshal(raw, &env); err != nil {
		return nil, fmt.Errorf("claudehook: decode envelope: %w", err)
	}
	name := env.HookEventName
	if name == "" {
		return RawEvent{Envelope: env, Raw: cloneRaw(raw)}, nil
	}
	if fn, ok := decoders[name]; ok {
		return fn(raw)
	}
	return RawEvent{Envelope: env, Raw: cloneRaw(raw)}, nil
}

// ParseEvent reads and decodes a hook payload from r.
func ParseEvent(r io.Reader) (Event, error) {
	raw, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("claudehook: read payload: %w", err)
	}
	return Decode(raw)
}

func cloneRaw(raw []byte) json.RawMessage {
	if len(raw) == 0 {
		return nil
	}
	return append(json.RawMessage(nil), raw...)
}

// RawBytes returns the untouched JSON for an event when available.
func RawBytes(ev Event) json.RawMessage {
	switch e := ev.(type) {
	case RawEvent:
		return cloneRaw(e.Raw)
	default:
		if raw := decodedRawFrom(ev); len(raw) > 0 {
			return cloneRaw(raw)
		}
		b, err := json.Marshal(ev)
		if err != nil {
			return nil
		}
		return b
	}
}

func decodedRawFrom(ev Event) json.RawMessage {
	rv := reflect.ValueOf(ev)
	if rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return nil
		}
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return nil
	}
	ef := rv.FieldByName("Envelope")
	if !ef.IsValid() {
		return nil
	}
	env, ok := ef.Interface().(Envelope)
	if !ok {
		return nil
	}
	return env.decodedRaw
}

// CloneRaw copies raw JSON for independent mutation.
func CloneRaw(raw json.RawMessage) json.RawMessage {
	return cloneRaw(raw)
}
