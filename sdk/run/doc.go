// Package run provides the shared hook handler registry and process lifecycle
// for wat hook scripts. Per-agent SDKs (claude, copilot, cursor) and the
// agnostic SDK register handlers via UseHooks into a registry (the process
// default from GetDefaultRegistry, or a caller-supplied *Registry); Main reads
// stdin, peeks the event name, decodes the payload once, dispatches all
// matching handlers with that event, merges their native JSON outputs, and
// exits with the appropriate code.
//
// WithDialect forces a dialect instead of auto-detection. WithGetenv injects
// environment lookup for dialect resolution and encode side effects. Payloads
// must include hook_event_name on the wire.
package run
