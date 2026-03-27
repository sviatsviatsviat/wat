# wat

**wat** is a small tool that accelerates hook writing by taking care of common boilerplate. It reads hook input, passes its data to templated commands or guards, and writes result back to host. It is meant to be:

- **Cross-platform** — runs on Linux, macOS, and Windows.
- **Fast and lightweight** — implemented in Go to have minimal runtime delay.

## Usage

Sample **Cursor** hook configuration (`.cursor/hooks.json`). Adjust the path to `wat` and the child command for your setup.

```json
{
  "version": 1,
  "hooks": {
    "afterFileEdit": [
      {
        "command": "wat cursor run echo __HOOK_EVENT_NAME__"
      }
    ]
  }
}
```

### `wat <host> run`

Run a templated hook subprocess: read hook JSON from stdin, substitute allowed `__PLACEHOLDER__` tokens in the command template, run that process, and write the host’s hook protocol line on success. The **first** argv word is the hook host (e.g. `cursor`); the **second** is the wat subcommand (`run` today).

```text
Usage:

	wat <host> run <command> [templated arguments]
	wat <host> run [-f <re>] <command> [templated arguments]
	wat <host> run [--file-pattern <re>] <command> [templated arguments]

Options (only for run, before the subprocess template):

	-f, --file-pattern <re>
	                      Optional; default * means no filter. If you pass a
	                      non-* value, <re> must be non-empty (Go regexp syntax).
	                      When stdin bindings include __FILE_PATH__ (Cursor
	                      afterFileEdit), `run` skips the subprocess if the path
	                      does not match <re>; other hook events ignore the flag
	                      for matching purposes.
```

Put `-f` / `--file-pattern` after `run` and before the subprocess command (for example `wat cursor run -f '[.]go$' …`). Flags are parsed inside **`run.NewRunCommand`**. If equivalent options are repeated, **the last value wins**.

