// Package agnostic provides a unified hook event model for Claude Code,
// GitHub Copilot, and Cursor. It defines normalized Event and Kind types,
// hook-response Result and Decision types, Dialect identification via ParseDialect
// and Detect, tool-name normalization, typed tool-input helpers, and dialect
// codecs that translate native stdin/stdout JSON.
//
// Agent-specific codecs live under sdk/agnostic/claude, sdk/agnostic/copilot,
// and sdk/agnostic/cursor. Shared event and result types live in
// sdk/agnostic/internal/model.
//
// Result is the outbound hook handler response. Event.Result is a different
// concept: incoming post-tool payload data on decoded events. Multiple handler
// results merge via Merge; Unsupported reports fields a dialect cannot encode.
// Use CodecFor to obtain a dialect codec; ClaudeCodec, CopilotCodec, and
// CursorCodec decode native stdin and encode Result for their respective agents.
//
// Hook scripts register handlers with On and OnAny (chainable via Chain), then
// call run.Main from github.com/sviatsviatsviat/wat/sdk/run. Handlers register
// into a shared singleton registry alongside per-agent SDK handlers. WithDialect,
// WithEvent, and WithGetenv configure run.Main and run.Serve.
package agnostic
