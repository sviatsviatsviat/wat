package cursor

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
)

// envelopeAccessor provides compile-time-checked envelope metadata access.
type envelopeAccessor interface {
	envelopePtr() *Envelope
}

func (e *Envelope) setEnvelopeMeta(received, canonical string, raw json.RawMessage) {
	e.receivedName = received
	e.canonical = canonical
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
	case SessionEnd:
		return &e
	case BeforeSubmitPrompt:
		return &e
	case PreToolUse:
		return &e
	case PostToolUse:
		return &e
	case PostToolUseFailure:
		return &e
	case BeforeShellExecution:
		return &e
	case AfterShellExecution:
		return &e
	case BeforeMCPExecution:
		return &e
	case AfterMCPExecution:
		return &e
	case BeforeReadFile:
		return &e
	case AfterFileEdit:
		return &e
	case SubagentStart:
		return &e
	case SubagentStop:
		return &e
	case Stop:
		return &e
	case PreCompact:
		return &e
	case AfterAgentResponse:
		return &e
	case AfterAgentThought:
		return &e
	case BeforeTabFileRead:
		return &e
	case AfterTabFileEdit:
		return &e
	case WorkspaceOpen:
		return &e
	case RawEvent:
		return &e
	default:
		panic("cursor: event missing envelopeAccessor")
	}
}

func envelopeAccessorForValue(v any) envelopeAccessor {
	switch e := v.(type) {
	case *SessionStart:
		return e
	case *SessionEnd:
		return e
	case *BeforeSubmitPrompt:
		return e
	case *PreToolUse:
		return e
	case *PostToolUse:
		return e
	case *PostToolUseFailure:
		return e
	case *BeforeShellExecution:
		return e
	case *AfterShellExecution:
		return e
	case *BeforeMCPExecution:
		return e
	case *AfterMCPExecution:
		return e
	case *BeforeReadFile:
		return e
	case *AfterFileEdit:
		return e
	case *SubagentStart:
		return e
	case *SubagentStop:
		return e
	case *Stop:
		return e
	case *PreCompact:
		return e
	case *AfterAgentResponse:
		return e
	case *AfterAgentThought:
		return e
	case *BeforeTabFileRead:
		return e
	case *AfterTabFileEdit:
		return e
	case *WorkspaceOpen:
		return e
	case *RawEvent:
		return e
	default:
		panic("cursor: value missing envelopeAccessor")
	}
}

var (
	_ envelopeAccessor = (*SessionStart)(nil)
	_ envelopeAccessor = (*SessionEnd)(nil)
	_ envelopeAccessor = (*BeforeSubmitPrompt)(nil)
	_ envelopeAccessor = (*PreToolUse)(nil)
	_ envelopeAccessor = (*PostToolUse)(nil)
	_ envelopeAccessor = (*PostToolUseFailure)(nil)
	_ envelopeAccessor = (*BeforeShellExecution)(nil)
	_ envelopeAccessor = (*AfterShellExecution)(nil)
	_ envelopeAccessor = (*BeforeMCPExecution)(nil)
	_ envelopeAccessor = (*AfterMCPExecution)(nil)
	_ envelopeAccessor = (*BeforeReadFile)(nil)
	_ envelopeAccessor = (*AfterFileEdit)(nil)
	_ envelopeAccessor = (*SubagentStart)(nil)
	_ envelopeAccessor = (*SubagentStop)(nil)
	_ envelopeAccessor = (*Stop)(nil)
	_ envelopeAccessor = (*PreCompact)(nil)
	_ envelopeAccessor = (*AfterAgentResponse)(nil)
	_ envelopeAccessor = (*AfterAgentThought)(nil)
	_ envelopeAccessor = (*BeforeTabFileRead)(nil)
	_ envelopeAccessor = (*AfterTabFileEdit)(nil)
	_ envelopeAccessor = (*WorkspaceOpen)(nil)
	_ envelopeAccessor = (*RawEvent)(nil)
)

func (e *SessionStart) envelopePtr() *Envelope         { return &e.Envelope }
func (e *SessionEnd) envelopePtr() *Envelope           { return &e.Envelope }
func (e *BeforeSubmitPrompt) envelopePtr() *Envelope   { return &e.Envelope }
func (e *PreToolUse) envelopePtr() *Envelope           { return &e.Envelope }
func (e *PostToolUse) envelopePtr() *Envelope          { return &e.Envelope }
func (e *PostToolUseFailure) envelopePtr() *Envelope   { return &e.Envelope }
func (e *BeforeShellExecution) envelopePtr() *Envelope { return &e.Envelope }
func (e *AfterShellExecution) envelopePtr() *Envelope  { return &e.Envelope }
func (e *BeforeMCPExecution) envelopePtr() *Envelope   { return &e.Envelope }
func (e *AfterMCPExecution) envelopePtr() *Envelope    { return &e.Envelope }
func (e *BeforeReadFile) envelopePtr() *Envelope       { return &e.Envelope }
func (e *AfterFileEdit) envelopePtr() *Envelope        { return &e.Envelope }
func (e *SubagentStart) envelopePtr() *Envelope        { return &e.Envelope }
func (e *SubagentStop) envelopePtr() *Envelope         { return &e.Envelope }
func (e *Stop) envelopePtr() *Envelope                 { return &e.Envelope }
func (e *PreCompact) envelopePtr() *Envelope           { return &e.Envelope }
func (e *AfterAgentResponse) envelopePtr() *Envelope   { return &e.Envelope }
func (e *AfterAgentThought) envelopePtr() *Envelope    { return &e.Envelope }
func (e *BeforeTabFileRead) envelopePtr() *Envelope    { return &e.Envelope }
func (e *AfterTabFileEdit) envelopePtr() *Envelope     { return &e.Envelope }
func (e *WorkspaceOpen) envelopePtr() *Envelope        { return &e.Envelope }
func (e *RawEvent) envelopePtr() *Envelope             { return &e.Envelope }
