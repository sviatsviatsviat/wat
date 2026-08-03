# Cursor protocol reference

This guide collects Cursor-specific hook behavior that cannot be represented
fully in wat's portable API. Use [SDK public API](../sdk.md) for handler and
result patterns, [Agent protocols and normalization](../agent-formats.md) for
cross-agent mappings, and package godoc for exact Go fields and method
signatures.

## Event capabilities

The table lists the supported `cursor.UseHooks()` surface. “Observe” handlers
receive `(context.Context, Event) error`; “result” handlers also receive a
hook-scoped results builder and return an output.

| Event | Handler | Host-consumed output and important constraints |
|---|---|---|
| `sessionStart` | Result | `env` and `additional_context`; fire-and-forget, so `continue` / `user_message` are not enforced; unavailable to cloud agents |
| `sessionEnd` | Observe | Response body is unused; IDE composer lifecycle only and unavailable to cloud agents |
| `workspaceOpen` | Result | `pluginPaths` with absolute plugin directories; skipped with zero workspace folders and unavailable to cloud agents |
| `beforeSubmitPrompt` | Result | `Block` writes `continue: false` and optional `user_message` with exit 0; hooks.json matcher value is `UserPromptSubmit` |
| `preToolUse` | Result | Allow/deny and schema-compatible ask; Cursor does not enforce ask on this event today |
| `postToolUse` | Result | Generic post-tool response; MCP output replacement uses `updated_mcp_tool_output` |
| `postToolUseFailure` | Observe | No host-consumed output fields |
| `beforeShellExecution` | Result | Allow/deny/ask; ask opens user approval |
| `afterShellExecution` | Observe | No host-consumed output fields; use for auditing and side effects |
| `beforeMCPExecution` | Result | Allow/deny/ask; prefer `failClosed: true` for security-critical gates; deferred for cloud agents |
| `afterMCPExecution` | Observe | No host-consumed output fields; cloud agents do not load MCP hooks |
| `beforeReadFile` | Result | Allow/deny; ask is coerced to deny |
| `afterFileEdit` | Observe | No host-consumed output fields; suitable for side effects such as formatting |
| `beforeTabFileRead` | Result | Tab-only allow/deny with no message fields; unavailable to cloud agents |
| `afterTabFileEdit` | Observe | Tab-only edit observation |
| `subagentStart` | Result | Allow/deny; ask is treated as deny |
| `subagentStop` | Result | Optional `followup_message`, consumed only when input `status` is `"completed"` |
| `stop` | Result | Optional `followup_message` |
| `afterAgentResponse` | Observe | hooks.json matcher value is `AgentResponse` |
| `afterAgentThought` | Observe | Optional `duration_ms`; hooks.json matcher value is `AgentThought` |
| `preCompact` | Result | Optional observational `user_message`; compaction cannot be blocked |

Cloud availability and matcher values are host constraints, not fields added by
wat. Exact exported event fields and builders remain documented in
`sdk/cursor` godoc.

For example, a workspace hook can return absolute plugin directories:

```go
func cursorWorkspace(
	ctx context.Context,
	event cursor.WorkspaceOpen,
	results cursor.WorkspaceOpenResults,
) (cursor.WorkspaceOpenOutput, error) {
	return results.PluginPaths([]string{"/abs/path/to/plugin"}), nil
}
```

## Permission and Ask / SoftDeny semantics

Cursor permission behavior is event-specific. Do not assume that one Ask or
Deny encoding works for every event.

| Event | Schema permissions | Ask / SoftDeny host behavior | Deny encoding (`sdk/cursor`) | Notes |
|---|---|---|---|---|
| `beforeShellExecution` | allow / deny / ask | **Enforced** (user approval) via `Ask` | `agent_message`, exit 2 by default; `WithUserMessage` for client-facing copy | Same Ask enforcement as `beforeMCPExecution` |
| `beforeMCPExecution` | allow / deny / ask | **Enforced** (user approval) via `Ask` | `agent_message`, exit 2 by default; `WithUserMessage` for client-facing copy | Contrasts with `preToolUse` |
| `preToolUse` | allow / deny / ask | **Not enforced** today — `Ask` still encodes `"ask"` | `agent_message`, exit 2 | Prefer `Allow`/`Deny` to gate |
| `beforeReadFile` | allow / deny (+ optional `user_message`) | **SDK:** no `Ask` builder — use `SoftDeny` (same as `Deny`). **Host:** a raw `"ask"` on the wire is coerced to deny | `user_message`, exit 0 (no `agent_message`) | Prefer `Deny` / `SoftDeny` |
| `beforeTabFileRead` | allow / deny only | **N/A** (no ask API) | Permission-only JSON, exit 0 (no message fields) | Chained message and `updated_input` helpers are ignored on the wire |
| `subagentStart` | allow / deny | **SDK:** no `Ask` builder — use `SoftDeny` (same as `Deny`). **Host:** a raw `"ask"` on the wire is treated as deny | `user_message`, exit 0 (no `agent_message`) | Prefer `Deny` / `SoftDeny`; exit 2 would re-wrap stdout as the user message |

