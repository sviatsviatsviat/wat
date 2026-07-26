# Agent protocols and normalization

This document records the cross-agent protocol facts that affect codecs,
portable mappings, fixtures, and public behavior. It is not a replacement for
the native agent documentation or the typed SDK godoc.

## Dialect detection and event identity

All supported payloads must include `hook_event_name`.

| Dialect | Detection signals used by wat |
|---|---|
| Claude Code | `session_id`, excluding Cursor and Copilot-specific signals |
| GitHub Copilot | `hook_event_name` plus `timestamp`, excluding Cursor signals |
| Cursor | `cursor_version`, `conversation_id`, or `CURSOR_VERSION` |

Each native SDK exports its string `Dialect` constant. Portable event
`Envelope.Agent` uses that value, while `Envelope.Name` preserves the native
event name.

The installed `wat run --agent ... --event ...` flags identify managed config
entries. Runtime dispatch still trusts the payload detected by `sdk/run`.

## Portable registration expansion

One portable registration may install more than one native event, especially
for Cursor's dedicated tool hooks.

| Portable method | Claude | Copilot | Cursor |
|---|---|---|---|
| `OnSessionStart` | `SessionStart` | `SessionStart` | `sessionStart` |
| `OnSessionEnd` | `SessionEnd` | `SessionEnd` | `sessionEnd` |
| `OnUserPrompt` | `UserPromptSubmit` | `UserPromptSubmitted` | `beforeSubmitPrompt` |
| `OnPreTool` | `PreToolUse` | `PreToolUse` | `preToolUse`, `beforeShellExecution`, `beforeMCPExecution`, `beforeReadFile` |
| `OnPostTool` | `PostToolUse` | `PostToolUse` | `postToolUse`, `afterMCPExecution`, `afterFileEdit` |
| `OnPostToolFailure` | `PostToolUseFailure` | `PostToolUseFailure` | `postToolUseFailure` |
| `OnSubagentStart` | `SubagentStart` | `SubagentStart` | `subagentStart` |
| `OnSubagentStop` | `SubagentStop` | `SubagentStop` and subagent-scoped `AgentStop` | `subagentStop` |
| `OnStop` | `Stop` | main-agent `AgentStop` | `stop` |
| `OnPreCompact` | `PreCompact` | `PreCompact` | `preCompact` |

Copilot's `AgentStop` wire event can describe either the main agent or a
subagent. The adapter routes it by its optional agent identity fields so
`OnStop` and `OnSubagentStop` do not both handle the same payload.

Events not representable on every agent are native-only. Examples include
permission requests and notifications (Claude/Copilot), Claude worktree and
elicitation events, and Cursor workspace/tab events.

## Tool-name normalization

Portable `ToolCall.Name` uses a canonical vocabulary. `ToolCall.Native` always
preserves the original value.

| Native examples | Canonical name |
|---|---|
| `Bash`, `bash`, `Shell`, `powershell` | `bash` |
| `Edit`, `edit`, `notebookedit` | `edit` |
| `Write`, `create` | `write` |
| `Read`, `view` | `read` |
| `Agent`, `task` | `task` |
| `Glob` | `glob` |
| `Grep` | `grep` |
| `WebFetch`, `web_fetch` | `web_fetch` |
| `WebSearch`, `web_search` | `web_search` |

Unknown names pass through unchanged. MCP detection recognizes Claude/Copilot
`mcp__...` names and Cursor `MCP:...` names; structured dedicated MCP events
also set `ToolCall.MCP`.

Dedicated Cursor events synthesize a portable tool identity:

| Cursor event | Portable tool |
|---|---|
| `beforeShellExecution` | `bash` with `ToolCall.Shell` |
| `afterShellExecution` (native observe-only) | not projected onto portable `OnPostTool` |
| `beforeReadFile` | `read` |
| `afterFileEdit` | `edit` |
| `beforeMCPExecution` / `afterMCPExecution` | Native MCP tool name with `MCP=true` |

Native names for these synthesized calls may be the event name rather than a
builtin tool name. Code that needs exact host identity should use `Native` and
`Envelope.Name`.

## Portable result projection

Portable builders deliberately expose the intersection of native behavior:

| Portable result | Native intent |
|---|---|
| `PreToolResults.Allow` | Explicitly allow the call |
| `PreToolResults.Deny(reason)` | Deny/block with agent-facing reason |
| `PreToolResults.Ask(reason)` | Escalate where supported |
| `PreToolResult.WithUpdatedInput` | Replace arguments where that native event supports it |
| `PostToolResults.Context` | Add model-facing context |
| `PostToolResult.WithUpdatedOutput` | Replace supported tool output |
| `PostToolFailureResults.Context` | Add recovery context |
| `StopResults.FollowUp` | Prevent completion and request more work |
| `SessionStartResults.Context` | Add startup context |

