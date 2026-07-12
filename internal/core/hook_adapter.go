package core

// HookAdapter carries host-specific parsed hook payload for [HookHandlerProvider] and [HookHandler].
// Capability views are obtained via As* methods;
type HookAdapter interface {
	AsAfterFileEdit() (AfterFileEditHook, bool)
	AsAfterShellExecution() (AfterShellExecutionHook, bool)
}
