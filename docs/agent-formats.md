# Agent hook formats

Reference for tool names, MCP naming, and payload conventions across Claude Code, GitHub Copilot, and Cursor. Use this when implementing normalization, codecs, matchers, or portconfig translation.

Sources: [GitHub Copilot hooks reference](https://docs.github.com/en/copilot/reference/hooks-reference), [copilot-sdk#869](https://github.com/github/copilot-sdk/issues/869).

## Tool name normalization

`agenthooks.NormalizeToolName` maps native names to a canonical vocabulary. `ToolCall.Native` always keeps the original string.

| Agent | Surface | Builtin example | Normalized | MCP example | MCP detection |
|-------|---------|-----------------|------------|-------------|---------------|
| Claude Code | `PreToolUse` | `Bash` | `bash` | `mcp__github__create_issue` | `mcp__` prefix → `mcp=true`, name unchanged |
| Copilot | PascalCase `PreToolUse` | `Bash` | `bash` | `mcp__github__create_issue` | same as Claude (Claude-format matchers) |
| Copilot | camelCase `preToolUse` | `bash` | `bash` | `my-server-list_items` | codec sets `ToolCall.MCP` from structured metadata (not inferred from hyphens) |
| Cursor | `preToolUse` / dedicated shell hooks | `Shell` | `bash` | `MCP:browser_navigate` | `MCP:` prefix → `mcp=true`, name unchanged |

### Builtin alias examples

| Native (any case where noted) | Normalized |
|-------------------------------|------------|
| `Bash`, `bash`, `Shell`, `powershell` | `bash` |
| `Edit`, `edit`, `notebookedit` | `edit` |
| `Write`, `create` | `write` |
| `Read`, `view` | `read` |
| `Agent`, `task` | `task` |
| `web_fetch` | `web_fetch` |

Unknown names pass through unchanged with `mcp=false` unless an `mcp__` or `MCP:` prefix matches.

### Out of scope for `NormalizeToolName`

- **Copilot `preMcpToolCall`** — separate hook event with `serverName` and bare `toolName` fields, not a combined tool name string.
- **Cursor dedicated events** (`beforeShellExecution`, `beforeMCPExecution`, …) — folded into unified `Kind` by codecs; native event name stays in `Event.Name`.

## Unified event envelope

All codecs produce `agenthooks.Event`:

- `Agent` — `Dialect` (Claude, Copilot, Cursor)
- `Kind` — normalized category (`KindPreTool`, `KindStop`, …)
- `Name` — native hook event name as received
- `Raw` — untouched native JSON payload

See `go doc github.com/sviatsviatsviat/wat/agenthooks Event` for sub-structs (`Tool`, `Result`, `Life`, …).

## Dialect detection

`agenthooks.Detect` infers the originating agent from a hook stdin payload and environment hints. `agenthooks.ParseDialect` parses explicit names from CLI flags (`claude`, `copilot`, `cursor` and aliases). When the dialect is already known, skip `Detect` and use the explicit value.

Payload shape is checked before environment variables because Cursor exports `CLAUDE_PROJECT_DIR` for Claude Code compatibility.

| Step | Signal | Result |
|------|--------|--------|
| 1 | `cursor_version` or `conversation_id` in JSON | Cursor |
| 2 | `sessionId` (camelCase) | Copilot camelCase |
| 3 | `hook_event_name` + `timestamp` | Copilot VS Code (Claude payloads carry no `timestamp`) |
| 4 | `session_id` | Claude Code |
| 5a | `CURSOR_VERSION` env | Cursor |
| 5b | `CLAUDE_PROJECT_DIR` env | Claude Code |
| 5c | `COPILOT_HOME` env | Copilot |
| else | — | Unknown |

**Ambiguous cases:** Copilot payloads always win over `CLAUDE_PROJECT_DIR` in the environment. Env-only detection with only `CLAUDE_PROJECT_DIR` may misidentify Copilot in repos that also have `.claude/settings.json`; prefer payload evidence or an explicit `--agent` override.

## Claude codec

`agenthooks.ClaudeCodec` decodes Claude Code hook stdin into `agenthooks.Event` and encodes `agenthooks.Result` into Claude stdout JSON. Blocking is expressed via JSON fields with exit code 0 (Claude ignores exit 2 with JSON).

### Event name mapping

| Claude `hook_event_name` | `Kind` |
|---|---|
| `SessionStart`, `SessionEnd`, `UserPromptSubmit`, `PreToolUse`, `PostToolUse`, `PostToolUseFailure`, `PermissionRequest`, `SubagentStart`, `SubagentStop`, `Stop`, `PreCompact`, `Notification`, `StopFailure` | Normalized (see `go doc agenthooks Kind`) |
| All other Claude events (`Setup`, `UserPromptExpansion`, `PostToolBatch`, `PermissionDenied`, `TaskCreated`, `TaskCompleted`, `TeammateIdle`, `MessageDisplay`, `InstructionsLoaded`, `ConfigChange`, `CwdChanged`, `FileChanged`, `WorktreeCreate`, `WorktreeRemove`, `PostCompact`, `Elicitation`, `ElicitationResult`, …) | `KindOther` — full payload preserved in `Event.Raw` |

`PreToolUse` with `tool_name: "Bash"` extracts `tool_input.command` into `ToolCall.Shell`.

### Encode surfaces

| Event kind | Key `Result` fields → Claude JSON |
|---|---|
| PreTool | `Decision`, `Reason` → `hookSpecificOutput.permissionDecision`; `UpdatedInput`, `Context` |
| UserPrompt | `BlockPrompt` / `Decision` → top-level `decision:block`; `Context`, `SetTitle` |
| Stop / SubagentStop | `FollowUp` → top-level `decision:block` + `reason`; `Context` |
| SessionStart | `Context`, `SetTitle`; `Env` → append `export KEY='value'` lines to `$CLAUDE_ENV_FILE` |
| Any | `HaltSession` → `continue:false`; `UserMessage` → `systemMessage` |

When `$CLAUDE_ENV_FILE` is set, env keys must match `[A-Za-z_][A-Za-z0-9_]*`; invalid keys return an encode error before any file write. When `$CLAUDE_ENV_FILE` is unset, `Result.Env` on SessionStart is a no-op (not an error) and key validation is skipped.

## Copilot codec

`agenthooks.CopilotCodec` decodes GitHub Copilot hook stdin in **camelCase CLI** or **VS Code compatible** (PascalCase event name, snake_case fields) format and encodes `agenthooks.Result` into flat camelCase stdout JSON.

### Wire formats

| Signal in payload | Format | Event name source |
|---|---|---|
| `sessionId` | camelCase CLI | Config key via `eventHint` (required except `notification`, which carries `hook_event_name`) |
| `hook_event_name` | VS Code compatible | Payload (`PreToolUse`, `Stop`, …) |

Timestamps decode as ms-epoch numbers (camelCase) or ISO-8601 strings (VS Code).

### Event name mapping

| Copilot event | `Kind` |
|---|---|
| `sessionStart` / `SessionStart` | `KindSessionStart` |
| `sessionEnd` / `SessionEnd` | `KindSessionEnd` |
| `userPromptSubmitted` / `UserPromptSubmit` | `KindUserPrompt` |
| `preToolUse` / `PreToolUse` | `KindPreTool` |
| `postToolUse` / `PostToolUse` | `KindPostTool` |
| `postToolUseFailure` / `PostToolUseFailure` | `KindPostToolFailure` |
| `permissionRequest` / `PermissionRequest` | `KindPermissionRequest` |
| `subagentStart` / `SubagentStart` | `KindSubagentStart` |
| `subagentStop` / `SubagentStop` | `KindSubagentStop` |
| `agentStop` / `Stop` | `KindStop` |
| `preCompact` / `PreCompact` | `KindPreCompact` |
| `notification` / `Notification` | `KindNotification` |
| `errorOccurred` / `ErrorOccurred` | `KindAgentError` |

`preToolUse` with `toolName: "bash"` or VS Code `tool_name: "Bash"` extracts shell `command` into `ToolCall.Shell`.

### Encode surfaces

| Event kind | Key `Result` fields → Copilot JSON | Exit code |
|---|---|---|
| PreTool | `Decision`, `Reason` → `permissionDecision`, `permissionDecisionReason`; `UpdatedInput` → `modifiedArgs` | `0` |
| PostTool | `UpdatedOutput` → `modifiedResult`; `Context` → `additionalContext` | `0` |
| Stop / SubagentStop | `FollowUp` → `decision:block` + `reason` | `0` |
| PermissionRequest | `Decision`, `Reason`, `HaltSession` → `behavior`, `message`, `interrupt` | `2` on deny |
| PostToolFailure | `Context` → stdout text (recovery guidance) | `2` |
| SessionStart / SubagentStart / Notification | `Context` → `additionalContext` | `0` |

### Exit codes

| Constant | Value | When |
|---|---|---|
| `CopilotPreToolErrorExit` | `1` | Runner should use when a `preToolUse` handler returns an error (fail-closed deny) |
| `CopilotWarnExit` | `2` | `Encode` returns this for documented `permissionRequest` deny and `postToolUseFailure` context paths |

Copilot-specific limitations (`BlockPrompt`, `Env`, most `HaltSession` cases) are reported via `Unsupported`.

## Cursor codec

`agenthooks.CursorCodec` decodes Cursor hook stdin into `agenthooks.Event` and encodes `agenthooks.Result` into Cursor stdout JSON.

Dedicated shell, MCP, and file events are **folded** into unified pre/post tool kinds so one `KindPreTool` handler receives shell, MCP, and read events with `Tool.Shell` / `Tool.MCP` populated. The native event name stays in `Event.Name`; the full payload stays in `Event.Raw`.

### Event name mapping

| Cursor event | `Kind` | Folding notes |
|---|---|---|
| `sessionStart` | `KindSessionStart` | |
| `sessionEnd` | `KindSessionEnd` | |
| `beforeSubmitPrompt` | `KindUserPrompt` | |
| `preToolUse` | `KindPreTool` | |
| `postToolUse` | `KindPostTool` | |
| `postToolUseFailure` | `KindPostToolFailure` | |
| `beforeShellExecution` | `KindPreTool` | → `Tool.Name=bash`, `Tool.Shell=command` |
| `afterShellExecution` | `KindPostTool` | → bash + terminal `output` in `Result.Text` |
| `beforeMCPExecution` | `KindPreTool` | → `Tool.MCP=true` |
| `afterMCPExecution` | `KindPostTool` | → MCP + `Result.Text=result_json` |
| `beforeReadFile` | `KindPreTool` | → `Tool.Name=read`; file content in `Tool.Input` |
| `afterFileEdit` | `KindPostTool` | → `Tool.Name=edit`; diffs in `Tool.Input` / `Result.Raw` |
| `subagentStart` | `KindSubagentStart` | |
| `subagentStop` | `KindSubagentStop` | `LoopCount` on `Subagent` |
| `stop` | `KindStop` | `LoopCount` on `Turn` |
| `preCompact` | `KindPreCompact` | |
| `afterAgentResponse` | `KindOther` | observe-only |
| `afterAgentThought` | `KindOther` | observe-only |
| `beforeTabFileRead` | `KindOther` | tab surface — not folded |
| `afterTabFileEdit` | `KindOther` | tab surface — not folded |
| `workspaceOpen` | `KindOther` | app lifecycle — not folded |

`preToolUse` with `tool_name: "Shell"` extracts `tool_input.command` into `ToolCall.Shell`.

### Encode surfaces

| Event kind | Key `Result` fields → Cursor JSON | Exit code |
|---|---|---|
| PreTool / dedicated pre-events / SubagentStart / `beforeTabFileRead` | `Decision` → `permission`; `UserMessage` → `user_message`; `Reason` → `agent_message`; `UpdatedInput` → `updated_input` (preToolUse only) | `2` on deny |
| UserPrompt (`beforeSubmitPrompt`) | `BlockPrompt` / `DecisionDeny` → `continue:false`; `UserMessage` → `user_message` | `0` |
| PostTool | `UpdatedOutput` → `updated_mcp_tool_output`; `Context` → `additional_context` | `0` |
| Stop / SubagentStop | `FollowUp` → `followup_message` | `0` |
| SessionStart | `Env` → `env`; `Context` → `additional_context` | `0` |
| PreCompact | `UserMessage` → `user_message` | `0` |

### Exit codes

| Constant | Value | When |
|---|---|---|
| `CursorHandlerErrorExit` | `1` | Runner should use when a handler returns an error under fail-open (default) |
| `CursorWarnExit` | `2` | `Encode` returns this for permission-gating deny |

Cursor-specific limitations (`HaltSession`, `SetTitle`, `Ask` on non-shell preToolUse) are reported via `Unsupported`.

## Related code

- Config porting: [`agenthooks/portconfig`](../agenthooks/portconfig/) — `Parse`, `Emit`, `Translate`
- Detection: [`agenthooks/dialect.go`](../agenthooks/dialect.go) — `ParseDialect`, `Detect`
- Tests: [`agenthooks/dialect_test.go`](../agenthooks/dialect_test.go)
- Normalization: [`agenthooks/event.go`](../agenthooks/event.go) — `NormalizeToolName`, `InputAs`
- Tests: [`agenthooks/event_test.go`](../agenthooks/event_test.go)
- Claude codec: [`agenthooks/claude.go`](../agenthooks/claude.go) — `ClaudeCodec`, `CodecFor`
- Tests: [`agenthooks/claude_test.go`](../agenthooks/claude_test.go)
- Copilot codec: [`agenthooks/copilot.go`](../agenthooks/copilot.go) — `CopilotCodec`, `CopilotPreToolErrorExit`, `CopilotWarnExit`
- Tests: [`agenthooks/copilot_test.go`](../agenthooks/copilot_test.go)
- Cursor codec: [`agenthooks/cursor.go`](../agenthooks/cursor.go) — `CursorCodec`, `CursorWarnExit`, `CursorHandlerErrorExit`
- Tests: [`agenthooks/cursor_test.go`](../agenthooks/cursor_test.go)
