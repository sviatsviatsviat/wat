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
| `UserPromptExpansion` | Result | Context injection and `Block` (`decision: "block"`, exit 0) |
| `PreToolUse` | Result | Allow / deny / ask via `hookSpecificOutput.permissionDecision`; optional `updatedInput` and `additionalContext` |
| `PostToolUse` | Result | Context and optional `decision: "block"` |
| `PostToolUseFailure` | Result | Recovery context and `Block` (`decision: "block"`, exit 0) |
| `PermissionRequest` | Result | `behavior` allow / deny on behalf of the user; optional interrupt |
| `PermissionDenied` | Result | `retry: true` only; exit codes are ignored because denial already happened |
| `SubagentStart` | Result | Context injection; exit 2 does not cancel the subagent |
| `SubagentStop` | Result | `FollowUp` / `Context`; honors `stop_hook_active` |
| `TaskCreated` / `TaskCompleted` / `TeammateIdle` | Result | Context injection; `Block` uses exit 2 + stderr (prefer over `WithContinue(false)`, which stops the teammate) |
| `Stop` | Result | `FollowUp` encodes `decision: "block"`; `Context` is non-blocking feedback |
| `Notification` | Result | Context / common fields |
| `MessageDisplay` | Result | Display-content override only; cannot block the message |
| `PreCompact` | Result | Context injection and `Block` (`decision: "block"`, exit 0) |
| `PostToolBatch` / `ConfigChange` | Result | Context injection and `Block` (`decision: "block"`, exit 0) |
| `WorktreeCreate` | Result | Must return the created worktree path; see [WorktreeCreate path return](#worktreecreate-path-return) |
| `Elicitation` / `ElicitationResult` | Result | Accept / decline / cancel |
| `Setup` / `CwdChanged` / `FileChanged` | Result | Context / common fields |
| `StopFailure` / `PostCompact` / `WorktreeRemove` / `InstructionsLoaded` | Observe | Response body unused |

Treat the fluent `UseHooks` method list as the supported author-facing
registration surface. Exact exported event fields and builders remain
documented in `sdk/claude` godoc.

## Permission, Block, and exit policy

Claude blocking is usually JSON on exit 0. Do not assume exit 2 is how wat
builders deny or block.

| Concern | Behavior |
|---|---|
| Successful JSON decisions | Builders encode allow / deny / ask / block fields and exit **0** |
| Handler / mux errors | Exit **1** — Claude treats this as a non-blocking error for most events (fail-open) |
| `TeammateIdle` / `TaskCreated` / `TaskCompleted` `Block` | Exit **2** (`BlockExit`) — stdout JSON ignored; `run.Serve` writes the reason to stderr for the model |
| Other host exit 2 | Claude may treat exit 2 as a blocking error for some events, but wat JSON deny/block builders do **not** use exit 2 |
| `WorktreeCreate` exception | **Any** non-zero exit aborts worktree creation |

| Event | Schema / builder | Notes |
|---|---|---|
| `PreToolUse` | allow / deny / ask / defer | Prefer `Deny` / `Ask` / `Defer` builders over exit 2. Empty / nil stdout leaves the normal permission flow |
| `PermissionRequest` | allow / deny only | No ask value. `Deny` short-circuits the permission dialog; empty output lets Claude prompt the user |
| `UserPromptSubmit` / `UserPromptExpansion` / `PostToolUse` / `PostToolUseFailure` / `PostToolBatch` / `ConfigChange` / `PreCompact` / `Stop` / `SubagentStop` | top-level `decision: "block"` | Encoded by `Block` or `FollowUp` with exit 0 |
| `TeammateIdle` / `TaskCreated` / `TaskCompleted` | exit 2 + stderr via `Block`; `WithContinue(false)` stops the teammate | Prefer `Block` when rolling back idle/create/complete without stopping the teammate |
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

wat installs command hooks only. `WorktreeCreateResults.Path` therefore encodes
the path as plain stdout bytes. Emitting JSON `worktreePath` on a command hook
is incorrect: the host treats the last non-empty stdout line as the directory
path, so a JSON object is not a usable worktree.

Shared `WithContinue` / `WithSystemMessage` and related Common fields remain on
the Go output type for merge consistency, but they are not written to stdout
for this event. Redirect diagnostic output to stderr so it does not become the
path line. Any non-zero process exit aborts creation.

## Portable projection boundaries

Portable registration expands onto Claude as described in the
[cross-agent expansion matrix](../agent-formats.md#portable-registration-expansion).
The following Claude-specific restrictions are particularly important:

- Native-only events (`PermissionRequest`, `PermissionDenied`,
  `UserPromptExpansion`, `Notification`, `MessageDisplay`, `WorktreeCreate`,
  `Elicitation`, task events, and decode-only types) have no portable
  registration.
- Portable `OnPreCompact` is observe-only even though Claude's native
  `PreCompact` can `Block`. Use `sdk/claude` when compaction must be denied.
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
  auto-deny without a prompt). The SDK decodes optional `permission_suggestions`
  and can emit `updatedPermissions` / `updatedInput` on allow via
  `WithUpdatedPermissions` / related builders.
- MCP tools appear under the usual tool events with `mcp__...` names.

### Lifecycle and workspace events

- `Stop` / `SubagentStop` include `last_assistant_message` and
  `stop_hook_active`. Claude Code v2.1.145+ may also send `background_tasks`
  and `session_crons`; see [Stop background metadata](#stop-and-subagentstop-background-metadata).
- `WorktreeCreate` is matcher-less and always fires when worktrees are created
  via `--worktree`, `isolation: "worktree"`, or background sessions.

## Stop and SubagentStop background metadata

Claude Code v2.1.145+ includes `background_tasks` and `session_crons` on `Stop`
and `SubagentStop` inputs so handlers can tell a finished turn from a pause
that will wake when background work or a session cron fires.

| Field | Meaning |
|---|---|
| `background_tasks` | In-flight tasks (`shell`, `subagent`, `monitor`, `workflow`, `teammate`, `cloud session`, `MCP task`, or a raw discriminant). Type-specific extras include `command`, `agent_type`, `server`/`tool`, and `name`. |
| `session_crons` | Session-scoped scheduled wakeups from `CronCreate`, `ScheduleWakeup`, and `/loop`, with `schedule`, `recurring`, and `prompt`. |

Both arrays are present when the task registry is reachable and empty when
nothing is in flight or scheduled. On `SubagentStop` they remain scoped to the
parent session, not the subagent. Older Claude Code builds omit the fields;
wat decodes that as nil slices.

`sdk/claude` exposes them on `Stop` / `SubagentStop` as `BackgroundTasks` and
`SessionCrons` (`BackgroundTask` / `SessionCron`). Portable `OnStop` /
`OnSubagentStop` do not project these Claude-only inputs. Use native Claude
handlers when FollowUp / continue logic must branch on background or cron
state, and still honor `stop_hook_active` to avoid runaway continuation loops.
