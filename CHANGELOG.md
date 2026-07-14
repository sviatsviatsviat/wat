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
- Kind-specific portable hook response types (`PreToolResult`, `PostToolResult`, `PostToolFailureResult`, `StopResult`, `SessionStartResult`) with hook-scoped builder interfaces (`PreToolResults`, `PostToolResults`, `StopResults`, and others) passed into typed `OnPreTool`, `OnPostTool`, and related registration methods; portable field matrix documented in `docs/agent-formats.md`
- `agnostic.ParseDialect` and `agnostic.Detect` for hook payload dialect sniffing
- `agnostic.ClaudeCodec` for Claude Code stdin/stdout translation
- `agnostic.CopilotCodec` for GitHub Copilot camelCase and VS Code hook payloads
- `agnostic.CursorCodec` decodes and encodes Cursor hooks, folding dedicated shell/MCP/file events into unified pre/post tool kinds
- Typed `OnPreTool`, `OnPostTool`, `OnStop`, and related handlers plus observe-only `OnAny`, `OnSessionEnd`, and chainable `Chain`; hook scripts call `run.Main` (reads `WAT_AGENT` / `WAT_EVENT` automatically)
- Serve options `WithDialect`, `WithEvent`, and `WithGetenv` for explicit dialect, Copilot event hints, and testable environment lookup
- `claude` package with typed Claude Code hook events, `Decode`/`Encode`, `On`/`Chain` registering into `sdk/run` (`On` panics on duplicate handler registration), and `sdk/claude/tools` lazy tool input helpers
- `copilot` package with typed GitHub Copilot hook events, dual-format `Decode`/`Encode`, `On`/`Chain` registering into `sdk/run` (`On` panics on duplicate handler registration), `WithEvent` for camelCase payloads, and `sdk/copilot/tools` lazy tool input helpers
- `cursor` package with typed Cursor hook events (21 surfaces), `Decode`/`Encode` (`ErrEventNameRequired` when `hook_event_name` and `WithEvent` are both absent), `On`/`Chain` registering into `sdk/run`, `WithEvent` for payloads missing `hook_event_name`, and `sdk/cursor/tools` lazy tool input helpers
- `claude.EnvelopeOf` and decode error sentinels (`ErrEmptyPayload`, `ErrDecodePayload`) for stable envelope access and error handling
- `claude.Handler.TimeoutSeconds` for hook config timeout lookup
- `claude.HandlerErrorExit` and `copilot.HandlerErrorExit` for handler-error exit codes
- `claude.FailBlockExit` for fail-closed blocking when `WithFailPolicy(FailBlock)` is active
- Decoder registry parity tests in `claude` and `copilot` (`envelope_meta_test.go`)
- Typed decode failures in `copilot` and `cursor` wrap `ErrDecodePayload` for stable `errors.Is` checks
