package cursor

import (
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/runtime"
)

// Dialect is the sdk/run registry name for Cursor hooks.
const Dialect = runtime.Dialect

var (
	// ErrEmptyPayload indicates stdin was empty.
	ErrEmptyPayload = runtime.ErrEmptyPayload
	// ErrDecodePayload indicates JSON parsing of the hook payload failed.
	ErrDecodePayload = runtime.ErrDecodePayload
	// ErrEventNameRequired indicates the payload omitted hook_event_name.
	ErrEventNameRequired = runtime.ErrEventNameRequired
)

// HandlerErrorExit is exit code 1. The runner should use this when a handler
// returns an error under Cursor's default fail-open policy.
const HandlerErrorExit = event.HandlerErrorExit

// PermissionDenyExit is exit code 2. Cursor treats it as block/deny on permission-gating
// events, equivalent to returning permission:"deny".
const PermissionDenyExit = event.PermissionDenyExit

// CanonicalEventName reports whether name is a known Cursor hook event name.
// Known names are returned unchanged; unknown non-empty names return false.
func CanonicalEventName(name string) (string, bool) {
	return runtime.CanonicalEventName(name)
}

// EventAliasMap returns a copy of known event name to itself (identity map).
func EventAliasMap() map[string]string {
	return runtime.EventAliasMap()
}