**Command template** — Everything after the optional flags is one command template: the subprocess program and its arguments. Use only `__PLACEHOLDER__` tokens documented for the current hook event in [Supported Cursor hook types](#supported-cursor-hook-types); any other `__TOKEN__` in the template is an error (exit code `2`).

**Exit status** — If the subprocess is started, wat exits with **that process’s exit code**. Otherwise wat uses own standard [Exit codes](#exit-codes).

### Exit codes

| Code | Meaning |
|------|---------|
| `0` | Success. For `run`, this means the templated command exited `0`, or `run` skipped the subprocess because `-f` / `--file-pattern` did not match `__FILE_PATH__`. |
| `1` | General failure — e.g. stdin JSON parse error, host/event rejected the payload, or the subprocess failed to run. |
| `2` | Bad input — invalid CLI usage, unknown host, unknown subcommand, missing `run` command, unknown `__PLACEHOLDER__`, or nothing left to execute after templating. |

If `run` **does** start a subprocess, the process exit code may match the child’s code, so `1` or `2` can mean either wat or the child; check stderr for context.

## Supported hosts

- **[Cursor](#cursor)** — supported today.

## Cursor

Cursor supplies hook JSON on stdin. Register hook commands in **`.cursor/hooks.json`**.

### Common Cursor placeholders

| Placeholder | Description |
|-------------|-------------|
| `__CONVERSATION_ID__` | Identifier for the current conversation; taken from the `conversation_id` property of the stdin JSON. |
| `__GENERATION_ID__` | Identifier for the generation step; taken from `generation_id`. |
| `__MODEL__` | Model name for the interaction; taken from `model`. |
| `__HOOK_EVENT_NAME__` | Which hook fired (for example `afterFileEdit`); taken from `hook_event_name`. |
| `__CURSOR_VERSION__` | Cursor app version string; taken from `cursor_version`. |
| `__WORKSPACE_ROOTS__` | Workspace root paths as a single string, joined with `;`; taken from the `workspace_roots` JSON array on stdin. |
| `__USER_EMAIL__` | Signed-in user email when present; taken from `user_email` (empty string if missing or `null`). |
| `__TRANSCRIPT_PATH__` | Transcript file path when present; taken from `transcript_path` (empty string if missing or `null`). |

### Supported Cursor hook types

#### `afterShellExecution`

Fires after a shell command runs in Cursor.

**Placeholders** — [Common Cursor placeholders](#common-cursor-placeholders) plus the additional placeholders below.

| Placeholder | Description |
|-------------|-------------|
| `__COMMAND__` | Full terminal command that Cursor executed; taken from `command`. |
| `__OUTPUT__` | Full terminal output captured by Cursor; taken from `output`. |
| `__DURATION__` | Duration in milliseconds spent executing the shell command; taken from `duration`. |
| `__SANDBOX__` | Whether the command ran in a sandboxed environment (`true` or `false`); taken from `sandbox`. |

**Returns** `{}`.

#### `afterMCPExecution`

Fires after MCP execution.

**Placeholders** — [Common Cursor placeholders](#common-cursor-placeholders).

**Returns** `{}`.

#### `afterFileEdit`

Fires after a file edit.

**Placeholders** — [Common Cursor placeholders](#common-cursor-placeholders) plus:

| Placeholder | Description |
|-------------|-------------|
| `__FILE_PATH__` | Absolute path of the edited file; taken from `file_path`. |

**Returns** `{}`.

When `wat cursor run …` includes `-f` / `--file-pattern` with a Go regexp, **`commands.runCommand.Execute`** applies the filter before invoking the subprocess when `__FILE_PATH__` is present in template bindings (other events omit that key, so the subprocess runs as usual). The regexp is matched against the hook’s `file_path` after path cleaning and normalizing separators to `/`.

#### `afterTabFileEdit`

Fires after a tab file edit.

**Placeholders** — [Common Cursor placeholders](#common-cursor-placeholders).

**Returns** `{}`.

#### `afterAgentResponse`

Fires after an agent response.

**Placeholders** — [Common Cursor placeholders](#common-cursor-placeholders).

**Returns** `{}`.

#### `afterAgentThought`

Fires after agent thought.

**Placeholders** — [Common Cursor placeholders](#common-cursor-placeholders).

**Returns** `{}`.

#### `sessionEnd`

Fires when the session ends.

**Placeholders** — [Common Cursor placeholders](#common-cursor-placeholders).

**Returns** `{}`.

## Development

Requires Go 1.26+ (see `go.mod`).

```bash
go test ./...
go vet ./...

# Local hook binary at repo root (gitignored; match the real filename in `.cursor/hooks.json`)
# Omit -o on Windows so the toolchain writes wat.exe in the current directory.
go build ./cmd/wat

# Build under bin/ (use a .exe suffix on Windows when using -o; -o uses the path literally)
go build -o bin/wat.exe ./cmd/wat
```

On **Windows**, `-o` is interpreted literally: `go build -o wat …` creates a file named `wat` with **no** `.exe`, which is a poor fit for hooks and `CreateProcess`.

- Prefer **`go build ./cmd/wat`** from the repo root (no `-o`) so the output is **`wat.exe`**, and point hooks at **`.\wat.exe`**.
- If you must pass `-o`, use **`-o wat.exe`**.

On **Unix**, `go build ./cmd/wat` writes **`wat`** in the current directory; use **`./wat`** in hooks.

CI runs `go test ./...`, `go vet ./...`, and `go build ./cmd/wat` across multiple `GOOS`/`GOARCH` targets.

### Versioning

This project follows [Semantic Versioning 2.0.0](https://semver.org/) and maintains a [Keep a Changelog](https://keepachangelog.com/en/1.1.0/) in `CHANGELOG.md`.

## Architecture overview

wat is layered so **hosts** (Cursor today) stay separate from **shared** CLI, templating, and subprocess execution.

### Core interfaces (`internal/core`)

These types define the host-neutral contract:

- **`HookHandlerFactory`** — Builds a `HookHandler` from raw hook stdin JSON bytes. The host chooses parsing, validation, and which events exist.
- **`HookHandler`** — Handles one invocation: receives the subcommand `Command`, fills `HookContext` (`HookHost`, host-specific `ParsedData`), calls `Command.Execute`, and returns `HookHandlerResult` (process exit `Code` and hook stdout `Output` string).
- **`Command`** — Subcommand implementation (`run` today): `Execute(ctx *HookContext) int`, returning the process exit code.
- **`TemplateBindings`** (`internal/template`) — `TemplateValue(key) (value, ok)` for keys matching the inner part of `__KEY__` in argv. If `ok` is false, `run` reports an unknown placeholder error. For Cursor, `run` builds this from `HookContext.ParsedData` (`*cursorcore.CursorHookRunData[T]` per event type `T`, or `T == struct{}` for common-only hooks).
- **`HookContext`** — Carries `HookHost` and `ParsedData` (`any`) into `Command.Execute`; the host handler sets both before `Execute`.

### Execution flow

1. **Entry** — **`cmd/wat`** `main` calls **`app.Execute`** with argv (minus program name), stdin, stdout, and stderr; **`Execute`** constructs **`cli.Console`** and **`watexec`** runner for the rest of the run.
2. **Host side** — The first argv word selects the hook host; **`app`** builds a host **`HookHandlerFactory`** and keeps the remaining argv slice for the wat subcommand.
3. **Hook handler** — **`cli.ReadHookStdinJSON`** reads hook event bytes from stdin, then **`HookHandlerFromJSON`** returns a **`HookHandler`** for that event (before the wat subcommand argv is turned into a **`Command`**).
4. **Command** — **`newHookCommand(subcommand, …, rest)`** builds **`hookCommand`** (`core.Command`, e.g. **`run.NewRunCommand`**, which parses **`run`** flags such as **`-f`** from **`rest`**).
5. **`app.Execute` → `HookHandler`** — **`Handle(hookCommand)`**; the handler sets **`HookContext`** (**`HookHost`**, **`ParsedData`**).
6. **`HookHandler` → `hookCommand` → `HookHandler` → `app.Execute`** — **`Execute(HookContext)`** (for **`run`**: build **`TemplateBindings`** from **`ParsedData`**, template render, **`watexec`** child, …) returns the subprocess exit code; **`HookHandler`** returns **`HookHandlerResult`** (**`Output`**, **`Code`**) to **`app.Execute`**.
7. **Finish** (diagram note over **`app.Execute`**) — write **`result.Output`** to hook stdout, return **`result.Code`** as the process exit code.

```mermaid
sequenceDiagram
  autonumber
  participant A as app.Execute
  participant F as HookHandlerFactory
  participant H as HookHandler
  participant C as hookCommand

  Note over A: Host factory, ReadHookStdinJSON,<br/>HookHandlerFromJSON, newHookCommand
  A->>F: HookHandlerFromJSON(hook event bytes)
  F-->>A: HookHandler
  A->>H: Handle(hookCommand)
  H->>C: Execute(HookContext)
  Note right of C: run: bindings from ParsedData, render, watexec, …
  C-->>H: exit code
  H-->>A: HookHandlerResult (Output, Code)
  Note over A: Write Output to hook stdout,<br/>return Code
```

### Cursor hook factory and handler (`internal/cursor`, `internal/cursor/core`)

This is how the **`HookHandlerFactory`** and **`HookHandler`** from the execution flow are implemented for Cursor today.

1. **Factory value** — **`cursor.NewHookHandlerFactory()`** returns **`cursor.HookHandlerFactory`** for **`HookHandlerFromJSON`** / event builders.
2. **`HookHandlerFromJSON`** — Rejects empty stdin (Cursor expects a JSON body). **`NewHookDataCommon`** unmarshals bytes into **`HookDataCommon`** (shared envelope: `conversation_id`, `hook_event_name`, etc.—see `hook_data.go`).
3. **Per-event dispatch** — **`hook_event_name`** selects an entry in **`cursorHookHandlerBuilders`** (`hook_handler_builders.go`). Missing events return an error (“not supported yet”).
4. **Building the handler** — Each registered builder is a **`HookHandlerBuilder`** `func(rawJSON []byte, hookData HookDataCommon) (core.HookHandler, error)`. Most events use **`NewDefaultHookHandler`** ( **`CursorHookRunData[struct{}]`** , no event payload) or **`NewHookHandlerFromEventFields[T]`** (parses **`HookDataWithCommon[T]`** and builds **`CursorHookRunData[T]`** with **`EventSpecific: &Fields`**).
5. **`CursorHookHandler[T].Handle`** — Builds **`HookContext`** with **`HookHost`** (**`cursorcore.HookHostCursor`**) and **`ParsedData`** pointing at **`CursorHookRunData[T]`**, calls **`cmd.Execute(ctx)`**, and returns **`HookHandlerResult`** with the subprocess exit **`Code`** and fixed hook stdout **`Output`** (**`cursorcore.DefaultHookResponseLine`**, i.e. `{}` plus newline).
6. **`TemplateBindings` in `internal/run`** (`cursor_bindings.go` plus `cursor_bindings_common.go`, `cursor_bindings_event.go`, and per-event files) — For Cursor, **`templateBindingsForCursor`** type-switches on **`*CursorHookRunData[T]`** and maps to **`template.TemplateBindings`**: common placeholders mirror shared stdin fields (**`CONVERSATION_ID`**, **`HOOK_EVENT_NAME`**, …—the inner names **`internal/template`** extracts from **`__KEY__`**). Event-specific keys (**`FILE_PATH`**, **`COMMAND`**, …) merge with common lookups. Optional JSON uses **`helpers.StringFromPtr`**; **`workspace_roots`** is joined with **`;`**. Missing map keys mean **`ok == false`** (unknown placeholder); known keys return **`ok == true`** even when the value is empty. Adding a new event type **`T`** adds a **`case`** branch (no change to **`CursorHookRunData`**’s shape).
7. **Where bindings run** — For **`wat <host> run`**, **`runCommand.Execute`** optionally filters on **`FILE_PATH`** when a **`-f`** pattern is set, then calls **`template.RenderTokens(argvTemplate, bindings)`**. **`RenderTokens`** scans each argv token for **`__KEY__`** substrings, calls **`TemplateValue`** per key, substitutes the returned string, and collects unknown keys; **`run`** turns any unknowns into a bad-input exit.

```mermaid
sequenceDiagram
  participant F as cursor.HookHandlerFactory
  participant Data as hookDataCommon
  participant Map as cursorHookHandlerBuilders
  participant Build as HookHandlerBuilder

  F->>F: require non-empty JSON
  F->>Data: NewHookDataCommon(bytes)
  F->>Map: lookup hook_event_name
  Map-->>F: hookHandlerBuilder
  F->>Build: builder(rawJSON, hookData)
  Build-->>F: HookHandler (e.g. CursorHookHandler)
```

```mermaid
sequenceDiagram
  participant H as CursorHookHandler
  participant Run as runCommand
  participant TB as templateBindingsForCursor
  participant R as template.RenderTokens

  H->>Run: Execute(HookContext ParsedData)
  Run->>TB: templateBindingsForCursor(CursorHookRunData)
  TB-->>Run: TemplateBindings
  Run->>R: RenderTokens(argvTemplate, bindings)
  loop each argv token / __KEY__
    R->>TB: TemplateValue(inner key)
    TB-->>R: value, ok
  end
  R-->>Run: rendered argv, unknown keys
  Note right of Run: then watexec.Run(rendered argv)
```

### Other packages

- **`internal/cli`** — Console (stderr vs hook stdout), help text, exit code constants, shared hook stdin JSON read.
- **`internal/run`** — Subcommand `run` as `core.Command` (`NewRunCommand`).
- **`internal/template`** — Replaces `__KEY__` tokens using `TemplateBindings`.
- **`internal/watexec`** — Subprocess runner (child stderr forwarded; child stdout discarded).
- **`internal/helpers`** — Small shared utilities.

### Extending wat

- **New host** — Add a package (like `internal/cursor`) implementing `HookHandlerFactory`, own JSON types, default hook stdout lines, and stdin policy. Register the factory in `app.newHookHandlerFactory`. Keep host protocol strings out of `internal/cli`.
- **New hook (event)** — For an existing host, register `hook_event_name` in that host’s handler-builder map (e.g. `cursorHookHandlerBuilders` in `internal/cursor/hook_handler_builders.go`), wiring an existing or new builder to a `HookHandler`. Define a new event field struct **`T`** in `cursorcore`, build **`CursorHookRunData[T]`** in the builder, and add a **`templateBindingsForCursor`** case for **`*CursorHookRunData[T]`** when `run` must support new **`__KEY__`** tokens. Document the event in the README.
- **New subcommand** — Implement `core.Command` under `internal/<subcommand>` (today `internal/run`) and wire argv construction in `app.newHookCommand`.
