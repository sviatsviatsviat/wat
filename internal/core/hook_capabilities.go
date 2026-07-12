package core

// HookBasicInfo is informational data shared across hook capability interfaces.
type HookBasicInfo interface {
	HookEventName() string
	TranscriptPath() *string
}

// AfterFileEditHook is the exec subcommand view of afterFileEdit / afterTabFileEdit payloads.
type AfterFileEditHook interface {
	HookBasicInfo
	WriteDefaultToHost()
	FilePath() string
}

// AfterShellExecutionHook is the exec subcommand view of afterShellExecution payloads.
type AfterShellExecutionHook interface {
	HookBasicInfo
	WriteDefaultToHost()
	Duration() float32
	RawCommand() string
	Sandbox() bool
}
