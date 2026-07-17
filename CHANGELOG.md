# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/).

## [Unreleased]

### Added

- `sdk/run` shared hook handler registry with `RegisterDialect`, `RegisterHandler`, `Serve`, and `Main`; agnostic and per-agent SDKs register into one singleton; fluent `Chain` methods in each SDK
- `wat` CLI with subcommands `init`, `install`, `run`, `port`, `test`, and `doctor`; root and per-command help; `--agent`, `--event`, and `--fail-closed` flags on the subcommands that need them
- `wat doctor` verifies Go toolchain, `.wat/` hook project compile, build cache, and installed hook entries; prints `PASS`/`FAIL`/`WARN` lines with fix suggestions; exits 4 when any check fails (warnings alone exit 0)
- `wat test --fixture` runs the user's hook script against a fixture JSON payload (file or stdin `-`); requires `--agent` (and `--event` for Copilot camelCase fixtures); prints fixture agent/event, hook stdout JSON, decision when present, and exit code; optional `--verbose` for hook stderr; sample fixtures under `testdata/fixtures/`
- `wat port --from` / `--to` translates hook configs between Claude Code, GitHub Copilot, and Cursor via `portconfig.Translate`; `-i` / `-o` select input and output files (default input: `.claude/settings.json`, `.github/hooks/wat.json`, or `.cursor/hooks.json` from cwd); warnings print to stderr and exit 0; translation errors exit 3
- Unified `agnostic.Event` with normalized `ToolCall.Name` values for hook authors; `ToolCall.Input` is a typed `tools.Input` (`sdk/agnostic/tools`) with `AsBash`, `AsWrite`, `AsEdit`, `AsRead`, and related accessors (path/`file_path` aliases)
- Kind-specific portable hook response types (`PreToolResult`, `PostToolResult`, `PostToolFailureResult`, `StopResult`, `SessionStartResult`) with hook-scoped builder interfaces (`PreToolResults`, `PostToolResults`, `StopResults`, and others) passed into typed `OnPreTool`, `OnPostTool`, and related registration methods; one exported builder per Chain/`On*` method; advanced fields via fluent `With*` (for example `WithUpdatedInput`); portable field matrix documented in `docs/agent-formats.md`
- Hook output and portable result types in `sdk/agnostic` are constructed via Chain/`On*`-injected `*Results` builders (and `With*`); `nil` is the no-op; host-specific wrappers live in `sdk/agnostic/internal/{claude,cursor,copilot}` and implement interfaces from `internal/model`
- Agnostic `On*` registration fans adapter handlers onto each agent SDK `Chain` (native Decode/Encode owned by `sdk/claude`, `sdk/copilot`, `sdk/cursor`); inbound mappers in `sdk/agnostic/internal/{claude,cursor,copilot}` convert native payloads to unified events
- Claude `Chain.Elicitation`, `MessageDisplay`, and `WorktreeCreate` accept hook-scoped result builders (`ElicitationResults`, `MessageDisplayResults`, `WorktreeCreateResults`)
- `run.Invocation` and hook wrappers: agnostic kind types (`PreToolHook`, …) and per-agent generic `Hook[E]`, carrying typed events plus serve-time settings (`Dialect`, `EventHint`, `Getenv`, `DialectConfig`) and `Raw()` for untouched native JSON; unified handler signatures `(ctx, hook, results)` for result hooks and `(ctx, hook) error` for observe hooks across all four SDKs
- Agnostic normalized typed events (`PreToolEvent`, `PostToolEvent`, `SessionEndEvent`, …) and per-kind handler types replacing bare `*Event` in handler signatures; expanded Chain coverage (copilot 13/13, cursor 21/21, claude `SessionEnd`)
- Typed `OnPreTool`, `OnPostTool`, `OnStop`, and related handlers plus observe-only `OnSessionEnd` and chainable `Chain`; hook scripts call `run.Main` (reads `WAT_AGENT` / `WAT_EVENT` automatically)
- Serve options `run.WithDialect` (string dialect name such as `claude.Dialect`), `run.WithEvent`, and `run.WithGetenv` for explicit dialect, Copilot event hints, and testable environment lookup
- Per-agent `Dialect` string constants (`claude.Dialect`, `copilot.Dialect`, `cursor.Dialect`) for `sdk/run` registration and `Event.Agent`
- `claude` package with typed Claude Code hook events, `Decode`/`Encode`, `Chain` handlers with hook-scoped result builders registering into `sdk/run`, and `sdk/claude/tools` event-bound tool input (`tools.Input` with `AsBash`, `AsWrite`, and related accessors)
- `copilot` package with typed GitHub Copilot hook events, dual-format `Decode`/`Encode`, `Chain` handlers with hook-scoped result builders registering into `sdk/run`, `WithEvent` for camelCase payloads, and `sdk/copilot/tools` event-bound tool input (`tools.Input` with `AsBash`, `AsCreate`, and related accessors)
- `cursor` package with typed Cursor hook events (21 surfaces), `Decode`/`Encode` (`ErrEventNameRequired` when `hook_event_name` and `WithEvent` are both absent), `Chain` handlers with hook-scoped result builders registering into `sdk/run`, `WithEvent` for payloads missing `hook_event_name`, and `sdk/cursor/tools` event-bound tool input (`tools.Input` with `AsShell`, `AsRead`, and related accessors)
- `claude.EnvelopeOf` and decode error sentinels (`ErrEmptyPayload`, `ErrDecodePayload`) for stable envelope access and error handling
- `claude.Handler.TimeoutSeconds` for hook config timeout lookup
- `claude.HandlerErrorExit` and `copilot.HandlerErrorExit` for handler-error exit codes
- `claude.FailBlockExit` for fail-closed blocking when `WithFailPolicy(FailBlock)` is active
- Decoder registry parity tests in `claude` and `copilot` (`envelope_meta_test.go`)
- Typed decode failures in `copilot` and `cursor` wrap `ErrDecodePayload` for stable `errors.Is` checks
