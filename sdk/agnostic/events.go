package agnostic

import (
	"encoding/json"
	"fmt"
)

// Envelope carries shared metadata present on every normalized hook event.
type Envelope struct {
	// Agent is the dialect that emitted this hook event.
	Agent Dialect
	// Name is the native event name as received.
	Name string
	// Session holds session_id, sessionId, or conversation_id from the native payload.
	Session string
	// Cwd is the working directory from the native payload.
	Cwd string
	// TranscriptPath is the conversation transcript path when provided.
	TranscriptPath string
	// Raw is the untouched native JSON payload.
	Raw json.RawMessage
}

func envelopeFrom(ev *Event) Envelope {
	if ev == nil {
		return Envelope{}
	}
	return Envelope{
		Agent:          ev.Agent,
		Name:           ev.Name,
		Session:        ev.Session,
		Cwd:            ev.Cwd,
		TranscriptPath: ev.TranscriptPath,
		Raw:            ev.Raw,
	}
}

// PreToolEvent is the normalized view of a PreTool hook invocation.
type PreToolEvent struct {
	Envelope
	// Tool holds tool invocation details.
	Tool *ToolCall
}

// PreToolEventFrom maps a decoded Event to PreToolEvent.
func PreToolEventFrom(ev *Event) (PreToolEvent, error) {
	if ev == nil {
		return PreToolEvent{}, fmt.Errorf("agnostic: nil event")
	}
	if ev.Kind != KindPreTool {
		return PreToolEvent{}, fmt.Errorf("agnostic: expected PreTool kind, got %s", ev.Kind)
	}
	return PreToolEvent{Envelope: envelopeFrom(ev), Tool: ev.Tool}, nil
}

// PostToolEvent is the normalized view of a PostTool hook invocation.
type PostToolEvent struct {
	Envelope
	Tool   *ToolCall
	Result *ToolResult
}

// PostToolEventFrom maps a decoded Event to PostToolEvent.
func PostToolEventFrom(ev *Event) (PostToolEvent, error) {
	if ev == nil {
		return PostToolEvent{}, fmt.Errorf("agnostic: nil event")
	}
	if ev.Kind != KindPostTool {
		return PostToolEvent{}, fmt.Errorf("agnostic: expected PostTool kind, got %s", ev.Kind)
	}
	return PostToolEvent{Envelope: envelopeFrom(ev), Tool: ev.Tool, Result: ev.Result}, nil
}

// PostToolFailureEvent is the normalized view of a PostToolFailure hook invocation.
type PostToolFailureEvent struct {
	Envelope
	Tool   *ToolCall
	Result *ToolResult
}

// PostToolFailureEventFrom maps a decoded Event to PostToolFailureEvent.
func PostToolFailureEventFrom(ev *Event) (PostToolFailureEvent, error) {
	if ev == nil {
		return PostToolFailureEvent{}, fmt.Errorf("agnostic: nil event")
	}
	if ev.Kind != KindPostToolFailure {
		return PostToolFailureEvent{}, fmt.Errorf("agnostic: expected PostToolFailure kind, got %s", ev.Kind)
	}
	return PostToolFailureEvent{Envelope: envelopeFrom(ev), Tool: ev.Tool, Result: ev.Result}, nil
}

// StopEvent is the normalized view of Stop and SubagentStop hook invocations.
type StopEvent struct {
	Envelope
	Turn     *TurnEnd
	Subagent *Subagent
}

// StopEventFrom maps a decoded Event to StopEvent.
func StopEventFrom(ev *Event) (StopEvent, error) {
	if ev == nil {
		return StopEvent{}, fmt.Errorf("agnostic: nil event")
	}
	if ev.Kind != KindStop && ev.Kind != KindSubagentStop {
		return StopEvent{}, fmt.Errorf("agnostic: expected Stop or SubagentStop kind, got %s", ev.Kind)
	}
	return StopEvent{Envelope: envelopeFrom(ev), Turn: ev.Turn, Subagent: ev.Subagent}, nil
}

// SessionStartEvent is the normalized view of a SessionStart hook invocation.
type SessionStartEvent struct {
	Envelope
	Life *Lifecycle
}

