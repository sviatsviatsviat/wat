# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/).

## [Unreleased]

### Added

- Shared `hookkit.Output` (`IsZero`, `AllowedEvents`, `Encode(eventName)`) aliased as `Output` in `sdk/claude`, `sdk/copilot`, and `sdk/cursor`; package `hookkit.Encoder` on `Codec` validates then calls `out.Encode`; Claude `SessionStart` `WithEnv` writes `CLAUDE_ENV_FILE` inside `sessionStartOutput.Encode`
- `sdk/run` shared hook handler registry with `GetDefaultRegistry`, `NewRegistry`, `EnsureDialect`, `RegisterHandler`, and `Main`; agnostic and per-agent SDKs register via `UseHooks` (optional registry; default when omitted) with fluent chaining; `Main` peeks the event name then decodes the payload once before dispatching `Producer`s with the decoded event
- `wat` CLI with subcommands `init`, `install`, `run`, `port`, `test`, and `doctor`; root and per-command help; `--agent`, `--event`, and `--fail-closed` flags on the subcommands that need them
- `wat doctor` verifies Go toolchain, `.wat/` hook project compile, build cache, and installed hook entries; prints `PASS`/`FAIL`/`WARN` lines with fix suggestions; exits 4 when any check fails (warnings alone exit 0)
- `wat test --fixture` runs the user's hook script against a fixture JSON payload (file or stdin `-`); requires `--agent`; prints fixture agent/event, hook stdout JSON, decision when present, and exit code; optional `--verbose` for hook stderr; sample fixtures under `testdata/fixtures/`
- `wat port --from` / `--to` translates hook configs between Claude Code, GitHub Copilot, and Cursor via `portconfig.Translate`; `-i` / `-o` select input and output files (default input: `.claude/settings.json`, `.github/hooks/wat.json`, or `.cursor/hooks.json` from cwd); warnings print to stderr and exit 0; translation errors exit 3
- Normalized `ToolCall.Name` values for hook authors; `ToolCall.Input` is a typed `tools.Input` (`sdk/agnostic/tools`) with `AsBash`, `AsWrite`, `AsEdit`, `AsRead`, and related accessors (path/`file_path` aliases)
- Kind-specific portable hook response types (`PreToolResult`, `PostToolResult`, `PostToolFailureResult`, `StopResult`, `SessionStartResult`) with hook-scoped builder interfaces (`PreToolResults`, `PostToolResults`, `StopResults`, and others) passed into typed `OnPreTool`, `OnPostTool`, and related registration methods; one exported builder per `On*` method; advanced fields via fluent `With*` (for example `WithUpdatedInput`); portable field matrix documented in `docs/agent-formats.md`
- Hook output and portable result types in `sdk/agnostic` are constructed via `On*`-injected `*Results` builders (and `With*`); `nil` is the no-op; host-specific wrappers live in `sdk/agnostic/internal/{claude,cursor,copilot}` and implement interfaces from `internal/model`
- Agnostic `UseHooks` registration fans adapter handlers onto each agent SDK via `UseHooks(r)` (native encode owned by `sdk/claude`, `sdk/copilot`, `sdk/cursor`; decode is internal to each SDK and `sdk/run`); inbound mappers in `sdk/agnostic/internal/{claude,cursor,copilot}` convert native payloads to unified events
- Claude `OnElicitation`, `OnMessageDisplay`, and `OnWorktreeCreate` accept hook-scoped result builders (`ElicitationResults`, `MessageDisplayResults`, `WorktreeCreateResults`)
- `run.Invocation` and hook wrappers: agnostic kind types (`PreToolHook`, …) and `run.Hook[E]`, carrying typed events plus serve-time settings (`Dialect`, `Getenv`, `DialectConfig`); unified handler signatures `(ctx, hook, results)` for result hooks and `(ctx, hook) error` for observe hooks across all four SDKs
- Agnostic normalized typed events (`PreToolEvent`, `PostToolEvent`, `SessionEndEvent`, …) and per-kind handler types; expanded `On*` coverage (copilot 13/13, cursor 21/21, claude `OnSessionEnd`)
- Typed `UseHooks().OnPreTool`, `OnPostTool`, `OnStop`, and related handlers plus observe-only `OnSessionEnd` with fluent chaining; hook scripts call `run.Main` (reads `WAT_AGENT` automatically)
- Main options `run.WithDialect` (string dialect name such as `claude.Dialect`) and `run.WithGetenv` for explicit dialect and testable environment lookup; payloads must include `hook_event_name` (Claude, Copilot, and Cursor)
- Per-agent `Dialect` string constants (`claude.Dialect`, `copilot.Dialect`, `cursor.Dialect`) for `sdk/run` registration and `Envelope.Agent`
- `claude` package with typed Claude Code hook events, package-internal encode, `UseHooks` chain methods with hook-scoped result builders registering into `sdk/run`, and `sdk/claude/tools` event-bound tool input (`tools.Input` with `AsBash`, `AsWrite`, and related accessors)
- `copilot` package with typed GitHub Copilot hook events (PascalCase `hook_event_name`, snake_case fields and package-internal encode), `UseHooks` chain methods with hook-scoped result builders registering into `sdk/run`, and `sdk/copilot/tools` event-bound tool input (`tools.Input` with `AsBash`, `AsCreate`, and related accessors)
- `cursor` package with typed Cursor hook events (21 surfaces), package-internal encode (`ErrEventNameRequired` when `hook_event_name` is absent), `UseHooks` chain methods with hook-scoped result builders registering into `sdk/run`, and `sdk/cursor/tools` event-bound tool input (`tools.Input` with `AsShell`, `AsRead`, and related accessors)
- Decode error sentinels (`ErrEmptyPayload`, `ErrDecodePayload`) for stable error handling; shared envelope fields are embedded on each event type (read via promoted fields such as `SessionID` / `Cwd`)
- Shared `run.Event` (`EventName()` only; defined in `internal/hookkit`) and `run.Hook` for typed handler context; `run.Codec.Decode` returns `Event`
- `claude.Handler.TimeoutSeconds` for hook config timeout lookup
- `claude.HandlerErrorExit` and `copilot.HandlerErrorExit` for handler-error exit codes
- Decoder registry parity tests in `claude` and `copilot` (`sdk/claude/envelope_test.go`, `sdk/copilot/envelope_test.go`)
- Typed decode failures in `copilot` and `cursor` wrap `ErrDecodePayload` for stable `errors.Is` checks
- `claude.ErrEventNameRequired` when `hook_event_name` is absent (aligned with Copilot and Cursor)
- Copilot wire `Stop` always decodes as `AgentStop` (optional `agent_name` / `agent_display_name` via `IsSubagent`); portable `OnStop` skips subagent-scoped `Stop`, and portable `OnSubagentStop` receives those payloads
- Per-agent decode returns an error if `hook_event_name` is not a registered typed event (Main never decodes unknown names when no handlers are registered)
