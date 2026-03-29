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
        "command": "wat cursor exec echo __HOOK_EVENT_NAME__"
      }
    ]
  }
}
```

### `wat <host> exec`

Run a templated hook subprocess: read hook JSON from stdin, substitute allowed `__PLACEHOLDER__` tokens in the command template, run that process, and write the host’s hook protocol line on success. The **first** argument is the hook host (e.g. `cursor`); the **second** is the wat subcommand (`exec`).

```text
Usage:

	wat <host> exec <command> [templated arguments]
	wat <host> exec [-f <re>] <command> [templated arguments]
	wat <host> exec [--file-pattern <re>] <command> [templated arguments]

Options (only for exec, before the subprocess template):

	-f, --file-pattern <re>
	                      Optional; default * means no filter. If you pass a
	                      non-* value, <re> must be non-empty (Go regexp syntax).
	                      When stdin bindings include __FILE_PATH__ (Cursor
	                      afterFileEdit), `exec` skips the subprocess if the path
	                      does not match <re>; other hook events ignore the flag
	                      for matching purposes.
```

Put `-f` / `--file-pattern` after `exec` and before the subprocess command (for example `wat cursor exec -f '[.]go$' …`). Flags are parsed when **`execcommand.NewExecHookHandlerProvider`** builds the provider. If equivalent options are repeated, **the last value wins**.

**Command template** — Everything after the optional flags is one command template: the subprocess program and its arguments. Use only `__PLACEHOLDER__` tokens documented under [Cursor exec template bindings](#cursor-exec-template-bindings); any other `__TOKEN__` in the template is an error (exit code `2`).

#### Cursor exec template bindings

Authoritative list of `__KEY__` segments for `wat cursor exec` (inner part between underscores). Optional JSON fields resolve to an empty string when missing or `null`.

**Common** — Available for every Cursor hook event that `exec` supports (including events that use the default adapter):

| Placeholder | Description |
|-------------|-------------|
| `__CONVERSATION_ID__` | From `conversation_id`. |
| `__GENERATION_ID__` | From `generation_id`. |
| `__MODEL__` | From `model`. |
| `__HOOK_EVENT_NAME__` | From `hook_event_name`. |
| `__CURSOR_VERSION__` | From `cursor_version`. |
| `__USER_EMAIL__` | From `user_email` when present (empty if missing or `null`). |
| `__TRANSCRIPT_PATH__` | From `transcript_path` when present (empty if missing or `null`). |

**Per event** — `exec` adds event-specific keys only when the hook adapter carries that data ([`exec_hook_handler_provider.go`](internal/execcommand/exec_hook_handler_provider.go)):

| Event | Additional placeholders |
|-------|-------------------------|
| `afterFileEdit` | `__FILE_PATH__` |
| `afterShellExecution` | `__DURATION__`, `__SANDBOX__` |
| Other registered events (default adapter) | None — common placeholders only. |

The built-in `wat … exec` help lists the union of placeholders across events; use the table above to see which tokens apply to the hook you are configuring.

**Exit status** — If the subprocess is started, wat exits with **that process’s exit code**. Otherwise wat uses own standard [Exit codes](#exit-codes).

### Exit codes

| Code | Meaning |
|------|---------|
| `0` | Success. For `exec`, this means the templated command exited `0`, or `exec` skipped the subprocess because `-f` / `--file-pattern` did not match `__FILE_PATH__`. |
| `1` | General failure — e.g. stdin JSON parse error, host/event rejected the payload, or the subprocess failed to run. |
| `2` | Bad input — invalid CLI usage, unknown host, unknown subcommand, missing `exec` command, unknown `__PLACEHOLDER__`, or nothing left to execute after templating. |

If `exec` **does** start a subprocess, the process exit code may match the child’s code, so `1` or `2` can mean either wat or the child; check stderr for context.

## Supported hosts

- **[Cursor](#cursor)** — supported today.

## Cursor

Cursor supplies hook JSON on stdin. Register hook commands in **`.cursor/hooks.json`**.

Shared and event-specific field types are defined in [`internal/cursor/hook_data.go`](internal/cursor/hook_data.go). Every event receives the shared **`HookDataCommon`** envelope (`conversation_id`, `generation_id`, `model`, `hook_event_name`, `cursor_version`, `workspace_roots`, optional `user_email`, optional `transcript_path`).

### Supported Cursor hook types

Each subsection describes what the **hook adapter** exposes from stdin for that `hook_event_name` (not which `exec` placeholders exist—see [Cursor exec template bindings](#cursor-exec-template-bindings)).

#### `afterShellExecution`

Fires after a shell command runs. **`AfterShellExecutionFields`** adds `command`, `output`, `duration`, and `sandbox` to the shared envelope.

**Returns** `{}`.

#### `afterMCPExecution`

Fires after MCP execution. Uses the default adapter with **no** separate event payload struct beyond **`HookDataCommon`**.

**Returns** `{}`.

#### `afterFileEdit`

Fires after a file edit. **`AfterFileEditFields`** adds `file_path` and `edits` (each edit is `old_string` / `new_string`).

**Returns** `{}`.

When `wat cursor exec …` includes `-f` / `--file-pattern` with a Go regexp, **`execcommand.NewExecHookHandlerProvider`** (in `internal/execcommand`) builds handlers whose **afterFileEdit** path applies the filter before invoking the subprocess when `__FILE_PATH__` is present in template bindings (other events omit that key, so the subprocess runs as usual). The regexp is matched against the hook’s `file_path` after path cleaning and normalizing separators to `/`.

#### `afterTabFileEdit`

Fires after a tab file edit. Uses the default adapter (**`HookDataCommon`** only).

**Returns** `{}`.

#### `afterAgentResponse`

Fires after an agent response. Uses the default adapter (**`HookDataCommon`** only).

**Returns** `{}`.

#### `afterAgentThought`

Fires after agent thought. Uses the default adapter (**`HookDataCommon`** only).

**Returns** `{}`.

#### `sessionEnd`

Fires when the session ends. Uses the default adapter (**`HookDataCommon`** only).

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
- **`Command`** — Subcommand implementation (`exec` today): `Execute(ctx *HookContext) int`, returning the process exit code.
- **Command argument placeholders (`internal/execcommand`)** — For each inner part of a `__KEY__` token in the subprocess template, the exec subcommand resolves a string via hook data; missing keys are treated as unknown placeholders (bad input). For Cursor, resolution is driven by `HookContext.ParsedData` (`*cursor.CursorHookRunData[T]` per event type `T`, or `T == struct{}` for common-only hooks).
- **`HookContext`** — Carries `HookHost` and `ParsedData` (`any`) into `Command.Execute`; the host handler sets both before `Execute`.

### Execution flow

1. **Entry** — **`cmd/wat`** `main` calls **`app.Execute`** with program arguments (minus binary name), stdin, stdout, and stderr; **`Execute`** constructs **`cli.Console`** for diagnostics and hook protocol output for the rest of the run.
2. **Host side** — The first argument selects the hook host; **`app`** builds a host **`HookHandlerFactory`** and keeps the remaining arguments for the wat subcommand.
3. **Hook handler** — **`cli.ReadHookStdinJSON`** reads hook event bytes from stdin, then **`HookHandlerFromJSON`** returns a **`HookHandler`** for that event (before the wat subcommand line is turned into a **`Command`**).
4. **Hook handler provider** — **`app.newHookHandlerProvider(subcommand, console, rest)`** builds a **`core.HookHandlerProvider`** (e.g. **`execcommand.NewExecHookHandlerProvider`**, which parses **`exec`** flags such as **`-f`** from **`rest`**).
5. **`app.Execute` → `HookHandler`** — **`Handle(hookCommand)`**; the handler sets **`HookContext`** (**`HookHost`**, **`ParsedData`**).
6. **`HookHandler` → `hookCommand` → `HookHandler` → `app.Execute`** — **`Execute(HookContext)`** (for **`exec`**: build placeholder bindings from **`ParsedData`**, expand template tokens, subprocess via **`runSubprocess`** with **`Console.ConnectErrorsFrom`**, …) returns the subprocess exit code; **`HookHandler`** returns **`HookHandlerResult`** (**`Output`**, **`Code`**) to **`app.Execute`**.
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
  Note right of C: exec: bindings from ParsedData, render, subprocess, …
  C-->>H: exit code
  H-->>A: HookHandlerResult (Output, Code)
  Note over A: Write Output to hook stdout,<br/>return Code
```

