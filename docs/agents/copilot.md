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
| `PreToolUse` | Result | allow / deny / ask; optional `modified_args`; cloud treats ask as deny |
| `PostToolUse` | Result | Optional `additional_context` |
| `PostToolUseFailure` | Result | Recovery context; command hooks may use exit 2 to carry context |
| `PermissionRequest` | Result | `behavior` allow / deny only; **CLI only** — cloud pre-approves tools |
| `SubagentStart` | Result | Optional context; built-in `general-purpose` agent does not emit this event |
| `SubagentStop` | Result | `FollowUp` / noop; see stop loops |
| `Stop` (`AgentStop`) | Result | `FollowUp` encodes `decision: "block"`; may be main-agent or subagent-scoped |
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
- Native-only events (`PermissionRequest`, `Notification`, `ErrorOccurred`) have
  no portable registration.
- Portable `OnUserPrompt` maps to observe-only `UserPromptSubmit` on Copilot.
- Portable `OnPreCompact` is observe-only, matching Copilot's native observe
  handler.
- Cloud-agent handling may downgrade portable `Ask` to a denial.

Observe-only portable events never emit host JSON.

## Stop follow-up loops

`StopResults.FollowUp` encodes `decision: "block"` with a `reason` string and
exit 0. Copilot forces another agent turn with that reason.

- Input `stop_hook_active` is true when a prior stop-hook block already forced
  continuation for this turn. Gate `FollowUp` on that signal when the payload
  includes it (Claude exposes it today as `StopHookActive`; check `sdk/copilot`
  godoc for current Copilot field coverage).
- Copilot also overrides the hook after several consecutive block continuations
  (host runaway guard).
- When `AgentStop` carries agent identity fields, treat it as subagent-scoped
  (`AgentStop.IsSubagent()`).

Portable handlers read Claude's mapped flag as `StopEvent.Turn.StopHookActive`
when the adapter supplies it.

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