// SessionStartEventFrom maps a decoded Event to SessionStartEvent.
func SessionStartEventFrom(ev *Event) (SessionStartEvent, error) {
	if ev == nil {
		return SessionStartEvent{}, fmt.Errorf("agnostic: nil event")
	}
	if ev.Kind != KindSessionStart {
		return SessionStartEvent{}, fmt.Errorf("agnostic: expected SessionStart kind, got %s", ev.Kind)
	}
	return SessionStartEvent{Envelope: envelopeFrom(ev), Life: ev.Life}, nil
}

// SessionEndEvent is the normalized view of a SessionEnd hook invocation.
type SessionEndEvent struct {
	Envelope
	Life *Lifecycle
}

// SessionEndEventFrom maps a decoded Event to SessionEndEvent.
func SessionEndEventFrom(ev *Event) (SessionEndEvent, error) {
	if ev == nil {
		return SessionEndEvent{}, fmt.Errorf("agnostic: nil event")
	}
	if ev.Kind != KindSessionEnd {
		return SessionEndEvent{}, fmt.Errorf("agnostic: expected SessionEnd kind, got %s", ev.Kind)
	}
	return SessionEndEvent{Envelope: envelopeFrom(ev), Life: ev.Life}, nil
}

// UserPromptEvent is the normalized view of a UserPrompt hook invocation.
type UserPromptEvent struct {
	Envelope
	Prompt string
}

// UserPromptEventFrom maps a decoded Event to UserPromptEvent.
func UserPromptEventFrom(ev *Event) (UserPromptEvent, error) {
	if ev == nil {
		return UserPromptEvent{}, fmt.Errorf("agnostic: nil event")
	}
	if ev.Kind != KindUserPrompt {
		return UserPromptEvent{}, fmt.Errorf("agnostic: expected UserPrompt kind, got %s", ev.Kind)
	}
	return UserPromptEvent{Envelope: envelopeFrom(ev), Prompt: ev.Prompt}, nil
}

// PreCompactEvent is the normalized view of a PreCompact hook invocation.
type PreCompactEvent struct {
	Envelope
	Compact *CompactInfo
}

// PreCompactEventFrom maps a decoded Event to PreCompactEvent.
func PreCompactEventFrom(ev *Event) (PreCompactEvent, error) {
	if ev == nil {
		return PreCompactEvent{}, fmt.Errorf("agnostic: nil event")
	}
	if ev.Kind != KindPreCompact {
		return PreCompactEvent{}, fmt.Errorf("agnostic: expected PreCompact kind, got %s", ev.Kind)
	}
	return PreCompactEvent{Envelope: envelopeFrom(ev), Compact: ev.Compact}, nil
}

// SubagentStartEvent is the normalized view of a SubagentStart hook invocation.
type SubagentStartEvent struct {
	Envelope
	Subagent *Subagent
}

// SubagentStartEventFrom maps a decoded Event to SubagentStartEvent.
func SubagentStartEventFrom(ev *Event) (SubagentStartEvent, error) {
	if ev == nil {
		return SubagentStartEvent{}, fmt.Errorf("agnostic: nil event")
	}
	if ev.Kind != KindSubagentStart {
		return SubagentStartEvent{}, fmt.Errorf("agnostic: expected SubagentStart kind, got %s", ev.Kind)
	}
	return SubagentStartEvent{Envelope: envelopeFrom(ev), Subagent: ev.Subagent}, nil
}

// AnyEvent is the catch-all normalized view for OnAny handlers.
type AnyEvent struct {
	Envelope
	Kind     Kind
	Prompt   string
	Tool     *ToolCall
	Result   *ToolResult
	Subagent *Subagent
	Turn     *TurnEnd
	Compact  *CompactInfo
	Note     *Note
	Life     *Lifecycle
}

// AnyEventFrom maps a decoded Event to AnyEvent.
func AnyEventFrom(ev *Event) AnyEvent {
	if ev == nil {
		return AnyEvent{}
	}
	return AnyEvent{
		Envelope: envelopeFrom(ev),
		Kind:     ev.Kind,
		Prompt:   ev.Prompt,
		Tool:     ev.Tool,
		Result:   ev.Result,
		Subagent: ev.Subagent,
		Turn:     ev.Turn,
		Compact:  ev.Compact,
		Note:     ev.Note,
		Life:     ev.Life,
	}
}
