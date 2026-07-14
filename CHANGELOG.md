# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/).

## [Unreleased]

### Added

- `sdk/run` shared hook handler registry with `RegisterDialect`, `RegisterHandler`, `Serve`, and `Main`; agnostic and per-agent SDKs register into one singleton; fluent `On().On()` chaining via `Chain` in each SDK
- `wat` CLI with subcommands `init`, `install`, `run`, `port`, `test`, and `doctor`; root and per-command help; `--agent`, `--event`, and `--fail-closed` flags on the subcommands that need them
- `wat doctor` verifies Go toolchain, `.wat/` hook project compile, build cache, and installed hook entries; prints `PASS`/`FAIL`/`WARN` lines with fix suggestions; exits 4 when any check fails (warnings alone exit 0)
- `wat test --fixture` runs the user's hook script against a fixture JSON payload (file or stdin `-`); prints decoded unified event summary, hook stdout JSON, decision when present, and exit code; optional `--verbose` for expanded event fields and hook stderr; sample fixtures under `testdata/fixtures/`
- `wat port --from` / `--to` translates hook configs between Claude Code, GitHub Copilot, and Cursor via `portconfig.Translate`; `-i` / `-o` select input and output files (default input: `.claude/settings.json`, `.github/hooks/wat.json`, or `.cursor/hooks.json` from cwd); warnings print to stderr and exit 0; translation errors exit 3
- Unified `agnostic.Event` / `Kind` types and tool-name normalization for hook authors
- Unified `agnostic.Result` / `Decision` types with `Merge` and `Unsupported` for hook handler responses
- `agnostic.ParseDialect` and `agnostic.Detect` for hook payload dialect sniffing
- `agnostic.Codec` and `agnostic.CodecFor` for dialect codec lookup
- `agnostic.ClaudeCodec` for Claude Code stdin/stdout translation
- `agnostic.CopilotCodec` for GitHub Copilot camelCase and VS Code hook payloads
- `agnostic.CursorCodec` decodes and encodes Cursor hooks, folding dedicated shell/MCP/file events into unified pre/post tool kinds
- `agnostic.On`, `OnAny`, and chainable `Chain` replace `agnostic.Mux`; hook scripts call `run.Main` (reads `WAT_AGENT` / `WAT_EVENT` automatically)
- Serve options `WithDialect`, `WithEvent`, and `WithGetenv` for explicit dialect, Copilot event hints, and testable environment lookup
- `sdk/agnostic/portconfig.Parse` and `Emit` for Claude Code, GitHub Copilot, and Cursor hook configuration files, preserving unmappable entries and native handler fields for same-dialect round-trip
- `sdk/agnostic/portconfig.Translate` for cross-agent hook config conversion with explicit warnings for lossy matchers, unsupported handler types, and unmappable events
- `claude` package with typed Claude Code hook events, `Decode`/`Encode`, `On`/`Chain` registering into `sdk/run` (`On` panics on duplicate handler registration), and `sdk/claude/tools` lazy tool input helpers
- `copilot` package with typed GitHub Copilot hook events, dual-format `Decode`/`Encode`, `On`/`Chain` registering into `sdk/run` (`On` panics on duplicate handler registration), `WithEvent` for camelCase payloads, and `sdk/copilot/tools` lazy tool input helpers
- `cursor` package with typed Cursor hook events (21 surfaces), `Decode`/`Encode` (`ErrEventNameRequired` when `hook_event_name` and `WithEvent` are both absent), `On`/`Chain` registering into `sdk/run`, `WithEvent` for payloads missing `hook_event_name`, and `sdk/cursor/tools` lazy tool input helpers
- `claude.EnvelopeOf` and decode error sentinels (`ErrEmptyPayload`, `ErrDecodePayload`) for stable envelope access and error handling
- `claude.Handler.TimeoutSeconds` for hook config timeout lookup
- `claude.HandlerErrorExit` and `copilot.HandlerErrorExit` for handler-error exit codes
- `claude.FailBlockExit` for fail-closed blocking when `WithFailPolicy(FailBlock)` is active
- Decoder registry parity tests in `claude` and `copilot` (`envelope_meta_test.go`)
- Typed decode failures in `copilot` and `cursor` wrap `ErrDecodePayload` for stable `errors.Is` checks
