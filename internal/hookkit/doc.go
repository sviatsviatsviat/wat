// Package hookkit provides shared helpers for per-agent hook SDKs:
// codecs and dialect handler bags, a HandlerQueue for deferred UseHooks
// installs, typed handler adapters, typed output merge helpers, tool input
// helpers, and JSON/shell utilities. Native hooks.json handler encoding lives
// in cmd/wat/internal/hookconfig. Hooks contribution types (run.Hooks /
// run.Registry) live in sdk/run.
package hookkit
