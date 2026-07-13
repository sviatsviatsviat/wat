package copilothook

import "encoding/json"

// envelopeAccessor provides compile-time-checked envelope metadata access.
type envelopeAccessor interface {
	envelopePtr() *Envelope
}

func (e *Envelope) setEnvelopeMeta(received, canonical string, raw json.RawMessage) {
	e.receivedName = received
	e.canonical = canonical
	e.decodedRaw = cloneRaw(raw)
}

func (e *Envelope) decodedRawBytes() json.RawMessage {
	return cloneRaw(e.decodedRaw)
}

func envelopeAccessorForEvent(ev Event) envelopeAccessor {
	switch e := ev.(type) {
	case SessionStart:
		return &e
	case SessionEnd:
		return &e
	case UserPromptSubmitted:
		return &e
	case PreToolUse:
		return &e
	case PostToolUse:
		return &e
	case PostToolUseFailure:
		return &e
	case PermissionRequest:
		return &e
	case SubagentStart:
		return &e
	case SubagentStop:
		return &e
	case AgentStop:
		return &e
	case PreCompact:
		return &e
	case Notification:
		return &e
	case ErrorOccurred:
		return &e
	case RawEvent:
		return &e
	default:
		panic("copilothook: event missing envelopeAccessor")
	}
}

func envelopeAccessorForValue(v any) envelopeAccessor {
	switch e := v.(type) {
	case *SessionStart:
		return e
	case *SessionEnd:
		return e
	case *UserPromptSubmitted:
		return e
	case *PreToolUse:
		return e
	case *PostToolUse:
		return e
	case *PostToolUseFailure:
		return e
	case *PermissionRequest:
		return e
	case *SubagentStart:
		return e
	case *SubagentStop:
		return e
	case *AgentStop:
		return e
	case *PreCompact:
		return e
	case *Notification:
		return e
	case *ErrorOccurred:
		return e
	case *RawEvent:
		return e
	default:
		panic("copilothook: value missing envelopeAccessor")
	}
}

var (
	_ envelopeAccessor = (*SessionStart)(nil)
	_ envelopeAccessor = (*SessionEnd)(nil)
	_ envelopeAccessor = (*UserPromptSubmitted)(nil)
	_ envelopeAccessor = (*PreToolUse)(nil)
	_ envelopeAccessor = (*PostToolUse)(nil)
	_ envelopeAccessor = (*PostToolUseFailure)(nil)
	_ envelopeAccessor = (*PermissionRequest)(nil)
	_ envelopeAccessor = (*SubagentStart)(nil)
	_ envelopeAccessor = (*SubagentStop)(nil)
	_ envelopeAccessor = (*AgentStop)(nil)
	_ envelopeAccessor = (*PreCompact)(nil)
	_ envelopeAccessor = (*Notification)(nil)
	_ envelopeAccessor = (*ErrorOccurred)(nil)
	_ envelopeAccessor = (*RawEvent)(nil)
)

func (e *SessionStart) envelopePtr() *Envelope        { return &e.Envelope }
func (e *SessionEnd) envelopePtr() *Envelope          { return &e.Envelope }
func (e *UserPromptSubmitted) envelopePtr() *Envelope { return &e.Envelope }
func (e *PreToolUse) envelopePtr() *Envelope          { return &e.Envelope }
func (e *PostToolUse) envelopePtr() *Envelope         { return &e.Envelope }
func (e *PostToolUseFailure) envelopePtr() *Envelope  { return &e.Envelope }
func (e *PermissionRequest) envelopePtr() *Envelope   { return &e.Envelope }
func (e *SubagentStart) envelopePtr() *Envelope       { return &e.Envelope }
func (e *SubagentStop) envelopePtr() *Envelope        { return &e.Envelope }
func (e *AgentStop) envelopePtr() *Envelope           { return &e.Envelope }
func (e *PreCompact) envelopePtr() *Envelope          { return &e.Envelope }
func (e *Notification) envelopePtr() *Envelope        { return &e.Envelope }
func (e *ErrorOccurred) envelopePtr() *Envelope       { return &e.Envelope }
func (e *RawEvent) envelopePtr() *Envelope            { return &e.Envelope }
