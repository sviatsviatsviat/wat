// Package agnostic provides a unified hook event model for Claude Code,
// GitHub Copilot, and Cursor. It defines normalized Event and Kind types,
// kind-specific hook-response result types, Dialect identification via
// ParseDialect and Detect, tool-name normalization, typed tool-input helpers, and
// dialect codecs that translate native stdin/stdout JSON.
//
// Agent-specific codecs live under sdk/agnostic/claude, sdk/agnostic/copilot,
// and sdk/agnostic/cursor. Shared event and result types live in
// sdk/agnostic/internal/model.
//
// Hook handlers receive a hook context wrapper (PreToolHook, PostToolHook, and
// others) embedding a normalized typed event (PreToolEvent, PostToolEvent, …)
// plus hook-scoped result builders (PreToolResults, PostToolResults, and others).
// Advanced fields use fluent With* methods on the returned result (for example
// WithUpdatedInput). Observe-only kinds use OnSessionEnd, OnUserPrompt,
// OnPreCompact, OnSubagentStart, or OnAny with per-kind handler types. Register
// handlers with OnPreTool, OnPostTool, and the other typed On methods, then call
// run.Main from github.com/sviatsviatsviat/wat/sdk/run.
//
// Agent-only capabilities (BlockPrompt, Env, HaltSession, and others) belong in
// sdk/claude, sdk/copilot, and sdk/cursor.
package agnostic
