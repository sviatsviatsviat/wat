# SDK public API

Wat exposes six public packages. Hook projects normally use
`sdk/agnostic` plus `sdk/run`, adding a native SDK only for agent-specific
behavior.

```text
github.com/sviatsviatsviat/wat/sdk/agnostic
github.com/sviatsviatsviat/wat/sdk/agnostic/tools
github.com/sviatsviatsviat/wat/sdk/claude
github.com/sviatsviatsviat/wat/sdk/copilot
github.com/sviatsviatsviat/wat/sdk/cursor
github.com/sviatsviatsviat/wat/sdk/run
```

## Hook project contract

`.wat/hooks.go` must be an importable Go package that exports:

```go
var Hooks []run.Hooks
```

Every value returned by an SDK's `UseHooks` implements `run.Hooks`. Registrars
are fluent and may combine portable and native handlers:

```go
var Hooks = []run.Hooks{
	agnostic.UseHooks().
		OnPreTool(preTool).
		OnStop(stop),
	claude.UseHooks().
		PermissionRequest(claudePermission),
	cursor.UseHooks().
		WorkspaceOpen(func(ctx context.Context, hook cursor.WorkspaceOpen, r cursor.WorkspaceOpenResults) (cursor.WorkspaceOpenOutput, error) {
			return r.PluginPaths([]string{"/abs/path/to/plugin"}), nil
		}),
}
```

The concrete registrar types are intentionally unexported. Call their exported
methods directly from `UseHooks`; do not store or expose the concrete type.

## Handler and result pattern

Result-producing handlers have three arguments:

```go
func(context.Context, Event, Results) (Output, error)
```

Observe-only handlers have two:

```go
func(context.Context, Event) error
```

The event is a typed value, not a generic envelope. Construct outputs through
the injected `Results` builder:

```go
func preTool(
	ctx context.Context,
	event agnostic.PreToolEvent,
	results agnostic.PreToolResults,
) (agnostic.PreToolResult, error) {
	if event.Tool == nil {
		return nil, nil
	}
	if input, ok := event.Tool.Input.AsBash(); ok && input.Command == "git clean -fdx" {
		return results.Deny("destructive clean is not allowed"), nil
	}
	return nil, nil
}
```

`nil` means no opinion and produces no hook stdout by itself. Do not implement
portable result interfaces yourself or reuse a result produced by a different
handler's builder. Adapters validate that results came from the injected native
builder.

When a result exposes fluent `With*` methods, call them on a builder-produced
value:

```go
return results.Allow().WithUpdatedInput(map[string]any{
	"command": "git push --dry-run",
}), nil
```

## Portable API: `sdk/agnostic`

Portable handlers fan out to native registrations for all three agents.

| Registration | Event | Result API | Behavior |
|---|---|---|---|
| `OnSessionStart` | `SessionStartEvent` | `Context(text)` | Add startup context |
| `OnSessionEnd` | `SessionEndEvent` | observe-only | Observe session completion |
| `OnUserPrompt` | `UserPromptEvent` | observe-only | Observe submitted prompt text |
| `OnPreTool` | `PreToolEvent` | `Allow`, `Deny`, `Ask`; `WithUpdatedInput` | Gate or rewrite a tool call (`Ask` is not enforced on Cursor `preToolUse` today) |
| `OnPostTool` | `PostToolEvent` | `Context`; `WithUpdatedOutput` | Add context or replace supported output (Cursor `afterFileEdit` / `afterMCPExecution` are observe-only; builders are no-ops there) |
| `OnPostToolFailure` | `PostToolFailureEvent` | `Context` | Add recovery guidance (ignored on Cursor; see agent formats) |
| `OnSubagentStart` | `SubagentStartEvent` | observe-only | Observe subagent creation |
| `OnSubagentStop` | `StopEvent` | `FollowUp` | Gate subagent completion |
| `OnStop` | `StopEvent` | `FollowUp` | Gate main-agent completion |
| `OnPreCompact` | `PreCompactEvent` | observe-only | Observe context compaction |

Portable `OnPreCompact` stays observe-only even though Cursor's native
`preCompact` can return an optional observational `user_message`. Use
`cursor.UseHooks().PreCompact` with `Results.UserMessage` when that host
message is required. Cursor-only compaction metrics
(`context_usage_percent`, `context_tokens`, `context_window_size`,
`message_count`, `messages_to_compact`, `is_first_compaction`) are likewise
available only on the native Cursor event, not on portable `Compact`.