Known limitations are part of the contract:

- Copilot cloud-agent handling may downgrade `Ask` to a denial.
- Cursor emits `updated_input` only for generic `preToolUse`, not for its
  dedicated pre-tool events.
- Cursor does not support `"ask"` on `subagentStart`; it is treated as
  `"deny"`. `sdk/cursor`'s `SubagentStart` `Deny` writes `user_message` and
  exits 0 so Cursor applies the JSON permission field (exit 2 would re-wrap
  stdout as the message). Prefer `Deny` over `Ask` for this event.
- Cursor `WithUpdatedOutput` maps to `updated_mcp_tool_output` on generic
  `postToolUse` for MCP tools only.
- Cursor observe-only post-tool events (Hooks docs list no consumed output
  fields):
  - `afterFileEdit`: `OnPostTool` still expands for edit observation, but
    portable `Context` / `WithUpdatedOutput` have no host effect. Native
    `sdk/cursor` registration is side-effects only (for example formatters).
  - `afterShellExecution`: not projected onto portable `OnPostTool`; use
    `sdk/cursor.AfterShellExecution` for auditing. Shell post-tool context
    remains via `postToolUse`. Decodes `sandbox`.
  - `afterMCPExecution`: `OnPostTool` expands for observation, but portable
    builders are discarded. Rewrite MCP tool output via `postToolUse`
    (`updated_mcp_tool_output`). Cloud agents do not load
    `beforeMCPExecution` / `afterMCPExecution`.
  - `postToolUseFailure`: native observe handler; portable
    `PostToolFailureResults.Context` is discarded on Cursor. Decodes
    `is_interrupt`.
- Observe-only portable events never emit host JSON.

Do not widen the portable interface until every dialect has a truthful mapping
and tests for it.

## Wire and process differences

| Concern | Claude Code | GitHub Copilot | Cursor |
|---|---|---|---|
| Event-name style | PascalCase | PascalCase | camelCase |
| Common JSON fields | Mostly snake_case with native output conventions | snake_case | snake_case/camelCase per native schema |
| Blocking | Usually encoded in JSON output fields | JSON decision plus native exit behavior | Permission denial may use exit 2 |
| Handler error | Exit 1 | Exit 1 | Exit 1 (host normally treats as fail-open) |
| Session environment | `CLAUDE_ENV_FILE` for supported output | stdout JSON | stdout JSON |
| Project config | `.claude/settings.json` matcher groups | `.github/hooks/wat.json` flat handlers | `.cursor/hooks.json` flat handlers |

Cursor's documented common input schema includes optional `model_id` and
`model_params` alongside the legacy `model` slug. `sdk/cursor.Envelope` decodes
those fields (with exported `ModelParam`) so handlers do not silently lose them.

Output encoding belongs to the native SDK. Portable adapters must return native
output values and must not serialize JSON themselves.

## Fixture and codec expectations

Codec changes require:

- a representative native JSON fixture or focused payload in tests;
- event-name peeking coverage;
- typed decode assertions;
- output JSON and exit-code assertions for result-producing events;
- a missing/unknown event test;
- portable mapping tests when the event participates in `sdk/agnostic`.

Unknown event names are decode errors when decoded directly. During normal
`run.Serve`, an event with no registered handler exits successfully before full
decode.

## Source locations

| Concern | Source |
|---|---|
| Shared codec/normalization machinery | [`internal/hookkit`](../internal/hookkit/) |
| Native event slices | [`sdk/claude/internal/hooks`](../sdk/claude/internal/hooks/), [`sdk/copilot/internal/hooks`](../sdk/copilot/internal/hooks/), [`sdk/cursor/internal/hooks`](../sdk/cursor/internal/hooks/) |
| Native public facades and registrars | [`sdk/claude`](../sdk/claude/), [`sdk/copilot`](../sdk/copilot/), [`sdk/cursor`](../sdk/cursor/) |
| Portable model and adapters | [`sdk/agnostic/internal`](../sdk/agnostic/internal/) |
| Portable typed tool inputs | [`sdk/agnostic/tools`](../sdk/agnostic/tools/) |
| Installation schemas | [`cmd/wat/internal/hostconfig`](../cmd/wat/internal/hostconfig/) |
| End-to-end fixtures | [`e2e/testdata`](../e2e/testdata/) |
| Project `wat test` fixtures | [`.wat/testdata`](../.wat/testdata/) |
