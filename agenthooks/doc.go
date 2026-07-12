// Package agenthooks provides a unified hook event model for Claude Code,
// GitHub Copilot, and Cursor. It defines normalized Event and Kind types,
// hook-response Result and Decision types, Dialect identification via ParseDialect
// and Detect, tool-name normalization, typed tool-input helpers, and dialect
// codecs that translate native stdin/stdout JSON.
//
// Result is the outbound hook handler response. Event.Result is a different
// concept: incoming post-tool payload data on decoded events. Multiple handler
// results merge via Merge; Unsupported reports fields a dialect cannot encode.
// Use CodecFor to obtain a dialect codec and ClaudeCodec for Claude Code hooks.
package agenthooks