Cursor defaults hook failures to fail-open. Set `failClosed: true` on the
`beforeMCPExecution` hooks.json handler when a crash, timeout, or invalid
response must block the tool.

## Portable projection boundaries

Portable registration expands onto Cursor as described in the
[cross-agent expansion matrix](../agent-formats.md#portable-registration-expansion).
The following Cursor-specific restrictions are particularly important:

- `OnPreTool` also observes dedicated shell, MCP, and read events. Updated input
  is emitted only for generic `preToolUse`.
- `OnPostTool` expands to `postToolUse`, `afterMCPExecution`, and
  `afterFileEdit`. The latter two are observation-only, so portable `Context`
  and `WithUpdatedOutput` have no host effect there.
- `afterShellExecution` is not projected onto portable `OnPostTool`; use the
  native observe handler for auditing. Generic `postToolUse` remains the route
  for shell post-tool context.
- Portable `OnPostToolFailure` context is discarded on Cursor because
  `postToolUseFailure` is observe-only.
- MCP output replacement is supported through generic `postToolUse`
  (`updated_mcp_tool_output`), not `afterMCPExecution`.
- Portable `OnPreCompact` remains observe-only and maps only shared compaction
  fields. Cursor-only metrics and native `UserMessage` output are not exposed
  through the portable API.

Observe-only portable events never emit host JSON.

## Stop follow-up loops

`StopResults.FollowUp` encodes a non-empty `followup_message` with exit 0, and
Cursor auto-submits it as the next user message.

- On `subagentStop`, Cursor consumes the message only when input `status` is
  `"completed"`. On `stop`, a non-empty message is always eligible.
- Input `loop_count` is how many automatic follow-ups the same script has
  already triggered for the conversation and starts at 0.
- Cursor enforces the per-script hooks.json `loop_limit` (default `5`; `null`
  means unlimited). Check `LoopCount` before returning another follow-up.

`loop_limit` is install configuration, not an SDK event field.

## Payload details

### Common envelope fields

Cursor's common input schema includes optional `model_id` and `model_params`
alongside the legacy `model` slug. `sdk/cursor.Envelope` exposes those fields,
using `ModelParam` for each parameter.

### Tool events

- `beforeMCPExecution` carries `tool_name` and `tool_input` plus either `url`
  for a remote MCP server or `command` for a stdio server.
- `preToolUse` may include input `agent_message`, the agent's pre-call
  narrative.
- `afterShellExecution` decodes `sandbox`; `postToolUseFailure` decodes
  `is_interrupt`.
- `afterShellExecution`, `afterMCPExecution`, `postToolUse`, and
  `postToolUseFailure` prefer documented `duration` when present, including
  explicit zero, and fall back to `duration_ms` only when it is absent.
- Agent `afterFileEdit` exposes `[]Edit` with `old_string` / `new_string`.
  Tab `afterTabFileEdit` exposes `[]TabEdit`, adding `range`, `old_line`, and
  `new_line`.

### Lifecycle, subagent, and compaction events

- `sessionStart` may include `composer_mode` (`"agent"`, `"ask"`, or
  `"edit"`).
- `sessionEnd` includes `duration_ms`, `final_status`, and optional
  `error_message` in addition to its reason and background-agent flag.
- `subagentStop` includes `description`, `duration_ms`, `message_count`,
  `tool_call_count`, and `modified_files`.
- `preCompact` includes `context_usage_percent`, `context_tokens`,
  `context_window_size`, `message_count`, `messages_to_compact`, and
  `is_first_compaction`.

### Live `subagentStart` compatibility

Payloads observed on Cursor 3.13.x differ from published Hooks schema examples:

- Automatic model selection may use `""`, `"auto"`, `"default"`, or
  `"inherit"` (case-insensitive after trimming). Treat those values as
  unpinned rather than concrete model IDs.
- A pinned Task sets both envelope `model` and `subagent_model` to the same
  concrete ID. Equality means pinned, not inherited.
- `subagent_type` may arrive as kebab-case (`general-purpose`) while hooks.json
  matchers use camelCase (`generalPurpose`). Normalize before comparison.
