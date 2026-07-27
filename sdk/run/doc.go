// Package run provides the hook process lifecycle and registration inspection
// for wat hook scripts.
// Hook authors build one or more UseHooks registrations on sdk/agnostic or a
// per-agent SDK, then pass them to Serve. Serve merges hooks that share a
// dialect, reads stdin, selects one agent dialect, peeks the event name,
// decodes the payload once, dispatches matching handlers, folds typed Output
// values via Merge/Stop, encodes once, and exits with the appropriate code.
// Inspect returns the same contributed handlers as a native registration
// manifest without invoking them.
//
// Dialect is auto-detected from the payload (and Cursor may also match via
// CURSOR_VERSION) unless --agent is passed on the process argv. Event selection
// peeks hook_event_name unless --event is passed; with --event, missing
// hook_event_name is allowed. Hint vs payload disagreements warn on stderr and
// still use the hint.
package run
