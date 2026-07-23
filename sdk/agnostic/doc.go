// Package agnostic provides a unified hook event model for Claude Code,
// GitHub Copilot, and Cursor. It defines typed events and hook-response result
// types, and UseHooks registration that fans out adapter handlers onto each
// agent SDK. Canonical tool names and typed tool-input helpers live in
// sdk/agnostic/tools. The Agent field on each normalized event is a plain string
// matching each per-agent SDK Dialect constant (claude.Dialect, …).
//
// Inbound native→unified mapping and host result wrappers live in
// sdk/agnostic/internal/{claude,cursor,copilot} with shared types and result
// interfaces in sdk/agnostic/internal/model. Agnostic depends on the
// per-agent SDKs; those SDKs remain usable without agnostic.
//
// Hook handlers receive a normalized typed event (PreToolEvent, PostToolEvent,
// and others) plus hook-scoped result builders (PreToolResults, PostToolResults,
// and others) that wrap the native agent Results. Advanced fields use fluent
// With* methods on the returned result (for example WithUpdatedInput).
// Observe-only kinds use OnSessionEnd, OnUserPrompt, OnPreCompact, and
// OnSubagentStart with per-kind handler types. Register handlers with
// UseHooks().OnPreTool, UseHooks().OnPostTool, and related chain methods, then
// pass the chain to run.Serve from github.com/sviatsviatsviat/wat/sdk/run.
// Main merges chains that share a dialect before dispatch.
//
// Agent-only capabilities (BlockPrompt, Env, HaltSession, and others) belong in
// sdk/claude, sdk/copilot, and sdk/cursor.
package agnostic
