# Agent hook formats

Reference for tool names, MCP naming, and payload conventions across Claude Code, GitHub Copilot, and Cursor. Use this when implementing normalization, codecs, matchers, or `wat port` config translation.

Sources: [GitHub Copilot hooks reference](https://docs.github.com/en/copilot/reference/hooks-reference), [copilot-sdk#869](https://github.com/github/copilot-sdk/issues/869).

## Per-agent SDK skeleton

`claude`, `copilot`, and `cursor` are standalone packages (stdlib only) with the same layout. Each can be used without `agnostic`. `agnostic` depends on them: `On*` registration fans adapter handlers onto each agent SDK via package-level `On*` helpers, wrapping native events and Results builders.

Hook logic is organized as **vertical slices at the package root** — one file per native `hook_event_name`, including the typed `On*` registration helper and chain method for that event. Shared decode helpers live in `internal/hookkit`. The shared `Event` interface (`EventName()` only) is defined in `hookkit` and re-exported as `run.Event`; per-agent SDKs do not define their own `Event` type. Typed handler context is `run.Hook`.

| Location | Role |
|----------|------|
| `doc.go` | Package overview |
| `<event>.go` | Event struct, output, results, `On*` helper + chain method, decode registration, encode for one hook |
| `registry.go` | Event-name constants, alias tables |
| `envelope.go` | Shared payload fields (embedded on each event; access via promoted fields) |
| `exit.go` | Handler/encode exit-code constants |
| `encode.go` | `Encode` router (wire mapping) |
| `register.go` | Dialect codec, `RegisterDialect` init, handler registration helpers |
| `chain.go` | Unexported fluent `chain` handle (obtained only via package-level `On*`) |
| `config.go` | Native hook config types (`Handler`, `Settings`/`File`) |
| `errors.go` | Decode error sentinels |
| `tools/` | Event-bound tool input (`Input` with `AsBash`, `AsWrite`, …) |

**Agnostic** uses the same root vertical-slice layout for portable hook kinds (`pretool.go`, `stop.go`, …):

| Location | Role |
|----------|------|
| `<kind>.go` | Portable kind slice (`PreToolEvent`, `OnPreTool`, adapters, …) |
| `chain.go` | Unexported fluent `chain` handle (obtained only via package-level `On*`) |
| `toolcall.go` / `result.go` | Shared leaf types (`ToolCall`, `ToolResult`) and portable result projection |
| `internal/model/<kind>.go` | Leaf definitions behind root aliases (`*Event`, `*Hook`, `*Handler`, `*Result`); same kind filenames as the package root |
| `internal/model/toolcall.go` / `leaf.go` / `envelope.go` | Shared leaf payloads (`ToolCall`, `Subagent`, …) and `Envelope` |
| `tools/` | Canonical tool names and typed `Input` with `AsBash`, `AsWrite`, … |

Shared wire shapes may live in dedicated root files (e.g. `stop.go`, `permission.go`, `common.go`) when multiple events reuse the same output type.

**Intentional protocol differences** (do not expect parity):

- **Event count** — Claude exposes ~30 events; Copilot exposes 13; Cursor exposes 21.
- **Wire format** — Claude, Copilot, and Cursor payloads include `hook_event_name`; Copilot uses PascalCase event names with snake_case fields.
- **Encode contract** — Copilot and Cursor `Encode` return `([]byte, exitCode, error)`; Claude `Encode` returns `([]byte, error)` with blocking in JSON fields.
- **Side effects** — Claude `SessionStartOutput.Env` writes `CLAUDE_ENV_FILE`; Copilot and Cursor pass `env` in stdout JSON.
- **Config schema** — Claude `settings.json` (`Settings`) vs Copilot/Cursor `hooks.json` (`File`).

## Tool name normalization

Inbound mappers map native names onto a canonical vocabulary via `hookkit.NormalizeToolName` (internal). `ToolCall.Name` carries the result; `ToolCall.Native` always keeps the original string.

| Agent | Surface | Builtin example | Normalized | MCP example | MCP detection |
|-------|---------|-----------------|------------|-------------|---------------|
| Claude Code | `PreToolUse` | `Bash` | `bash` | `mcp__github__create_issue` | `mcp__` prefix → `mcp=true`, name unchanged |
| Copilot | PascalCase `PreToolUse` | `Bash` / `bash` | `bash` | `mcp__github__create_issue` or structured MCP metadata | `mcp__` prefix or codec metadata |
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

Typed handlers receive kind-specific events that embed a shared envelope:

- `Agent` — string dialect name (`claude.Dialect`, `copilot.Dialect`, or `cursor.Dialect`)
- `Name` — native hook event name as received

Config porting uses an internal `Kind` taxonomy (`KindPreTool`, `KindStop`, …) in `cmd/wat/internal/portconfig/model` — not part of the hook-author SDK.

See `go doc github.com/sviatsviatsviat/wat/sdk/agnostic` for typed events (`PreToolEvent`, `StopEvent`, …) and leaf structs (`ToolCall`, `Lifecycle`, …).

## Dialect identification

Each per-agent SDK exports a `Dialect` string constant (`claude.Dialect` = `"claude"`, and likewise for copilot/cursor) used for `sdk/run` registration and `Envelope.Agent`. CLI and config parse agent names via `cmd/wat/internal/dialect.Parse` (aliases like `claude-code`, `gh`). Hook serve resolves dialect via `run.WithDialect` / `WAT_AGENT`, or by walking each per-agent SDK’s registered `DialectOps.Detect` in `sdk/run` (payload shape and agent env hints such as `CURSOR_VERSION`).

## Portable agnostic API

Typed registration methods (`OnPreTool`, `OnPostTool`, `OnStop`, and others) accept only **portable** event kinds — those present on Claude Code, GitHub Copilot, and Cursor. Each kind has its own result type (`PreToolResult`, `PostToolResult`, `StopResult`, …) so hook authors can only set fields every agent can encode. Use `sdk/claude`, `sdk/copilot`, or `sdk/cursor` directly for agent-only capabilities.

Observe-only kinds (`SessionEnd`, `UserPrompt`, `PreCompact`, `SubagentStart`) take per-kind observe handlers that return only `error` — no hook response. Each handler receives a hook wrapper (typed event + `run.Invocation`).

### Handler signatures

All four SDKs use the same handler shapes. Result-producing handlers take `(ctx, hook, results)`; observe-only handlers take `(ctx, hook) error`. Access event fields on the hook (embedded typed event in agnostic; `hook.Event` in per-agent SDKs). Use `hook.Invocation()` for serve-time settings (`Dialect`, `Getenv`, `DialectConfig`).

| Category | agnostic | claude / copilot / cursor |
|---|---|---|
| Result | `(ctx, hook PreToolHook, r PreToolResults) (PreToolResult, error)` — hook embeds **`PreToolEvent`** (normalized) | `(ctx, hook Hook[PreToolUse], r PreToolUseResults) (PreToolUseOutput, error)` — hook carries **native typed event** via `hook.Event` |
| Observe | `(ctx, hook SessionEndHook) error` — hook embeds **`SessionEndEvent`** | `(ctx, hook Hook[SessionEnd]) error` — hook carries native typed event via `hook.Event` |

### Hook-scoped result builders

Each result-producing `On*` registration injects a builder interface scoped to that hook only — one exported builder type per `On*` helper, even when method sets are identical. Shared unexported implementation is fine; shared exported interfaces are not. Examples:

- `PostToolUseFailureResults` exposes only `Context` (not `Block` from post-success hooks).
- Claude `SubagentStartResults`, `NotificationResults`, and `PreCompactResults` are separate types (not a shared `CommonResults`).
- Cursor `BeforeReadFileResults` and `BeforeTabFileReadResults` are separate from shell/MCP `PermissionResults` where encode surfaces differ.

Use builder methods for common verbs (`r.Deny`, `r.Context`, `r.FollowUp`, …). Return `nil` for no opinion (silent stdout). Set advanced fields with fluent `With*` methods on the value returned by the builder (for example `r.Allow().WithUpdatedInput(args)`). Construct results only via the injected `*Results` builders (and `With*`); host-specific wrappers live in `sdk/agnostic/internal/{claude,cursor,copilot}`. Cursor notes: `WithUpdatedInput` emits `updated_input` only for `preToolUse`; `WithUpdatedOutput` maps to `updated_mcp_tool_output` (MCP post-tool only).

### Event support matrix

| Unified `Kind` | Claude | Copilot | Cursor | Portable handler |
|---|---|---|---|---|
| `SessionStart` | yes | yes | yes | `OnSessionStart` |
| `SessionEnd` | yes | yes | yes | `OnSessionEnd` (observe-only) |
| `UserPrompt` | yes | yes | yes | `OnUserPrompt` (observe-only) |
| `PreTool` | yes | yes | yes (+ dedicated pre-events) | `OnPreTool` |
| `PostTool` | yes | yes | yes (+ dedicated post-events) | `OnPostTool` |
| `PostToolFailure` | yes | yes | yes | `OnPostToolFailure` |
| `SubagentStart` | yes | yes | yes | `OnSubagentStart` (observe-only) |
| `SubagentStop` | yes | yes | yes | `OnSubagentStop` |
| `Stop` | yes | yes | yes | `OnStop` |
| `PreCompact` | yes | yes | yes | `OnPreCompact` (observe-only) |
| `PermissionRequest` | yes | yes | no | no — use `OnPermissionRequest` on `sdk/claude` / `sdk/copilot` |
| `Notification` | yes | yes | no | no — use `OnNotification` on `sdk/claude` / `sdk/copilot` |
| `AgentError` | yes | yes | no | no — decode-only typed events (`StopFailure`, `errorOccurred`) |
| `Other` | yes | yes | yes | no |

Observe-only handlers accept decoded events but produce no hook stdout JSON.

### Portable result types

| Kind | Result type | Builder / fields |
|---|---|---|
| `PreTool` | `PreToolResult` | `Allow`/`Deny`/`Ask`; `WithUpdatedInput` |
| `PostTool` | `PostToolResult` | `Context`; `WithUpdatedOutput` |
| `PostToolFailure` | `PostToolFailureResult` | `Context` |
| `Stop`, `SubagentStop` | `StopResult` | `FollowUp` |
| `SessionStart` | `SessionStartResult` | `Context` |
| `SessionEnd`, `UserPrompt`, `PreCompact`, `SubagentStart` | — | observe-only (no result) |

Each result-producing handler receives a hook-scoped builder interface (`PreToolResults`, `PostToolResults`, `StopResults`, and others) as its third parameter. See **Hook-scoped result builders** above. Multiple handlers for the same kind merge at the native JSON layer.

### Agent-only capabilities

Register with the per-agent SDK when you need features outside the portable intersection:

| Capability | Agents | SDK |
|---|---|---|
| `BlockPrompt`, `SetTitle` | Claude, Cursor (prompt) | `sdk/claude`, `sdk/cursor` |
| `Env` | Claude (`CLAUDE_ENV_FILE`), Cursor (stdout JSON) | `sdk/claude`, `sdk/cursor` |
| `HaltSession` | Claude (broad), Copilot (`permissionRequest` interrupt) | `sdk/claude`, `sdk/copilot` |
| `UserMessage` | Claude (`systemMessage`), Cursor (`user_message`) | `sdk/claude`, `sdk/cursor` |
| `PermissionRequest` handlers | Claude, Copilot | `sdk/claude`, `sdk/copilot` |
| `Context` on `PreTool` | Claude only | `sdk/claude` |
| `Decision` on `SubagentStart` | Cursor only | `sdk/cursor` |

### Per-agent On* coverage

Each per-agent SDK exposes package-level `On*` helpers (and fluent chain methods) with one entry per native hook surface (or shared builder where wire encode is identical). Claude-only long-tail events decode but have no dedicated portable `On*` — handle them with `sdk/claude` `On*` helpers when available.

| SDK | On* helpers | Notes |
|---|---|---|
| **claude** | `OnPreToolUse`, `OnPostToolUse`, `OnPostToolUseFailure`, `OnPermissionRequest`, `OnUserPromptSubmit`, `OnStop`, `OnSubagentStop`, `OnSessionStart`, `OnSubagentStart`, `OnNotification`, `OnPreCompact`, `OnSessionEnd`, plus Claude-only surfaces (`OnElicitation`, …) | ~18 additional Claude events (`Setup`, `TaskCreated`, …) decode as long-tail |
| **copilot** | `OnPreToolUse`, `OnPostToolUse`, `OnPostToolUseFailure`, `OnAgentStop`, `OnSubagentStop`, `OnPermissionRequest`, `OnSessionStart`, `OnSubagentStart`, `OnNotification`, `OnSessionEnd`, `OnUserPromptSubmitted`, `OnPreCompact`, `OnErrorOccurred` | Full 13-event surface covered |
| **cursor** | All 21 native events including `OnPreToolUse`, dedicated shell/MCP/file/tab hooks, lifecycle/telemetry observe helpers | Full 21-event surface covered |

## Claude inbound mapping

Portable `On*` handlers fan out onto `sdk/claude` via package-level `On*` helpers with unexported inbound mapping in `sdk/agnostic`; native decode, encode, and exit behavior stay in `sdk/claude`. Blocking is expressed via JSON fields with exit code 0 (Claude ignores exit 2 with JSON).

### Event name mapping

| Claude `hook_event_name` | `Kind` |
|---|---|
| `SessionStart`, `SessionEnd`, `UserPromptSubmit`, `PreToolUse`, `PostToolUse`, `PostToolUseFailure`, `PermissionRequest`, `SubagentStart`, `SubagentStop`, `Stop`, `PreCompact`, `Notification`, `StopFailure` | Normalized (see `go doc agnostic Kind`) |
| All other Claude events (`Setup`, `UserPromptExpansion`, `PostToolBatch`, `PermissionDenied`, `TaskCreated`, `TaskCompleted`, `TeammateIdle`, `MessageDisplay`, `InstructionsLoaded`, `ConfigChange`, `CwdChanged`, `FileChanged`, `WorktreeCreate`, `WorktreeRemove`, `PostCompact`, `Elicitation`, `ElicitationResult`, …) | `KindOther` |

`PreToolUse` with `tool_name: "Bash"` extracts `tool_input.command` into `ToolCall.Shell`.

### Encode surfaces

Portable Results wrappers delegate to `sdk/claude` `*Results` builders. Full Claude response shapes (including `BlockPrompt`, `Env`, `HaltSession`, `SetTitle`, `UserMessage`, and `PermissionRequest`) are available through `sdk/claude` directly.

| Event kind | Portable fields → Claude JSON |
|---|---|
| PreTool | `Decision`, `Reason` → `hookSpecificOutput.permissionDecision`; `UpdatedInput` |
| PostTool | `UpdatedOutput` → `updatedToolOutput`; `Context` → `additionalContext` |
| PostToolFailure | `Context` → `additionalContext` |
| Stop / SubagentStop | `FollowUp` → top-level `decision:block` + `reason` |
| SessionStart | `Context` → `additionalContext` |

Agent-native encode surfaces (use `sdk/claude` directly):

| Event kind | Key fields |
|---|---|
| UserPrompt | `BlockPrompt`, `Context`, `SetTitle` |
| PermissionRequest | `Decision`, `UpdatedInput`, `HaltSession`, `Context` |
| SessionStart | `Env` → `$CLAUDE_ENV_FILE` append |
| Any | `HaltSession` → `continue:false`; `UserMessage` → `systemMessage` |

### Exit codes

| Constant | Value | When |
|---|---|---|
| `claude.HandlerErrorExit` | `1` | Runner should use when a handler returns an error under fail-open (default) |
| `claude.FailBlockExit` | `2` | Runner should use when `WithFailPolicy(FailBlock)` is active |

## Copilot inbound mapping

Portable `On*` handlers fan out onto `sdk/copilot` via package-level `On*` helpers with unexported inbound mapping in `sdk/agnostic` (PascalCase `hook_event_name`, snake_case fields). Native decode stays in `sdk/copilot`.

### Wire format

Payloads require `hook_event_name` (for example `PreToolUse`, `Stop`). Shared fields use snake_case (`session_id`, `tool_name`, `tool_input`, …). Timestamps are RFC3339 strings.

Wire `Stop` always decodes as `AgentStop`. When the host scopes the stop to a subagent, `agent_name` / `agent_display_name` are set on that payload (`AgentStop.IsSubagent`); there is no decode-time remap to `SubagentStop`. Explicit `hook_event_name: "SubagentStop"` still decodes as `SubagentStop`. Portable `OnStop` receives only agent-scoped `Stop`; portable `OnSubagentStop` receives explicit `SubagentStop` and subagent-scoped `Stop`.

### Event name mapping

| Copilot event | `Kind` |
|---|---|
| `SessionStart` | `KindSessionStart` |
| `SessionEnd` | `KindSessionEnd` |
| `UserPromptSubmit` | `KindUserPrompt` |
| `PreToolUse` | `KindPreTool` |
| `PostToolUse` | `KindPostTool` |
| `PostToolUseFailure` | `KindPostToolFailure` |
| `PermissionRequest` | `KindPermissionRequest` |
| `SubagentStart` | `KindSubagentStart` |
| `SubagentStop` | `KindSubagentStop` |
| `Stop` | `KindStop` |
| `PreCompact` | `KindPreCompact` |
| `Notification` | `KindNotification` |
| `ErrorOccurred` | `KindAgentError` |

`PreToolUse` with `tool_name: "Bash"` (or `bash`) extracts shell `command` into `ToolCall.Shell`.

### Encode surfaces

Portable Results wrappers delegate to `sdk/copilot` `*Results` builders. Full Copilot response shapes are available through `sdk/copilot` directly.

| Event kind | Portable fields → Copilot JSON | Exit code |
|---|---|---|
| PreTool | `Decision`, `Reason` → `permission_decision`, `permission_decision_reason`; `UpdatedInput` → `modified_args` | `0` |
| PostTool | `UpdatedOutput` → `modified_result`; `Context` → `additional_context` | `0` |
| Stop / SubagentStop | `FollowUp` → `decision:block` + `reason` | `0` |
| PostToolFailure | `Context` → stdout text (recovery guidance) | `2` |
| SessionStart | `Context` → `additional_context` | `0` |

Agent-native encode surfaces (use `sdk/copilot` directly):

| Event kind | Key fields | Exit code |
|---|---|---|
| PermissionRequest | `Decision`, `Reason`, `HaltSession` → `behavior`, `message`, `interrupt` | `2` on deny |
| SubagentStart / Notification | `Context` → `additional_context` | `0` |

### Exit codes

| Constant | Value | When |
|---|---|---|
| `copilot.HandlerErrorExit` | `1` | Runner should use when a handler returns an error under fail-open (default) |
| `copilot.PreToolErrorExit` | `1` | Same value; use when a `PreToolUse` handler returns an error (fail-closed deny) |
| `copilot.WarnExit` | `2` | `Encode` returns this for documented `postToolUseFailure` context paths |

## Cursor inbound mapping

Portable `On*` handlers fan out onto `sdk/cursor` via package-level `On*` helpers with unexported inbound mapping in `sdk/agnostic`. Native decode stays in `sdk/cursor`.

Dedicated shell, MCP, and file events are **folded** into unified pre/post tool kinds so one `KindPreTool` handler receives shell, MCP, and read events with `Tool.Shell` / `Tool.MCP` populated. The native event name stays in `Event.Name`.

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
| `beforeTabFileRead` | `KindOther` | tab surface — permission-gating only, not folded |
| `afterTabFileEdit` | `KindOther` | tab surface — not folded |
| `workspaceOpen` | `KindOther` | app lifecycle — not folded |

`preToolUse` with `tool_name: "Shell"` extracts `tool_input.command` into `ToolCall.Shell`.

### Encode surfaces

Portable Results wrappers delegate to `sdk/cursor` `*Results` builders. Full Cursor response shapes are available through `sdk/cursor` directly.

| Event kind | Portable fields → Cursor JSON | Exit code |
|---|---|---|
| PreTool / dedicated pre-events | `Decision` → `permission`; `Reason` → `agent_message`; `UpdatedInput` → `updated_input` (preToolUse only) | `2` on deny |
| PostTool | `UpdatedOutput` → `updated_mcp_tool_output`; `Context` → `additional_context` | `0` |
| PostToolFailure | `Context` → `additional_context` | `0` |
| Stop / SubagentStop | `FollowUp` → `followup_message` | `0` |
| SessionStart | `Context` → `additional_context` | `0` |

Agent-native encode surfaces (use `sdk/cursor` directly):

| Event kind | Key fields | Exit code |
|---|---|---|
| SubagentStart / `beforeTabFileRead` | `Decision`, `UserMessage`, `Reason`, `UpdatedInput` | `2` on deny |
| UserPrompt | `BlockPrompt`, `UserMessage` → `continue:false` | `0` |
| SessionStart | `Env` → `env` | `0` |
| PreCompact | `UserMessage` | `0` |

### Exit codes

| Constant | Value | When |
|---|---|---|
| `cursor.HandlerErrorExit` | `1` | Runner should use when a handler returns an error under fail-open (default) |
| `cursor.PermissionDenyExit` | `2` | `Encode` returns this for permission-gating deny |

## Related code

- Config porting: `wat port --from` / `--to` (see [`cmd/wat/port.go`](../cmd/wat/port.go))
- Dialect: per-agent `Dialect` constants in [`sdk/claude`](../sdk/claude/), [`sdk/copilot`](../sdk/copilot/), [`sdk/cursor`](../sdk/cursor/); CLI name parsing in [`cmd/wat/internal/dialect`](../cmd/wat/internal/dialect/)
- Tests: [`cmd/wat/internal/dialect`](../cmd/wat/internal/dialect/)
- Normalization: [`internal/hookkit/toolname.go`](../internal/hookkit/toolname.go) — `NormalizeToolName`; [`sdk/agnostic/tools`](../sdk/agnostic/tools/) — `Input` with `AsBash`, `AsWrite`, and related accessors
- Tests: [`internal/hookkit/toolname_test.go`](../internal/hookkit/toolname_test.go); [`sdk/agnostic/tools`](../sdk/agnostic/tools/) (`input_test.go`)
- Port kind/event registries: [`cmd/wat/internal/portconfig/`](../cmd/wat/internal/portconfig/) (`claude/`, `copilot/`, `cursor/`)
- Serve / fan-out: [`sdk/agnostic/runner_test.go`](../sdk/agnostic/runner_test.go)
- Cursor SDK: [`sdk/cursor/`](../sdk/cursor/) — typed events, `Encode`, `On*` registration into [`sdk/run`](../sdk/run/), `sdk/cursor/tools` event-bound tool input (`AsShell`, …)
- Tests: [`sdk/cursor/cursor_test.go`](../sdk/cursor/cursor_test.go)