### Cursor hook factory and handler (`internal/cursor`)

This is how the **`HookHandlerFactory`** and **`HookHandler`** from the execution flow are implemented for Cursor today.

1. **Factory value** — **`cursor.NewHookHandlerFactory()`** returns **`cursor.HookHandlerFactory`** for **`HookHandlerFromJSON`** / event builders.
2. **`HookHandlerFromJSON`** — Rejects empty stdin (Cursor expects a JSON body). **`NewHookDataCommon`** unmarshals bytes into **`HookDataCommon`** (shared envelope: `conversation_id`, `hook_event_name`, etc.—see `hook_data.go`).
3. **Per-event dispatch** — **`hook_event_name`** selects an entry in **`cursorHookAdapterBuilders`** (`hook_adapter_builders.go`). Missing events return an error (“not supported yet”).
4. **Building the handler** — Each registered builder is a **`HookHandlerBuilder`** `func(rawJSON []byte, hookData HookDataCommon) (core.HookHandler, error)`. Most events use **`NewDefaultHookHandler`** ( **`CursorHookRunData[struct{}]`** , no event payload) or **`NewHookHandlerFromEventFields[T]`** (parses **`HookDataWithCommon[T]`** and builds **`CursorHookRunData[T]`** with **`EventSpecific: &Fields`**).
5. **`CursorHookHandler[T].Handle`** — Builds **`HookContext`** with **`HookHost`** (**`cursor.HookHostCursor`**) and **`ParsedData`** pointing at **`CursorHookRunData[T]`**, calls **`cmd.Execute(ctx)`**, and returns **`HookHandlerResult`** with the subprocess exit **`Code`** and fixed hook stdout **`Output`** (**`cursor.DefaultHookResponseLine`**, i.e. `{}` plus newline).
6. **Placeholder bindings in `internal/execcommand`** (`cursor_bindings_common.go`, `cursor_bindings_event.go`, and per-event extractor maps) — For Cursor, **`NewExecHookHandlerProvider`**’s **`HookHandlerFor`** selects an exec handler; bindings come from **`newTemplateBindingsCommon`** and/or **`templateBindingsFromCursorEventPayload`** (common keys such as **`CONVERSATION_ID`**, **`HOOK_EVENT_NAME`**, …—the inner part of each **`__KEY__`** token, plus event keys like **`FILE_PATH`** or **`DURATION`** / **`SANDBOX`** where defined). Optional JSON strings use **`helpers.StringFromPtr`**. Missing map keys mean the placeholder is unknown; known keys resolve even when the value is empty. Supporting a new adapter type in **`exec`** adds a **`case`** in **`HookHandlerFor`** and usually a small extractor map (no change to **`CursorHookRunData`**’s shape).
7. **Where bindings run** — For **`wat <host> exec`**, the exec **`HookHandler`** (from **`NewExecHookHandlerProvider`**) optionally filters on **`FILE_PATH`** when a **`-f`** pattern is set, then substitutes **`__KEY__`** segments in each template token using the bindings and collects any unknown keys; the exec subcommand turns unknowns into a bad-input exit.

