// Package agnostic provides a unified hook event model for Claude Code,
// GitHub Copilot, and Cursor. It defines normalized Event and Kind types,
// kind-specific hook-response result types, and On* registration that fans out
// adapter handlers onto each agent SDK Chain. Canonical tool names and typed
// tool-input helpers live in sdk/agnostic/tools. Event.Agent is a plain string
// matching each per-agent SDK Dialect constant (claude.Dialect, …).
//
// Inbound native→unified mapping and host result wrappers live in
// sdk/agnostic/internal/{claude,cursor,copilot} with shared types and result
// interfaces in sdk/agnostic/internal/model. Agnostic depends on the
// per-agent SDKs; those SDKs remain usable without agnostic.
//
// Hook handlers receive a hook context wrapper (PreToolHook, PostToolHook, and
// others) embedding a normalized typed event (PreToolEvent, PostToolEvent, …)
// plus hook-scoped result builders (PreToolResults, PostToolResults, and others)
// that wrap the native agent Results. Advanced fields use fluent With* methods
// on the returned result (for example WithUpdatedInput). Observe-only kinds use
// OnSessionEnd, OnUserPrompt, OnPreCompact, and OnSubagentStart with
// per-kind handler types. Register handlers with OnPreTool, OnPostTool, and the
// other typed On methods, then call run.Main from
// github.com/sviatsviatsviat/wat/sdk/run.
//
// Agent-only capabilities (BlockPrompt, Env, HaltSession, and others) belong in
// sdk/claude, sdk/copilot, and sdk/cursor.
package agnostic
