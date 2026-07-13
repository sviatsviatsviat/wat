// Package agnostic provides a unified hook event model for Claude Code,
// GitHub Copilot, and Cursor. It defines normalized Event and Kind types,
// hook-response Result and Decision types, Dialect identification via ParseDialect
// and Detect, tool-name normalization, typed tool-input helpers, and dialect
// codecs that translate native stdin/stdout JSON.
//
// Result is the outbound hook handler response. Event.Result is a different
// concept: incoming post-tool payload data on decoded events. Multiple handler
// results merge via Merge; Unsupported reports fields a dialect cannot encode.
// Use CodecFor to obtain a dialect codec; ClaudeCodec, CopilotCodec, and
// CursorCodec decode native stdin and encode Result for their respective agents.
//
// Hook scripts register handlers on a Mux with On and OnAny, then call Main
// (or Serve for testable I/O). Serve reads stdin, detects or accepts a dialect
// (WithDialect), decodes via CodecFor, runs OnAny then kind-specific handlers
// and merges results, encodes stdout, and returns the exit code. WithEvent
// supplies the Copilot camelCase event name when the payload omits it;
// WithGetenv injects environment lookup for Detect and ClaudeCodec encode.
package agnostic