### Normalized event shape

Every portable event embeds a shared envelope:

| Field | Meaning |
|---|---|
| `Agent` | `"claude"`, `"copilot"`, or `"cursor"` |
| `Name` | Native `hook_event_name` |
| `Session` | Native session or conversation identifier |
| `Cwd` | Working directory when supplied |
| `TranscriptPath` | Transcript path when supplied |

Kind-specific leaf pointers are optional because not every native payload
supplies the same detail:

- tool events use `Tool` and, after execution, `Result`;
- lifecycle events use `Life`;
- stop events use `Turn` and optionally `Subagent`;
- compaction uses `Compact`;
- user prompts expose `Prompt`.

`ToolCall` and `ToolResult` are named aliases at the package root. The other
leaf values are currently exposed through event fields but do not have
package-root type aliases; consumers can read their exported fields without
naming or constructing those internal model types. Always nil-check optional
pointers.

### `ToolCall` and portable inputs

`ToolCall` preserves both normalized and native data:

| Field | Meaning |
|---|---|
| `Name` | Canonical name such as `bash`, `edit`, `write`, or `read` |
| `Native` | Original native name |
| `Input` | `tools.Input` with raw and typed access |
| `ID` | Native tool-call identifier when available |
| `Shell` | Extracted shell command when available |
| `MCP` | Whether the call is recognized as MCP |

Compare names with constants from `sdk/agnostic/tools`, not string literals.
Typed accessors currently cover:

| Constant | Accessor |
|---|---|
| `ToolBash` | `AsBash()` |
| `ToolEdit` | `AsEdit()` |
| `ToolWrite` | `AsWrite()` |
| `ToolRead` | `AsRead()` |
| `ToolGlob` | `AsGlob()` |
| `ToolGrep` | `AsGrep()` |
| `ToolTask` | `AsTask()` |
| `ToolWebFetch` | `AsWebFetch()` |
| `ToolWebSearch` | `AsWebSearch()` |

`Input.Raw()` returns a copy of the native JSON. `Name()`, `Native()`, and
`HasRaw()` expose its binding metadata.

## Native SDKs

Native SDKs use the exact host event names, payload types, result builders, and
output semantics. Their common entry pattern is:

```go
claude.UseHooks().PreToolUse(handler)
copilot.UseHooks().PreToolUse(handler)
cursor.UseHooks().BeforeShellExecution(handler)
```

The following are the currently registerable fluent methods.

### Claude Code

| Domain | Methods |
|---|---|
| Session/prompt | `SessionStart`, `SessionEnd`, `UserPromptSubmit`, `UserPromptExpansion` |
| Tool/permission | `PreToolUse`, `PostToolUse`, `PostToolUseFailure`, `PermissionRequest`, `PermissionDenied` |
| Agent/task/stop | `SubagentStart`, `SubagentStop`, `TaskCreated`, `TaskCompleted`, `Stop` |
| UI/compact/workspace | `Notification`, `MessageDisplay`, `PreCompact`, `WorktreeCreate`, `Elicitation` |

Claude also exports typed decode models for additional native events that do
not yet have `UseHooks` registration methods. Treat the fluent method list as
the supported author-facing registration surface.

### GitHub Copilot

| Domain | Methods |
|---|---|
| Session/prompt | `SessionStart`, `SessionEnd`, `UserPromptSubmitted` |
| Tool/permission | `PreToolUse`, `PostToolUse`, `PostToolUseFailure`, `PermissionRequest` |
| Agent/stop | `SubagentStart`, `SubagentStop`, `AgentStop` |
| Other | `PreCompact`, `Notification`, `ErrorOccurred` |

### Cursor

| Domain | Methods |
|---|---|
| Session/prompt | `SessionStart`, `SessionEnd`, `WorkspaceOpen`, `BeforeSubmitPrompt` |
| Generic tools | `PreToolUse`, `PostToolUse`, `PostToolUseFailure` (observe-only) |
| Dedicated tools | `BeforeShellExecution`, `AfterShellExecution` (observe-only), `BeforeMCPExecution`, `AfterMCPExecution` (observe-only), `BeforeReadFile`, `AfterFileEdit` (observe-only), `BeforeTabFileRead`, `AfterTabFileEdit` (observe-only) |
| Agent/stop/compact | `SubagentStart`, `SubagentStop`, `Stop`, `AfterAgentResponse`, `AfterAgentThought`, `PreCompact` |

