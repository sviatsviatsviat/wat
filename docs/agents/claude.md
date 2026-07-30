# Claude Code protocol reference

This guide collects Claude Code-specific hook behavior that cannot be
represented fully in wat's portable API. Use [SDK public API](../sdk.md) for
handler and result patterns, [Agent protocols and normalization](../agent-formats.md)
for cross-agent mappings, and package godoc for exact Go fields and method
signatures.

Official host references:

- [Hooks reference](https://code.claude.com/docs/en/hooks.md)
- [Hooks guide](https://code.claude.com/docs/en/hooks-guide)

## Event capabilities

The table lists the supported `claude.UseHooks()` surface. “Observe” handlers
receive `(context.Context, Event) error`; “result” handlers also receive a
hook-scoped results builder and return an output.

| Event | Handler | Host-consumed output and important constraints |
|---|---|---|
| `SessionStart` | Result | `additionalContext`, optional env via `CLAUDE_ENV_FILE`, and related session fields; cannot block the session |
| `SessionEnd` | Observe | Response body is unused |
| `UserPromptSubmit` | Result | `Block` writes top-level `decision: "block"` with exit 0; context and common fields also supported |
| `UserPromptExpansion` | Result | Context injection today; host also documents `decision: "block"` (not exposed by the builder yet) |
| `PreToolUse` | Result | Allow / deny / ask via `hookSpecificOutput.permissionDecision`; optional `updatedInput` and `additionalContext` |
| `PostToolUse` | Result | Context and optional `decision: "block"` |
| `PostToolUseFailure` | Result | Recovery context (and optional block) |
| `PermissionRequest` | Result | `behavior` allow / deny on behalf of the user; optional interrupt |
| `PermissionDenied` | Result | `retry: true` only; exit codes are ignored because denial already happened |
| `SubagentStart` | Result | Context injection; exit 2 does not cancel the subagent |
| `SubagentStop` | Result | `FollowUp` / `Context`; honors `stop_hook_active` |
| `TaskCreated` / `TaskCompleted` | Result | Context injection; host also documents exit-2 / `continue: false` teammate control (not exposed by the builders yet) |
| `Stop` | Result | `FollowUp` encodes `decision: "block"`; `Context` is non-blocking feedback |
| `Notification` | Result | Context / common fields |
| `MessageDisplay` | Result | Display-content override only; cannot block the message |
| `PreCompact` | Result | Context injection today; host also documents `decision: "block"` (not exposed by the builder yet) |
| `WorktreeCreate` | Result | Must return the created worktree path; see [WorktreeCreate path return](#worktreecreate-path-return) |
| `Elicitation` | Result | Accept / decline / cancel |

Claude also exports typed decode models for additional native events that do
not yet have `UseHooks` registration methods. Treat the fluent method list as
the supported author-facing registration surface. Exact exported event fields
and builders remain documented in `sdk/claude` godoc.

## Permission, Block, and exit policy

Claude blocking is usually JSON on exit 0. Do not assume exit 2 is how wat
builders deny or block.

| Concern | Behavior |
|---|---|
| Successful JSON decisions | Builders encode allow / deny / ask / block fields and exit **0** |
| Handler / mux errors | Exit **1** — Claude treats this as a non-blocking error for most events (fail-open) |
| Host exit 2 | Claude ignores stdout JSON and feeds stderr to the model as a blocking error for events that support it (`PreToolUse`, `UserPromptSubmit`, `PermissionRequest`, and others). wat result builders do **not** use exit 2 for structured decisions |
| `WorktreeCreate` exception | **Any** non-zero exit aborts worktree creation |

| Event | Schema / builder | Notes |
|---|---|---|
| `PreToolUse` | allow / deny / ask | Prefer `Deny` / `Ask` builders over exit 2. Empty / nil stdout leaves the normal permission flow |
| `PermissionRequest` | allow / deny only | No ask value. `Deny` short-circuits the permission dialog; empty output lets Claude prompt the user |
| `UserPromptSubmit` / `PostToolUse` / `Stop` / `SubagentStop` | top-level `decision: "block"` | Encoded by `Block` or `FollowUp` with exit 0 |
| `PermissionDenied` | `retry` | Observational after denial; use `Retry()` to tell the model it may retry |

HTTP hooks (not what `wat install` writes) cannot signal a block through HTTP
status alone; they must return 2xx with JSON decision fields. wat installs
command hooks only.

## WorktreeCreate path return

`WorktreeCreate` replaces Claude Code's default `git worktree` creation. The
hook must return the absolute path of the created worktree directory. Claude
Code does **not** use the usual allow/block JSON decision model for this event.

| Transport | How the path is returned |
|---|---|
| Command hooks (`type: "command"`) — what `wat install` writes | Print the path as **plain stdout text** (last non-empty line). No JSON wrapper |
| HTTP hooks (`type: "http"`) | JSON body with `hookSpecificOutput.hookEventName` = `"WorktreeCreate"` and `hookSpecificOutput.worktreePath` |

`WorktreeCreateResults.Path` currently encodes the HTTP-style JSON
`worktreePath` wrapper. That is incorrect for command hooks: the host treats
the last non-empty stdout line as the directory path, so a JSON object is not a
usable worktree. Redirect diagnostic output to stderr so it does not become the
path line. Any non-zero process exit aborts creation.

## Portable projection boundaries

Portable registration expands onto Claude as described in the
[cross-agent expansion matrix](../agent-formats.md#portable-registration-expansion).
The following Claude-specific restrictions are particularly important:

- Native-only events (`PermissionRequest`, `PermissionDenied`,
  `UserPromptExpansion`, `Notification`, `MessageDisplay`, `WorktreeCreate`,
  `Elicitation`, task events, and decode-only types) have no portable
  registration.
- Portable `OnPreCompact` is observe-only. Claude's native `PreCompact` builder
  can inject context, but host `decision: "block"` is not projected through the
  portable API.
- Portable `OnStop` / `OnSubagentStop` map `FollowUp` onto Claude's
  `decision: "block"` encoding and project `stop_hook_active` when present.
- Session environment variables written through `CLAUDE_ENV_FILE` are Claude-only;
  portable session start exposes context, not env-file writes.

Observe-only portable events never emit host JSON.

## Stop follow-up loops

`StopResults.FollowUp` encodes top-level `decision: "block"` with a `reason`
and exit 0. Claude forces another turn with that reason.

- Input `stop_hook_active` is true when a prior stop-hook block already forced
  continuation for this turn. It is exposed as `Stop.StopHookActive` /
  `SubagentStop.StopHookActive`.
- Gate `FollowUp` on `StopHookActive` to avoid runaway continuation loops.
- `StopResults.Context` injects non-blocking `additionalContext` and does not
  prevent completion.

Portable handlers read the same signal as `StopEvent.Turn.StopHookActive`.

## Payload details

### Common envelope and process model

- Dialects are detected by `session_id` without Cursor/Copilot-specific signals.
- Event names are PascalCase (`PreToolUse`, `UserPromptSubmit`, …).
- Shared output conventions include `continue`, `stopReason`, `systemMessage`,
  `suppressOutput`, and `terminalSequence` on events that apply `Common`.
- Session environment for supported outputs uses `CLAUDE_ENV_FILE`, not stdout
  JSON env maps.

### Tool and permission events

- `PreToolUse` uses `hookSpecificOutput.permissionDecision` values `allow`,
  `deny`, and `ask` (legacy top-level `decision` values are deprecated for this
  event on the host).
- `PermissionRequest` runs only when Claude is about to prompt (or would
  auto-deny without a prompt). Host docs also describe
  `permission_suggestions` input and `updatedPermissions` / `updatedInput`
  allow-side fields; prefer godoc for what the current SDK decodes and emits.
- MCP tools appear under the usual tool events with `mcp__...` names.

### Lifecycle and workspace events

- `Stop` / `SubagentStop` include `last_assistant_message` and
  `stop_hook_active`. Newer Claude Code builds may also send background-task
  and session-cron metadata; check `sdk/claude` godoc for current field
  coverage.
- `WorktreeCreate` is matcher-less and always fires when worktrees are created
  via `--worktree`, `isolation: "worktree"`, or background sessions.
