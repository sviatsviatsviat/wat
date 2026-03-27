## Context

**wat** is a Go CLI (`github.com/sviatsviatsviat/wat`, `go 1.26`) built from `cmd/wat`. Current product focus is `wat run` templating, tests, SemVer/VERSION, Keep a Changelog–style changelog, README, and GitHub Actions; fuller config is deferred.

**Packages:**

- **`internal/core`** — Host-neutral hook types and behavior: `HookContext`, `Command`, `HookHandler`, `HookHandlerResult` (see `hook_handler_result.go`), `HookHandlerFactory`.
- **`internal/run`** — The `run` wat subcommand as [core.Command] (`NewRunCommand`), flag parsing, argv templating for Cursor (`TemplateBindings` built from `cursor.CursorHookRunData[T]` on `HookContext.ParsedData`), and subprocess execution.
- **`internal/cursor`** — Cursor hook stdin models (`HookDataCommon`, `CursorHookRunData[T]`, event field structs), `CursorHookHandler[T]`, `HookHandlerBuilder` plumbing, `HookHandlerFactory`, and the per-event handler registry.

**I/O and execution:** User-facing stderr and hook stdout go through `cli.Console` (`cli.NewConsole` in production; `cli.NewMockConsole` in tests with `StderrBufferWriter` / `StdoutBufferWriter` and `StderrString`, `StdoutString`, etc.). Subprocesses use `watexec.NewRunner(stderr, console)` as `watexec.SubprocessRunner`, wired once in `app.Execute` and passed into `run.NewRunCommand`. The runner runs `exec.LookPath` on `args[0]` and falls back to the system shell when the first token is not a PATH executable (e.g. Windows shell builtins such as `dir` or `copy`).

**Hook stdin:** `app.Execute` reads hook stdin with `cli.ReadHookStdinJSON` (shared JSON-object validation; whitespace-only input may yield a nil body). Cursor requires a non-empty JSON payload and treats empty or missing JSON as an error; other hosts may allow empty bodies. For Cursor, **`internal/cursor`** orchestrates the hook path: `HookHandlerFromJSON` unmarshals shared stdin into **`HookDataCommon`** via **`NewHookDataCommon`**, builds the **`HookHandler`** with parsed **`CursorHookRunData[T]`** on **`HookContext.ParsedData`**, then the wired command runs **`Command.Execute`** (for **`run`**, **`internal/run`** type-switches on concrete `*CursorHookRunData[T]` and builds **`template.TemplateBindings`**).

**Build and CI:** For **Cursor hooks**, prefer a repo-root binary (`go build -o wat ./cmd/wat`; gitignored as `/wat` and `*.exe`). For **local test builds**, use `bin/` (`go build -o bin/wat ./cmd/wat`). CI lives under `.github/workflows` (e.g. `ci.yml`).

**Runtime hook pipeline (current design):** Parse argv for `wat <host> <subcommand> …` (`parseHost` / `parseSubcommand` in `app`), obtain a host `HookHandlerFactory`, read hook stdin JSON and `HookHandlerFromJSON`, then build the subcommand `Command` from the argv tail (e.g. `run.NewRunCommand` parses `run` flags), then `Handle(cmd)` → `core.HookHandlerResult` (`Code`, `Output`). `app.Execute` writes `Output` with `console.Write` (hook stdout, including when empty) and returns `Code`. The host handler sets `Output` (Cursor default and event handlers use the `{}` + newline protocol from `internal/cursor`). Early failures (unsupported host, stdin JSON parse error, `HookHandlerFromJSON` error) go to stderr only; hook stdout stays empty.

## Rules

- **CLI and help:** Follow Go-style tooling: short description, `Usage:` block, command list. Use raw multiline strings when they keep output clearer.
- **Tests:** Assert real behavior; no empty or placeholder tests. When validating the real binary path, prefer building and running `cmd/wat` over only calling internal helpers.
- **Naming:** Use names that match roles (`entrypoint`, `console`, `subprocessRunner`, `parseCommandAndCommonParameters`, dedicated factory filenames). Avoid cryptic short names (`d`, `rest`) and single-letter parameters for important API values. Split growing factory maps/registries into dedicated files (e.g. hook handler builders beside the factory). For `cli.MockConsole`, name buffers and helpers after the streams (`stdoutBuf`, `StderrString` / `StdoutString`); do not relabel stdout/stderr.
- **Exit codes:** Use named exit-code constants; document exit codes in README in a dedicated list. Keep semantics aligned with the product (e.g. unsupported host name, empty command after templating as bad input where specified).
- **Layering:** Keep host-neutral abstractions in `internal/core`. Put host-specific JSON, payload types, handler factories, and Cursor hook protocol details in each host package (e.g. `internal/cursor`), not in `cli`. Template placeholder binding for **`… run`** lives in **`internal/run`** (Cursor-shaped input today). Avoid one monolithic stdin envelope for every host. Enforce host-specific stdin rules in that host’s layer (e.g. Cursor non-empty JSON); shared readers alone are not enough.
- **Preserve the pipeline contract:** Do not break the orchestration above: early errors stderr-only and empty hook stdout; successful path writes `Output` via `console.Write`.
- **Host argv:** The first positional argument is the hook host name (e.g. `cursor`); the second is the wat subcommand. Remaining argv is passed into subcommand construction (e.g. `run` parses `-f` / `--file-pattern` from that slice).
- **Help on errors:** Root help for shared/common parsing failures; subcommands print their own help for subcommand-specific validation (e.g. empty run template when building the run command).
- **Wiring:** Prefer explicit factories in `cmd/wat` over blank imports and global registration for `HookHandlerFactory` implementations.
- **API surface:** Keep types, methods, and constructors unexported unless another package needs them. No dead code or exports only for tests if the same behavior can be asserted from package-internal tests. Put small shared utilities in `internal/helpers` instead of duplicating on domain types.
