# Changelog

All notable user-visible changes to wat are documented here.

The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).
The project intends to use [Semantic Versioning](https://semver.org/).

## [Unreleased]

### Added

- Claude `UseHooks` registration for previously decode-only events: `Setup`,
  `InstructionsLoaded` (observe-only), `PostToolBatch`, `TeammateIdle`,
  `ConfigChange`, `CwdChanged`, `FileChanged`, `WorktreeRemove` (observe-only),
  `PostCompact` (observe-only), `ElicitationResult`, and `StopFailure`
  (observe-only). Install/doctor/`run.Inspect` advertise these when handlers
  are registered.
- Claude `UserPromptExpansion`, `PreCompact`, `PostToolUseFailure`,
  `PostToolBatch`, and `ConfigChange` expose event-scoped `Block` builders that
  encode JSON `decision:"block"` with exit 0 (`DecisionOutput`). `TeammateIdle` /
  `TaskCreated` / `TaskCompleted` `Block` uses Claude exit 2 with stderr
  feedback (`continue:false` still stops the teammate). Handler errors remain
  exit 1 (fail-open for most events). See `docs/agents/claude.md`.
- Claude `PermissionRequest` decodes optional `permission_suggestions` and can
  return `updatedPermissions` on allow via `WithUpdatedPermissions` (same entry
  shape: add/replace/remove rules, setMode, add/removeDirectories). Echoing a
  suggestion applies that "always allow" option without prompting.

### Fixed

- Claude `TaskCreated` / `TaskCompleted` decode flat `task_id`, `task_subject`,
  and related task fields instead of a nested `task` blob. `StopFailure`
  decodes `error`, `error_details`, and `last_assistant_message`.
  `Elicitation` / `ElicitationResult` use `mcp_server_name` and optional
  mode/url/id fields. `WorktreeCreate` decodes `name`. `FileChanged` decodes
  `event`, `ConfigChange` decodes optional `file_path`, and `PostCompact`
  decodes `compact_summary`.
- Claude `WorktreeCreate` encodes `Path` as plain stdout text for command hooks
  (the transport wat installs). JSON `hookSpecificOutput.worktreePath` was
  incorrect for that transport and could prevent the host from using the
  returned directory.

## [v0.2.0-alpha] - 2026-07-27

### Added

- `wat run --agent` / `--event` are forwarded to the hooks binary as dispatch
  hints: when set, Serve selects dialect/event from the hint (not from
  detection / `hook_event_name` peek) and may still inspect the payload only to
  warn on disagreements. Missing `hook_event_name` is allowed when `--event` is
  present. Mismatches do not fail the run. `wat test` forwards the same flags.
  `wat doctor` warns when a wat run command’s flags disagree with the native
  hook config map key or agent file.
- `wat version` (and root `-version` / `--version`) prints the module version
  string used for `wat init` pinning and the hook build cache key.
- Cursor hook `Envelope` decodes optional common-schema fields `model_id` and
  `model_params`, using the exported `ModelParam` element type.
- Cursor `afterShellExecution` decodes `sandbox`. Cursor `postToolUseFailure`
  decodes `is_interrupt`. `afterShellExecution`, `afterMCPExecution`,
  `postToolUse`, and `postToolUseFailure` prefer documented `duration` via
  `DurationMillis()` when that field is present (including explicit `0`),
  falling back to `duration_ms` only when `duration` is absent.
- Cursor `PreToolUse` decodes optional input `agent_message` (pre-call
  narrative from the agent).
- `sdk/cursor`'s `SubagentStop` event decodes the full native `subagentStop`
  telemetry payload: `description`, `duration_ms`, `message_count`,
  `tool_call_count`, and `modified_files`, in addition to the existing
  identity, status, summary, loop, and transcript fields.
- Cursor `AfterTabFileEdit` decodes Tab-specific edit detail (`range`,
  `old_line`, `new_line`) via `TabEdit` / `EditRange`, while Agent
  `AfterFileEdit` continues to use `Edit` with `old_string` / `new_string`.
- Cursor `WorkspaceOpen` handlers can return `pluginPaths` (absolute plugin
  directories to load for the workspace). This is a desktop/CLI lifecycle hook
  and does not run in cloud agents.
- `sdk/cursor`'s `PreCompact` event decodes the full native `preCompact`
  compaction metrics: `context_usage_percent`, `context_tokens`,
  `context_window_size`, `message_count`, `messages_to_compact`, and
  `is_first_compaction`, in addition to `trigger`. Observational
  `user_message` output is unchanged. Portable `OnPreCompact` remains
  observe-only and does not expose Cursor-only metrics or `user_message`.
- Cursor `AfterAgentThought` decodes optional `duration_ms` for the completed
  thinking block. The event remains observe-only.
- Cursor `SessionStart` decodes optional `composer_mode` (`"agent"`, `"ask"`,
  or `"edit"`). The hook is fire-and-forget (the host does not enforce
  `continue` / `user_message`) and is not available for cloud agents.
- `sdk/cursor`'s `SessionEnd` event decodes the full native `sessionEnd`
  payload: `duration_ms`, `final_status`, and `error_message`, in addition to
  the existing `reason` and `is_background_agent` fields. The event remains
  observe-only.

### Changed

- Cursor `AfterFileEdit`, `AfterShellExecution`, `AfterMCPExecution`, and
  `PostToolUseFailure` handlers are now observe-only. Existing handlers must
  migrate from `(context.Context, Event, PostToolResults) (PostToolOutput,
  error)` to `(context.Context, Event) error`. These events no longer emit
  host-consumed JSON.
- Cursor `WorkspaceOpen` handlers are now result-producing. Existing
  `(context.Context, WorkspaceOpen) error` handlers must accept
  `WorkspaceOpenResults` and return `(WorkspaceOpenOutput, error)`.
- Portable `OnPostTool` still expands to `afterFileEdit` and
  `afterMCPExecution` for observation (`Context` / `WithUpdatedOutput` are
  no-ops there), but no longer expands to `afterShellExecution`; use
  `sdk/cursor.AfterShellExecution` for native auditing. MCP tool output rewrite
  stays on `postToolUse` (`updated_mcp_tool_output`). Portable
  `OnPostToolFailure` context is ignored on Cursor; Claude and Copilot still
  apply it.
- Cursor `PreToolUse` handlers now receive `PreToolUseResults`. `Ask` still
  encodes `"permission":"ask"` for schema compatibility, but Cursor does not
  enforce ask for `preToolUse` today; prefer `Allow` or `Deny` when gating is
  required. Cursor does enforce `"ask"` on `beforeShellExecution` and
  `beforeMCPExecution`.

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
