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
elicitation events, and Cursor workspace/tab events. Cursor `workspaceOpen` is
an app lifecycle hook (desktop/CLI): it can return `pluginPaths`, is skipped
when there are zero workspace folders, and does not run in cloud agents.

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

Cursor Tab edit hooks are native-only (`beforeTabFileRead`,
`afterTabFileEdit`). On `afterTabFileEdit`, each edit includes `old_string`,
`new_string`, `range` (`start_line_number`, `start_column`, `end_line_number`,
`end_column`), `old_line`, and `new_line`. Agent `afterFileEdit` keeps the
simpler `old_string` / `new_string` shape; `sdk/cursor` exposes that as `Edit`
and the Tab shape as `TabEdit`.

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

Cursor `stop` / `subagentStop` follow-up loops:

- `FollowUp` encodes non-empty `followup_message` with exit 0; Cursor
  auto-submits that text as the next user message.
- For `subagentStop`, Cursor only consumes `followup_message` when the input
  `status` is `"completed"`. For `stop`, a non-empty message is always
  eligible.
- Input `loop_count` is how many automatic follow-ups the same script has
  already triggered for the conversation (starts at 0).
- Cursor enforces a per-script `loop_limit` from `hooks.json` (default `5`;
  `null` means unlimited). Authors should check `loop_count` before emitting
  another follow-up. That option is install config, not an SDK field.

Known limitations are part of the contract:

- Copilot cloud-agent handling may downgrade `Ask` to a denial.
- Cursor emits `updated_input` only for generic `preToolUse`, not for its
  dedicated pre-tool events.
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
- Cursor `subagentStop` decodes the documented telemetry fields
  (`description`, `duration_ms`, `message_count`, `tool_call_count`,
  `modified_files`) in addition to identity and status fields.
- Cursor `afterAgentThought` is observe-only and decodes optional
  `duration_ms` for the completed thinking block. Cursor's hooks.json
  matcher for this event is the fixed value `AgentThought`.
- Cursor `sessionStart` is fire-and-forget: the agent loop does not wait for or
  enforce a blocking response. The schema accepts `continue` / `user_message`,
  but callers do not enforce them; session creation is not blocked when
  `continue` is `false`. Meaningful outputs are `env` and `additional_context`.
  Input may include optional `composer_mode` (`"agent"`, `"ask"`, or `"edit"`).
  The hook is not available for Cursor cloud agents.
- Cursor `sessionEnd` is tied to the IDE composer session and is not available
  for cloud agents. The response body is unused (observe-only).
- Observe-only portable events never emit host JSON.
- Portable `OnPreCompact` is observe-only and maps only shared compaction
  fields (`trigger`, plus Claude/Copilot `custom_instructions` when present).
  Cursor's native `preCompact` also carries compaction metrics
  (`context_usage_percent`, `context_tokens`, `context_window_size`,
  `message_count`, `messages_to_compact`, `is_first_compaction`) and may emit
  an observational `user_message` via `sdk/cursor`'s `PreCompact` /
  `UserMessage`. Those Cursor-only inputs and the `user_message` output stay
  on the native SDK; they are not part of the portable contract.

### Cursor permission and Ask semantics

Cursor permission-gating is event-specific. Use this matrix; do not assume a
single Ask or Deny encoding across events.

| Event | Schema permissions | Ask host behavior | Deny encoding (`sdk/cursor`) | Notes |
|---|---|---|---|---|
| `beforeShellExecution` | allow / deny / ask | **Enforced** (user approval) | `agent_message`, exit 2 by default; `WithUserMessage` for client-facing copy | Same Ask enforcement as `beforeMCPExecution` |
| `beforeMCPExecution` | allow / deny / ask | **Enforced** (user approval) | `agent_message`, exit 2 by default; `WithUserMessage` for client-facing copy | Contrasts with `preToolUse` |
| `preToolUse` | allow / deny / ask | **Not enforced** today | `agent_message`, exit 2 | `PreToolUseResults.Ask` still encodes `"ask"` for schema compatibility; prefer `Allow`/`Deny` to gate. Optional input `agent_message` (pre-call narrative) is decoded |
| `beforeReadFile` | allow / deny (+ optional `user_message`) | **Coerced to deny** (no `"ask"`) | `user_message`, exit 0 (no `agent_message`) | Prefer `Deny` over `Ask` |
| `beforeTabFileRead` | allow / deny only | **N/A** (no ask API) | permission-only JSON, exit 0 (no message fields) | Tab-only; not available in cloud agents. Chained message/`updated_input` helpers are ignored on the wire |
| `subagentStart` | allow / deny | **Treated as deny** (no `"ask"`) | `user_message`, exit 0 (no `agent_message`) | Prefer `Deny` over `Ask`. Exit 2 would re-wrap stdout as the user message |

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
