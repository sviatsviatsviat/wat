package claude

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
)

// envelopeAccessor provides compile-time-checked envelope metadata access.
type envelopeAccessor interface {
	envelopePtr() *Envelope
}

func (e *Envelope) setDecodedRaw(raw json.RawMessage) {
	e.decodedRaw = hookkit.CloneRaw(raw)
}

// DecodedRaw returns the untouched JSON stored on the envelope.
func (e *Envelope) DecodedRaw() json.RawMessage {
	return hookkit.CloneRaw(e.decodedRaw)
}

func envelopeAccessorForEvent(ev Event) envelopeAccessor {
	switch e := ev.(type) {
	case SessionStart:
		return &e
	case Setup:
		return &e
	case SessionEnd:
		return &e
	case UserPromptSubmit:
		return &e
	case UserPromptExpansion:
		return &e
	case PreToolUse:
		return &e
	case PostToolUse:
		return &e
	case PostToolUseFailure:
		return &e
	case PostToolBatch:
		return &e
	case PermissionRequest:
		return &e
	case PermissionDenied:
		return &e
	case SubagentStart:
		return &e
	case SubagentStop:
		return &e
	case TaskCreated:
		return &e
	case TaskCompleted:
		return &e
	case Stop:
		return &e
	case StopFailure:
		return &e
	case TeammateIdle:
		return &e
	case Notification:
		return &e
	case MessageDisplay:
		return &e
	case InstructionsLoaded:
		return &e
	case ConfigChange:
		return &e
	case CwdChanged:
		return &e
	case FileChanged:
		return &e
	case WorktreeCreate:
		return &e
	case WorktreeRemove:
		return &e
	case PreCompact:
		return &e
	case PostCompact:
		return &e
	case Elicitation:
		return &e
	case ElicitationResult:
		return &e
	case RawEvent:
		return &e
	default:
		panic("claude: event missing envelopeAccessor")
	}
}

func envelopeAccessorForValue(v any) envelopeAccessor {
	switch e := v.(type) {
	case *SessionStart:
		return e
	case *Setup:
		return e
	case *SessionEnd:
		return e
	case *UserPromptSubmit:
		return e
	case *UserPromptExpansion:
		return e
	case *PreToolUse:
		return e
	case *PostToolUse:
		return e
	case *PostToolUseFailure:
		return e
	case *PostToolBatch:
		return e
	case *PermissionRequest:
		return e
	case *PermissionDenied:
		return e
	case *SubagentStart:
		return e
	case *SubagentStop:
		return e
	case *TaskCreated:
		return e
	case *TaskCompleted:
		return e
	case *Stop:
		return e
	case *StopFailure:
		return e
	case *TeammateIdle:
		return e
	case *Notification:
		return e
	case *MessageDisplay:
		return e
	case *InstructionsLoaded:
		return e
	case *ConfigChange:
		return e
	case *CwdChanged:
		return e
	case *FileChanged:
		return e
	case *WorktreeCreate:
		return e
	case *WorktreeRemove:
		return e
	case *PreCompact:
		return e
	case *PostCompact:
		return e
	case *Elicitation:
		return e
	case *ElicitationResult:
		return e
	case *RawEvent:
		return e
	default:
		panic("claude: value missing envelopeAccessor")
	}
}

