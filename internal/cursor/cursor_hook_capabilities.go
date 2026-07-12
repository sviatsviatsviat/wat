package cursor

import "github.com/sviatsviatsviat/wat/internal/core"

// afterFileEditHookBridge implements [core.AfterFileEditHook] for Cursor afterFileEdit / afterTabFileEdit payloads.
type afterFileEditHookBridge struct {
	a *cursorHook[AfterFileEditFields]
}

func (b *afterFileEditHookBridge) HookEventName() string {
	return b.a.CommonInput.HookEventName
}

func (b *afterFileEditHookBridge) TranscriptPath() *string {
	return b.a.CommonInput.TranscriptPath
}

func (b *afterFileEditHookBridge) WriteDefaultToHost() {
	b.a.writeDefaultHookResponse()
}

func (b *afterFileEditHookBridge) FilePath() string {
	return b.a.EventSpecificInput.FilePath
}

// afterShellExecutionHookBridge implements [core.AfterShellExecutionHook] for Cursor afterShellExecution payloads.
type afterShellExecutionHookBridge struct {
	a *cursorHook[AfterShellExecutionFields]
}

func (b *afterShellExecutionHookBridge) HookEventName() string {
	return b.a.CommonInput.HookEventName
}

func (b *afterShellExecutionHookBridge) TranscriptPath() *string {
	return b.a.CommonInput.TranscriptPath
}

func (b *afterShellExecutionHookBridge) WriteDefaultToHost() {
	b.a.writeDefaultHookResponse()
}

func (b *afterShellExecutionHookBridge) Duration() float32 {
	return b.a.EventSpecificInput.Duration
}

func (b *afterShellExecutionHookBridge) RawCommand() string {
	return b.a.EventSpecificInput.Command
}

func (b *afterShellExecutionHookBridge) Sandbox() bool {
	return b.a.EventSpecificInput.Sandbox
}

func (a *cursorHook[T]) AsAfterFileEdit() (core.AfterFileEditHook, bool) {
	switch x := any(a).(type) {
	case *cursorHook[AfterFileEditFields]:
		if x.EventSpecificInput == nil {
			return nil, false
		}
		return &afterFileEditHookBridge{a: x}, true
	default:
		return nil, false
	}
}

func (a *cursorHook[T]) AsAfterShellExecution() (core.AfterShellExecutionHook, bool) {
	switch x := any(a).(type) {
	case *cursorHook[AfterShellExecutionFields]:
		if x.EventSpecificInput == nil {
			return nil, false
		}
		return &afterShellExecutionHookBridge{a: x}, true
	default:
		return nil, false
	}
}