`AfterShellExecution`, `AfterMCPExecution`, `AfterFileEdit`,
`PostToolUseFailure`, and `AfterTabFileEdit` are observe-only where Hooks docs
list no consumed output. Rewrite MCP tool output with `PostToolUse`
(`updated_mcp_tool_output`). Cloud agents do not load MCP hooks.

Cursor `AfterFileEdit` uses `[]Edit` (`old_string` / `new_string`). Cursor
`AfterTabFileEdit` uses `[]TabEdit`, which adds Tab-specific `range`,
`old_line`, and `new_line` fields from the native payload.

`SubagentStop` decodes Cursor's documented telemetry fields (`description`,
`duration_ms`, `message_count`, `tool_call_count`, `modified_files`).
Cursor `Stop` / `SubagentStop` `FollowUp` auto-submits a next user message
(`followup_message`); for `subagentStop`, Cursor only consumes it when the
input `status` is `"completed"`. Cursor caps those loops with the hooks.json
`loop_limit` option (default `5`; `null` means unlimited). Use `LoopCount` on
the event (starts at 0) before returning another follow-up. See
[Agent protocols](agent-formats.md).

`WorkspaceOpen` is result-capable: handlers receive `WorkspaceOpenResults` and
may return `pluginPaths` (absolute plugin directories for the workspace). It is an
app lifecycle hook for the Cursor desktop app and CLI, skipped when there are
zero workspace folders, and does not run in cloud agents.

Cursor `PreCompact` decodes the full documented compaction metrics and exposes
`UserMessage` for an optional observational stdout message. Compaction cannot
be blocked. See [Agent protocols](agent-formats.md) for how this narrows
relative to portable `OnPreCompact`.

### Native constants and helpers

Each native package exports:

- `Dialect`;
- `Event*` constants for its known native event names;
- `Tool*` constants and a native `Input` type with agent-specific `As*`
  accessors;
- typed envelopes, events, result builders, and outputs;
- stable decode error sentinels such as `ErrEmptyPayload`,
  `ErrDecodePayload`, and `ErrEventNameRequired`;
- documented process exit constants.

Use `go doc` for complete event fields and builder methods:

```bash
go doc github.com/sviatsviatsviat/wat/sdk/claude PreToolUse
go doc github.com/sviatsviatsviat/wat/sdk/copilot PreToolResults
go doc github.com/sviatsviatsviat/wat/sdk/cursor PreToolUseResults
go doc github.com/sviatsviatsviat/wat/sdk/cursor PermissionResults
```

## Process API: `sdk/run`

### `Serve`

```go
func Serve(hooks ...Hooks)
```

`Serve` reads one payload from `os.Stdin`, detects the dialect, dispatches the
matching native event handlers, writes one merged response to `os.Stdout`, and
terminates the process with the native exit code. It calls `os.Exit`, so it is
an executable boundary rather than an embeddable long-running server.

Standalone hook executables can use it directly:

```go
func main() {
	run.Serve(
		agnostic.UseHooks().OnPreTool(preTool),
	)
}
```

The `wat` CLI generates an equivalent bootstrap for `.wat/` projects.

### `Inspect`

```go
func Inspect(hooks ...Hooks) Manifest
```

`Inspect` contributes the same registrations without reading stdin or invoking
handlers. `Manifest.Registrations` contains sorted native
`Dialect`/`Event`/`HandlerCount` entries. `EventsFor` returns registered events
for one dialect, and `Has` checks a dialect/event pair.

Use this API for tooling that needs the expanded native registration set.
Installation must not infer native coverage from portable method names.

## Composition and conflicts

Multiple `run.Hooks` values targeting the same dialect append to the same
handler bag. A portable handler and a native handler can therefore both receive
one native event.

Outputs merge in order. Additive context is combined according to native
semantics. Replacement fields use the later value and emit a warning. Terminal
outputs stop dispatch. These rules are implemented by native output types, so
an agent may have stricter behavior than the portable subset.

See [Agent protocols](agent-formats.md) for the portable-to-native expansion
matrix and known protocol differences.
