// Package hookkit provides shared helpers for per-agent hook SDKs and cmd/wat,
// including codecs, dialect handler bags, a HandlerQueue for deferred UseHooks
// installs, and typed handler adapters used by SDK bind code. Hooks
// contribution types (run.Hooks / run.Registry) live in sdk/run.
package hookkit
