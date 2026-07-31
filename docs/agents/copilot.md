# GitHub Copilot protocol reference

This guide collects Copilot-specific hook behavior that cannot be represented
fully in wat's portable API. Use [SDK public API](../sdk.md) for handler and
result patterns, [Agent protocols and normalization](../agent-formats.md) for
cross-agent mappings, and package godoc for exact Go fields and method
signatures.

Official host references:

- [Copilot Hooks reference](https://docs.github.com/en/copilot/reference/hooks-reference)
- [About hooks for GitHub Copilot](https://docs.github.com/en/copilot/concepts/agents/hooks)

## Supported protocol surface

`sdk/copilot` implements one wire dialect: the hooks-reference
**PascalCase / snake_case** shape (the reference calls this the “VS Code
compatible” payload format).

| Concern | wat behavior |
|---|---|
| Config event keys | PascalCase in `.github/hooks/wat.json` (for example `PreToolUse`) |
| `wat run --event` | Same PascalCase string as the config key / `Event*` constant value |
| stdin `hook_event_name` | PascalCase (for example `"PreToolUse"`, `"Stop"`) |
| stdin / stdout JSON fields | snake_case (for example `session_id`, `tool_name`, `permission_decision`) |
| stdout shape | Flat JSON object encoded by `sdk/copilot` builders |

`wat install` writes PascalCase keys and
`wat run --agent copilot --event <PascalCase>`. Fixtures under
`.wat/testdata/copilot/` and `e2e/testdata/` use the same shape.

### Explicit non-goals

These surfaces appear in upstream docs or product samples but are **not**
wat’s Copilot dialect today:

| Surface | Why it does not apply |
|---|---|
| Copilot CLI **camelCase** config / payloads (`preToolUse`, `sessionId`, …) | Alternate hooks-reference format. wat does not install camelCase keys and does not decode camelCase stdin as Copilot |
| [VS Code agent-customization hooks](https://code.visualstudio.com/docs/agent-customization/hooks) **`hookSpecificOutput`** wrapper | Divergent preview wire (nested `hookSpecificOutput`). Not a second `sdk/copilot` dialect |
| Ad-hoc casing aliases (`preToolUse` ↔ `PreToolUse` decode) | Deferred optional enhancement; authors must use PascalCase / snake_case |

Do not assume that a sample that “works in Copilot CLI docs” or “works in VS
Code hooks docs” will decode or install under wat without matching the table
above.

The hooks reference labels PascalCase + snake_case “VS Code compatible.” That
name means the Copilot extension-style **flat** stdin fields—not the separate
VS Code agent-customization output envelope built around `hookSpecificOutput`.

### Event-name matrix

Upstream CLI samples often use camelCase event ids. wat config keys,
`--event` hints, and `hook_event_name` use the PascalCase column.

| Copilot CLI camelCase (docs samples) | wat / hooks-reference PascalCase | `sdk/copilot` notes |
|---|---|---|
| `sessionStart` | `SessionStart` | `EventSessionStart` |
| `sessionEnd` | `SessionEnd` | `EventSessionEnd` |
| `userPromptSubmitted` | `UserPromptSubmit` | Go API `UserPromptSubmitted` / `EventUserPromptSubmitted` → wire `"UserPromptSubmit"` |
| `userPromptTransformed` | `UserPromptTransformed` | `EventUserPromptTransformed` |
| `preToolUse` | `PreToolUse` | `EventPreToolUse` |
| `postToolUse` | `PostToolUse` | `EventPostToolUse` |
| `postToolUseFailure` | `PostToolUseFailure` | `EventPostToolUseFailure` |
| `permissionRequest` | `PermissionRequest` | `EventPermissionRequest` |
| `subagentStart` | `SubagentStart` | `EventSubagentStart` |
| `subagentStop` | `SubagentStop` | `EventSubagentStop` |
| `agentStop` | `Stop` | Go API `AgentStop` / `EventAgentStop` → wire `"Stop"` |
| `preCompact` | `PreCompact` | `EventPreCompact` |
| `notification` | `Notification` | `EventNotification` |
| `errorOccurred` | `ErrorOccurred` | `EventErrorOccurred` |

Go method and type names sometimes differ from the wire string (notably
`UserPromptSubmitted` and `AgentStop`). Always use the `Event*` constant value
for install keys, `--event`, and fixtures.

## Event capabilities

The table lists the supported `copilot.UseHooks()` surface. “Observe” handlers
receive `(context.Context, Event) error`; “result” handlers also receive a
hook-scoped results builder and return an output.

| Event | Handler | Host-consumed output and important constraints |
|---|---|---|
| `SessionStart` | Result | Optional `additional_context`; fires once per cloud job as a new session |
| `SessionEnd` | Observe | Response body is unused |
| `UserPromptSubmit` | Observe | Go API `UserPromptSubmitted`; observe-only today |
| `UserPromptTransformed` | Result | Optional `modified_transformed_prompt` rewrite only (cannot block) |
| `PreToolUse` | Result | allow / deny / ask; optional `modified_args`; cloud treats ask as deny |
| `PostToolUse` | Result | Optional `additional_context` and success `modified_result` via `WithModifiedResult` (`result_type: "success"`). Failure-shaped `modified_result` is not supported until the host schema is verified — use `PostToolUseFailure` after a real tool failure |
| `PostToolUseFailure` | Result | Recovery context; command hooks may use exit 2 to carry context |
| `PermissionRequest` | Result | `behavior` allow / deny only; **CLI only** — cloud pre-approves tools |
| `SubagentStart` | Result | Optional context; built-in `general-purpose` agent does not emit this event |
| `SubagentStop` | Result | `FollowUp` and/or `ModifiedResponse` (`modifiedResponse`); decodes `last_assistant_message`; see stop loops |
| `Stop` (`AgentStop`) | Result | `FollowUp` encodes `decision: "block"`; may be main-agent or subagent-scoped; decodes `stop_hook_active` |
| `PreCompact` | Observe | Observe-only today |
| `Notification` | Result | Optional context; **CLI only** (does not fire under cloud agent) |
| `ErrorOccurred` | Observe | Response body is unused |

Cloud availability and timeout policy are host constraints, not fields added by
wat. Exact exported builders remain documented in `sdk/copilot` godoc.

## Permission and Ask semantics

Copilot permission behavior is event-specific. Do not assume that `Ask` opens a
confirmation UI.

| Event | Schema decisions | Ask / soft-deny host behavior | Deny encoding (`sdk/copilot`) | Notes |
|---|---|---|---|---|
| `PermissionRequest` | `behavior` allow \| deny only | **Soft deny** — `Ask` encodes `"behavior":"deny"` with exit 0 (no WarnExit) and does not `Stop` later handlers. It does **not** escalate to the user | `behavior` + `message`, WarnExit (2); `Stop` skips remaining handlers | Prefer `Deny` to block. Prefer `Noop` / nil / empty stdout to fall through to the host permission service (rules, session approvals, **user prompting**). Optional `interrupt` with deny stops the session |
| `PreToolUse` | `permission_decision` allow \| deny \| ask | Encodes `"ask"` on the wire. **Cloud agent treats ask as deny** (no user) | deny + reason; exit 0 for structured JSON | Prefer `Allow` / `Deny` when gating must not depend on a user |

### `PermissionRequest.Ask` footgun

Authors often read `Ask` as “prompt the user.” On Copilot `PermissionRequest`,
the schema has no ask value. The SDK keeps `Ask` as an explicit soft deny for
compatibility with older handlers:

- Wire: `"behavior":"deny"` plus optional `message`
- Process exit: `0` (unlike hard `Deny`, which uses WarnExit `2`)
- Merge / `Stop`: soft deny does not stop remaining handlers; hard `Deny` does

To let Copilot prompt the user (or apply its normal permission rules), return
`Noop()`, `nil`, or empty stdout—not `Ask`.

## Cloud agent and timeouts

| Concern | Host behavior |
|---|---|
| `PermissionRequest` under cloud agent | Tool calls are pre-approved; the hook either does not fire or has no effect. Gate with `PreToolUse` instead |
| `PreToolUse` `"ask"` under cloud agent | Treated as `"deny"` because no user is available |
| `Notification` under cloud agent | Does not fire |
| Command-hook timeouts | **Always fail-open**, including for `PreToolUse` and policy hooks: a timed-out hook warns and lets the tool proceed through the normal permission flow |
| `PreToolUse` non-timeout errors | Command hooks fail-closed (deny), including exit 2 and other non-zero exits |
| Network / sandbox | Cloud outbound network is firewalled; only repo `.github/hooks/*.json` is discovered by default |

## Portable projection boundaries

Portable registration expands onto Copilot as described in the
[cross-agent expansion matrix](../agent-formats.md#portable-registration-expansion).
The following Copilot-specific restrictions are particularly important:

- Copilot's `Stop` wire event can describe either the main agent or a subagent.
  The portable adapter routes by optional agent identity fields so `OnStop` and
  `OnSubagentStop` do not both handle the same payload.
- Native-only events (`PermissionRequest`, `Notification`, `ErrorOccurred`,
  `UserPromptTransformed`) have no portable registration.
- Portable `OnUserPrompt` maps to observe-only `UserPromptSubmit` on Copilot.
- Portable `OnPreCompact` is observe-only, matching Copilot's native observe
  handler.
- Portable `OnSubagentStop` maps Copilot `last_assistant_message` onto
  `Turn.LastAssistantMessage` / `Subagent.Summary` but still exposes only
  `FollowUp` (not `ModifiedResponse`).
- Cloud-agent handling may downgrade portable `Ask` to a denial.

Observe-only portable events never emit host JSON.

## Stop follow-up loops

`AgentStop` uses `StopResults.FollowUp`; `SubagentStop` uses
`SubagentStopResults.FollowUp`. Both encode `decision: "block"` with a `reason`
string and exit 0. Copilot forces another agent turn with that reason.
`SubagentStop` may also rewrite the final message with `ModifiedResponse`
(`modifiedResponse`); a FollowUp/block decision wins over a rewrite.

- Input `stop_hook_active` (`AgentStop.StopHookActive`) is true when a prior
  stop-hook block already forced continuation for this turn. Gate `FollowUp` on
  that signal.
- Copilot also overrides the hook after several consecutive block continuations
  (host runaway guard).
- When `AgentStop` carries agent identity fields, treat it as subagent-scoped
  (`AgentStop.IsSubagent()`).

Portable handlers read the mapped flag as `StopEvent.Turn.StopHookActive` when
the adapter supplies it.

## Field casing checklist

When reading upstream samples, map CLI camelCase fields to wat snake_case:

| camelCase (CLI docs) | snake_case (wat) |
|---|---|
| `sessionId` | `session_id` |
| `toolName` / `toolArgs` | `tool_name` / `tool_input` |
| `permissionDecision` | `permission_decision` |
| `permissionDecisionReason` | `permission_decision_reason` |
| `modifiedArgs` | `modified_args` |
| `additionalContext` | `additional_context` |
| `stopHookActive` | `stop_hook_active` |

Exact builder fields remain in `sdk/copilot` godoc. Prefer result builders over
hand-written stdout JSON.

## Matchers and typed tool inputs

`wat install` writes catch-all command handlers with **no** `matcher` field so
one `wat run` process sees every tool. Optional host `matcher` regexes in
`.github/hooks/*.json` are supported if you edit the file yourself.

PascalCase `PreToolUse` / `PermissionRequest` configs may use Claude-format
matcher and `tool_name` values (`Bash`, `Edit`, `Write`, …). Do not copy Claude
matcher strings onto Copilot configs without checking the host tool-name table.
Filter in Go with `NativeToolName()` / `Input.As*` instead of install-time
matchers when using wat-managed handlers.

| Accessor | Runtime names | Claude-format / aliases |
|---|---|---|
| `AsBash` | `bash`, `powershell`, `shell` | `Bash` (case-insensitive) |
| `AsView` | `view` | `Read` |
| `AsCreate` | `create` | `Write` |
| `AsEdit` | `edit`, `str_replace_editor`, `apply_patch` | `Edit` |
| `AsGlob` | `glob` | `Glob` |
| `AsGrep` | `grep`, `rg` | `Grep` |
| `AsWebFetch` | `web_fetch` | `WebFetch` |
| `AsWebSearch` | `web_search` | `WebSearch` |
| `AsTask` | `task` | `Agent`, `Task` |
| `AsAskUser` | `ask_user` | `AskUserQuestion` |
| `AsUpdateTodo` | `update_todo` | `TodoWrite` |

Other tool names stay available via `Input.Name()` / `Input.Raw()`. Portable
handlers may also use `sdk/agnostic/tools` accessors after name normalization.

