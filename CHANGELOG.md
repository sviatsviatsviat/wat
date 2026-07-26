# Changelog

All notable user-visible changes to wat are documented here.

The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).
The project intends to use [Semantic Versioning](https://semver.org/).

## [Unreleased]

### Added

- Cursor hook `Envelope` decodes optional common-schema fields `model_id` and
  `model_params`, using the exported `ModelParam` element type.
- Cursor `afterShellExecution` decodes `sandbox`. Cursor `postToolUseFailure`
  decodes `is_interrupt`. `afterShellExecution`, `afterMCPExecution`,
  `postToolUse`, and `postToolUseFailure` prefer documented `duration` via
  `DurationMillis()` when that field is present (including explicit `0`),
  falling back to `duration_ms` only when `duration` is absent.
- Cursor `PreToolUse` decodes optional input `agent_message` (pre-call
  narrative from the agent).

### Changed

- Cursor post-tool events match Hooks docs observe-only contracts:
  `afterFileEdit`, `afterShellExecution`, `afterMCPExecution`, and
  `postToolUseFailure` no longer emit host-consumed JSON. Portable
  `OnPostTool` still expands to `afterFileEdit` and `afterMCPExecution` for
  observation (`Context` / `WithUpdatedOutput` are no-ops there) but no longer
  expands to `afterShellExecution` (use `sdk/cursor.AfterShellExecution`). MCP
  tool output rewrite stays on `postToolUse` (`updated_mcp_tool_output`).
  Portable `OnPostToolFailure` `Context` is ignored on Cursor; Claude and
  Copilot still apply recovery context. Cloud agents do not load MCP hooks.
- Cursor `PreToolUse` handlers now receive `PreToolUseResults`. `Ask` still
  encodes `"permission":"ask"` for schema compatibility, but godoc and protocol
  docs warn that Cursor does not enforce ask for `preToolUse` today; prefer
  `Allow` or `Deny` when gating is required. Protocol docs also clarify that
  Cursor enforces `"ask"` on `beforeShellExecution` and `beforeMCPExecution`.

### Fixed

- Cursor `BeforeReadFile` `Deny` emits `permission` + `user_message` with exit 0
  (no `agent_message` / exit 2), matching Cursor's beforeReadFile schema.
  `Ask` is coerced to the same deny encoding because the host does not support
  `"ask"` on this event.
- Cursor `BeforeTabFileRead` `Allow`/`Deny` emit `permission` allow|deny only
  with exit 0 (no ask, no `user_message` / `agent_message`), matching Cursor's
  beforeTabFileRead schema.

## [v0.1.1-alpha] - 2026-07-25

### Fixed

- Cursor `SubagentStart` `Deny` emits `permission` + `user_message` with exit 0
  (no `agent_message` / exit 2), matching Cursor's subagentStart schema so the
  host does not re-wrap stdout as the user message.

## [v0.1.0-alpha] - 2026-07-25

### Added

- A `wat` CLI for the complete hook-project lifecycle:
  - `wat init [path]` creates an independent `.wat/` Go module with a minimal
    portable `OnSessionStart` example, starter `.wat/testdata/` fixtures with
    expect sidecars, and preserves existing files unless `--force` is
    explicitly used for `hooks.go`.
  - `wat install` inspects authored registrations and reconciles only the
    matching wat-managed entries in Claude Code, GitHub Copilot, and Cursor
    project configuration. Unrelated hooks are preserved and stale wat-managed
    events are removed.
  - `wat run` builds a content-addressed hook executable on demand, reuses it on
    warm invocations, and passes the native hook process streams and exit code
    through. `--fail-closed` maps build failures to exit 2. Hook module path
    resolution uses the package's owning module so builds work inside a
    `go.work` workspace that lists multiple modules. The cache key hashes
    `.wat/` sources except `.cache/` and `testdata/`.
  - `wat test` runs a registered native event against a JSON fixture, including
    stdin fixtures, and reports the dialect, event, stdout, recognized decision,
    and exit code. Optional `<fixture>.expect.json` sidecars (or `--expect`)
    assert exit code, decision, and stdout; a matching expect run exits 0.
  - `wat doctor` checks the Go toolchain, project files, compilation, cache,
    authored registration manifest, and installed configurations, with
    actionable fixes and exit 4 when a check fails. Missing install wiring and
    `wat` not on `PATH` are warnings (hooks will not run until installed).
    Status labels are colored on a TTY.

- Upward `.wat/` project discovery for CLI commands, with
  `WAT_PROJECT_DIR` available to select a workspace root explicitly.

- A portable `sdk/agnostic` hook API for Claude Code, GitHub Copilot, and
  Cursor:
  - typed session start/end, user prompt, pre/post tool, post-tool failure,
    subagent start/stop, stop, and pre-compact events;
  - fluent `UseHooks().On*` registration that expands to the correct native
    events for each agent, including Cursor's dedicated shell, MCP, read, and
    edit hooks;
  - hook-scoped result builders for portable allow/deny/ask, context injection,
    input/output updates, and stop follow-up behavior;
  - observe-only handlers for events that cannot emit a portable response;
  - normalized event envelopes and tool calls that preserve the native event,
    native tool name, raw input, call ID, shell command, and MCP identity.

- Canonical typed tool inputs in `sdk/agnostic/tools` for shell, edit, write,
  read, glob, grep, task, web fetch, and web search operations.

- Native typed hook SDKs for Claude Code, GitHub Copilot, and Cursor, including
  native event constants and payloads, fluent `UseHooks` registration,
  hook-scoped response builders, advanced fluent output fields, typed tool
  inputs, stable decode errors, and native output/exit behavior.

- `sdk/cursor`'s `SubagentStart` event decodes the full native
  `subagentStart` payload: `parent_conversation_id`, `tool_call_id`,
  `subagent_model`, `is_parallel_worker`, and `git_branch`, in addition to
  the existing `subagent_id`, `subagent_type`, and `task`.

- This repository's own `.wat/hooks.go` now includes a real Cursor
  `subagentStart` hook that denies a subagent spawn pinned to a model other
  than Cursor's automatic selection.

- `sdk/run` for standalone hook executables and tooling:
  - `Serve` detects one native dialect, decodes a registered event once, invokes
    handlers in registration order, merges their typed results, stops on a
    terminal result, and encodes one native response;
  - `Inspect` returns the expanded native dialect/event manifest and handler
    counts without invoking handlers.

- Deterministic multi-handler composition. Additive fields merge, later
  replacement fields win with a stderr warning, deny/block and
  `continue: false` stop later handlers, and nil/zero results remain silent.