var (
	_ envelopeAccessor = (*SessionStart)(nil)
	_ envelopeAccessor = (*Setup)(nil)
	_ envelopeAccessor = (*SessionEnd)(nil)
	_ envelopeAccessor = (*UserPromptSubmit)(nil)
	_ envelopeAccessor = (*UserPromptExpansion)(nil)
	_ envelopeAccessor = (*PreToolUse)(nil)
	_ envelopeAccessor = (*PostToolUse)(nil)
	_ envelopeAccessor = (*PostToolUseFailure)(nil)
	_ envelopeAccessor = (*PostToolBatch)(nil)
	_ envelopeAccessor = (*PermissionRequest)(nil)
	_ envelopeAccessor = (*PermissionDenied)(nil)
	_ envelopeAccessor = (*SubagentStart)(nil)
	_ envelopeAccessor = (*SubagentStop)(nil)
	_ envelopeAccessor = (*TaskCreated)(nil)
	_ envelopeAccessor = (*TaskCompleted)(nil)
	_ envelopeAccessor = (*Stop)(nil)
	_ envelopeAccessor = (*StopFailure)(nil)
	_ envelopeAccessor = (*TeammateIdle)(nil)
	_ envelopeAccessor = (*Notification)(nil)
	_ envelopeAccessor = (*MessageDisplay)(nil)
	_ envelopeAccessor = (*InstructionsLoaded)(nil)
	_ envelopeAccessor = (*ConfigChange)(nil)
	_ envelopeAccessor = (*CwdChanged)(nil)
	_ envelopeAccessor = (*FileChanged)(nil)
	_ envelopeAccessor = (*WorktreeCreate)(nil)
	_ envelopeAccessor = (*WorktreeRemove)(nil)
	_ envelopeAccessor = (*PreCompact)(nil)
	_ envelopeAccessor = (*PostCompact)(nil)
	_ envelopeAccessor = (*Elicitation)(nil)
	_ envelopeAccessor = (*ElicitationResult)(nil)
	_ envelopeAccessor = (*RawEvent)(nil)
)

func (e *SessionStart) envelopePtr() *Envelope        { return &e.Envelope }
func (e *Setup) envelopePtr() *Envelope               { return &e.Envelope }
func (e *SessionEnd) envelopePtr() *Envelope          { return &e.Envelope }
func (e *UserPromptSubmit) envelopePtr() *Envelope    { return &e.Envelope }
func (e *UserPromptExpansion) envelopePtr() *Envelope { return &e.Envelope }
func (e *PreToolUse) envelopePtr() *Envelope          { return &e.Envelope }
func (e *PostToolUse) envelopePtr() *Envelope         { return &e.Envelope }
func (e *PostToolUseFailure) envelopePtr() *Envelope  { return &e.Envelope }
func (e *PostToolBatch) envelopePtr() *Envelope       { return &e.Envelope }
func (e *PermissionRequest) envelopePtr() *Envelope   { return &e.Envelope }
func (e *PermissionDenied) envelopePtr() *Envelope    { return &e.Envelope }
func (e *SubagentStart) envelopePtr() *Envelope       { return &e.Envelope }
func (e *SubagentStop) envelopePtr() *Envelope        { return &e.Envelope }
func (e *TaskCreated) envelopePtr() *Envelope         { return &e.Envelope }
func (e *TaskCompleted) envelopePtr() *Envelope       { return &e.Envelope }
func (e *Stop) envelopePtr() *Envelope                { return &e.Envelope }
func (e *StopFailure) envelopePtr() *Envelope         { return &e.Envelope }
func (e *TeammateIdle) envelopePtr() *Envelope        { return &e.Envelope }
func (e *Notification) envelopePtr() *Envelope        { return &e.Envelope }
func (e *MessageDisplay) envelopePtr() *Envelope      { return &e.Envelope }
func (e *InstructionsLoaded) envelopePtr() *Envelope  { return &e.Envelope }
func (e *ConfigChange) envelopePtr() *Envelope        { return &e.Envelope }
func (e *CwdChanged) envelopePtr() *Envelope          { return &e.Envelope }
func (e *FileChanged) envelopePtr() *Envelope         { return &e.Envelope }
func (e *WorktreeCreate) envelopePtr() *Envelope      { return &e.Envelope }
func (e *WorktreeRemove) envelopePtr() *Envelope      { return &e.Envelope }
func (e *PreCompact) envelopePtr() *Envelope          { return &e.Envelope }
func (e *PostCompact) envelopePtr() *Envelope         { return &e.Envelope }
func (e *Elicitation) envelopePtr() *Envelope         { return &e.Envelope }
func (e *ElicitationResult) envelopePtr() *Envelope   { return &e.Envelope }
func (e *RawEvent) envelopePtr() *Envelope            { return &e.Envelope }
