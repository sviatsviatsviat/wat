# Agent protocols and normalization

This document records the cross-agent protocol facts that affect codecs,
portable mappings, fixtures, and public behavior. It is not a replacement for
the native agent documentation or the typed SDK godoc.

## Dialect detection and event identity

When `--event` is absent, supported payloads must include `hook_event_name`.
When `wat run` (or the hooks binary) is invoked with `--event`, Serve selects
that event without peeking; missing `hook_event_name` is allowed.

| Dialect | Detection signals used by wat |
|---|---|
| Claude Code | `session_id`, excluding Cursor and Copilot-specific signals |
| GitHub Copilot | `hook_event_name` plus `timestamp`, excluding Cursor signals |
| Cursor | `cursor_version`, `conversation_id`, or `CURSOR_VERSION` |

Each native SDK exports its string `Dialect` constant. Portable event
`Envelope.Agent` uses that value, while `Envelope.Name` preserves the native
event name.

Installed `wat run --agent ... --event ...` flags identify managed config
entries and are forwarded to the hooks binary as dispatch hints: `--agent` /
`--event` select dialect and event without using payload detection or
`hook_event_name` peek for that choice. Serve may still inspect the payload
only to warn when a hint disagrees; mismatches do not fail the run.
`wat doctor` warns when a command’s flags disagree with the native config map
key or agent file.

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

Copilot `AgentStop` (including subagent-scoped Stop) also decodes
`stop_hook_active`. When that flag is true, a prior stop-hook `decision:
"block"` already forced continuation for the turn; authors must gate
`FollowUp` on it to avoid runaway loops. Explicit `SubagentStop` wire events
do not document this field today. Copilot additionally caps consecutive block
continuations after several iterations. See the

Keep the Copilot wire split intact:

- Explicit `hook_event_name: "SubagentStop"` decodes as `sdk/copilot.SubagentStop`
  and carries `last_assistant_message` (full final subagent response text; the
  VS Code name for camelCase `response`). Native handlers use
  `SubagentStopResults` (`FollowUp` / `ModifiedResponse`).
- VS Code-style `hook_event_name: "Stop"` payloads with `agent_name` /
  `agent_display_name` still decode as `AgentStop` (`IsSubagent`); they do not
  become `SubagentStop` and do not expose `modifiedResponse`. Use
  `StopResults.FollowUp` there.

`modifiedResponse` is Copilot-native on `SubagentStop` only: a `block`
decision wins over a rewrite, and rewrites do not compose across handlers
(last non-empty rewrite wins; each handler still sees the original response).
Portable `OnSubagentStop` maps the message onto `Turn.LastAssistantMessage` /
`Subagent.Summary` and still exposes only `FollowUp`.

Events not representable on every agent are native-only. Examples include
permission requests and notifications (Claude/Copilot), Copilot
`UserPromptTransformed` (post-transform prompt rewrite), Claude worktree and
elicitation events, and Cursor workspace/tab events. See the
[Cursor protocol reference](agents/cursor.md) for its native-only lifecycle,
permission, and Tab behavior.

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

On Claude and Copilot stop events, check `Turn.StopHookActive` (native
`stop_hook_active`) before emitting another `FollowUp`. Cursor uses
`Turn.LoopCount` / hooks.json `loop_limit` instead; see
[Cursor stop follow-up loops](agents/cursor.md#stop-follow-up-loops) and

Known limitations are part of the contract:

- Copilot cloud-agent handling treats `PreToolUse` `"ask"` as a denial (no
  user). Native Copilot `PermissionRequest.Ask` is a soft deny
  (`behavior:"deny"`, exit 0), not a confirmation UI; use `Noop` to fall
  through to the host permission prompt, and prefer `Deny` to block. Cloud
  agents do not apply `permissionRequest` (gate with `PreToolUse`). Command-hook
  timeouts are fail-open.
- Cursor emits `updated_input` only for generic `preToolUse`, not for its
  dedicated pre-tool events.
- Cursor `WithUpdatedOutput` maps to `updated_mcp_tool_output` on generic
  `postToolUse` for MCP tools only.
- Cursor observe-only post-tool events discard portable results, and
  `afterShellExecution` is not projected onto portable `OnPostTool`.
- Observe-only portable events never emit host JSON.
- Portable `OnPreCompact` is observe-only and maps only shared compaction
  fields. Cursor-only metrics and native output remain in `sdk/cursor`.

For the complete event-by-event Cursor behavior—including observe-only
handlers, permission encodings, matcher values, cloud availability, follow-up
loops, and live payload differences—use the
[Cursor protocol reference](agents/cursor.md).

Do not widen the portable interface until every dialect has a truthful mapping
and tests for it.

## Wire and process differences

| Concern | Claude Code | GitHub Copilot | Cursor |
|---|---|---|---|
| Event-name style | PascalCase | PascalCase | camelCase |
| Common JSON fields | Mostly snake_case with native output conventions | snake_case | snake_case/camelCase per native schema |
| Blocking | Usually encoded in JSON output fields; `WorktreeCreate` uses plain-path stdout and any non-zero exit | JSON decision plus native exit behavior | Permission denial may use exit 2; `beforeSubmitPrompt` uses `continue` JSON with exit 0 |
| Handler error | Exit 1 | Exit 1 | Exit 1 (host normally treats as fail-open) |
| Session environment | `CLAUDE_ENV_FILE` for supported output | stdout JSON | stdout JSON |
| Project config | `.claude/settings.json` matcher groups | `.github/hooks/wat.json` flat handlers | `.cursor/hooks.json` flat handlers |

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
