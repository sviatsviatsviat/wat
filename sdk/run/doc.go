// Package run provides the hook process lifecycle for wat hook scripts.
// Hook authors register handlers with UseHooks on sdk/agnostic or a per-agent
// SDK, then call Main. Main reads stdin, selects one agent dialect, peeks the
// event name, decodes the payload once, dispatches matching handlers, folds
// typed Output values via Merge/Stop, encodes once, and exits with the
// appropriate code.
//
// Dialect is auto-detected from the payload (and Cursor may also match via
// CURSOR_VERSION). Payloads must include hook_event_name on the wire.
package run