```mermaid
sequenceDiagram
  participant F as cursor.HookHandlerFactory
  participant Data as hookDataCommon
  participant Map as cursorHookAdapterBuilders
  participant Build as HookAdapterBuilder

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
  participant Run as execHookHandler
  participant TB as exec_bindings
  participant R as placeholder_substitution

  H->>Run: Handle()
  Run->>TB: build bindings from ParsedData
  TB-->>Run: bindings
  Run->>R: substitute __KEY__ in argsTemplate
  loop each template token / __KEY__
    R->>TB: TemplateValue(inner key)
    TB-->>R: value, ok
  end
  R-->>Run: rendered args, unknown keys
  Note right of Run: then runSubprocess(console, rendered args)
```

### Other packages

- **`internal/cli`** — Console (stderr vs hook stdout, including **`ConnectErrorsFrom`** for child stderr), help text, exit code constants, shared hook stdin JSON read.
- **`internal/execcommand`** — Subcommand `exec` as **`core.HookHandlerProvider`** (`NewExecHookHandlerProvider`), `__KEY__` placeholder expansion in the command template, and subprocess execution (PATH lookup, shell fallback, child stdout discarded).
- **`internal/helpers`** — Small shared utilities.

### Extending wat

- **New host** — Add a package (like `internal/cursor`) implementing `HookHandlerFactory`, own JSON types, default hook stdout lines, and stdin policy. Register the factory in `app.newHookHandlerFactory`. Keep host protocol strings out of `internal/cli`.
- **New hook (event)** — For an existing host, register `hook_event_name` in that host’s adapter-builder map (e.g. `cursorHookAdapterBuilders` in `internal/cursor/hook_adapter_builders.go`), wiring an existing or new builder to a `HookAdapter`. Define a new event field struct **`T`** in `internal/cursor`, build **`CursorHookRunData[T]`** in the builder, and add a **`HookHandlerFor`** case in **`NewExecHookHandlerProvider`** (plus extractor maps) when `exec` must support new **`__KEY__`** tokens. Document the event in the README.
- **New subcommand** — Implement `core.HookHandlerProvider` under `internal/<subcommand>` (today `internal/execcommand`) and wire it in `app.newHookHandlerProvider`.
