# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/).

## [Unreleased]

### Added

- `wat` CLI with subcommands `init`, `install`, `run`, `port`, `test`, and `doctor`; root and per-command help; `--agent`, `--event`, and `--fail-closed` flags on the subcommands that need them; stub commands exit 2 until implemented
- Unified `agenthooks.Event` / `Kind` types and tool-name normalization for hook authors
- Unified `agenthooks.Result` / `Decision` types with `Merge` and `Unsupported` for hook handler responses
- `agenthooks.ParseDialect` and `agenthooks.Detect` for hook payload dialect sniffing
- `agenthooks.Codec` and `agenthooks.CodecFor` for dialect codec lookup
- `agenthooks.ClaudeCodec` for Claude Code stdin/stdout translation
- `agenthooks.CopilotCodec` for GitHub Copilot camelCase and VS Code hook payloads
- `agenthooks.CursorCodec` decodes and encodes Cursor hooks, folding dedicated shell/MCP/file events into unified pre/post tool kinds
- `agenthooks.Mux` with `On`, `OnAny`, `Serve`, and `Main` for hook handler registration and the hook process lifecycle
- Serve options `WithDialect`, `WithEvent`, and `WithGetenv` for explicit dialect, Copilot event hints, and testable environment lookup
- `agenthooks/portconfig.Parse` and `Emit` for Claude Code, GitHub Copilot, and Cursor hook configuration files, preserving unmappable entries and native handler fields for same-dialect round-trip
- `agenthooks/portconfig.Translate` for cross-agent hook config conversion with explicit warnings for lossy matchers, unsupported handler types, and unmappable events
- `claudehook` package with typed Claude Code hook events, `Decode`/`Encode`, generic `Mux`/`On`/`Serve`/`Main` (`On` panics on duplicate handler registration), and `claudehook/tools` lazy tool input helpers
- `agenthooks.As` to re-decode unified events into native `claudehook` types for long-tail Claude events
- `copilothook` package with typed GitHub Copilot hook events, dual-format `Decode`/`Encode`, generic `Mux`/`On`/`Serve`/`Main` (`On` panics on duplicate handler registration), `WithEvent` for camelCase payloads, and `copilothook/tools` lazy tool input helpers
- `cursorhook` package with typed Cursor hook events (21 surfaces), `Decode`/`Encode` (`ErrEventNameRequired` when `hook_event_name` and `WithEvent` are both absent), generic `Mux`/`On`/`Serve`/`Main`, `WithEvent` for payloads missing `hook_event_name`, `cursorhook/tools` lazy tool input helpers, and env helpers (`CURSOR_PROJECT_DIR`, `CURSOR_VERSION`)
- `agenthooks.AsCursor` to re-decode unified events into native `cursorhook` types
- `claudehook.EnvelopeOf` and decode error sentinels (`ErrEmptyPayload`, `ErrDecodePayload`) for stable envelope access and error handling
- `claudehook.Handler.TimeoutSeconds` for hook config timeout lookup
- `claudehook.HandlerErrorExit`, `claudehook.FailBlockExit`, `copilothook.HandlerErrorExit`, and `agenthooks.CopilotHandlerErrorExit` for named handler-error exit codes
- Decoder registry parity tests in `claudehook` and `copilothook` (`envelope_meta_test.go`)
- Typed decode failures in `copilothook` and `cursorhook` wrap `ErrDecodePayload` for stable `errors.Is` checks
