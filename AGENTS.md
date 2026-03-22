## Context

**wat** is a Go CLI (`github.com/sviatsviatsviat/wat`, `go 1.26`) built from `cmd/wat`. Current product focus is `wat run` templating, tests, SemVer/VERSION, Keep a Changelog–style changelog, README, and GitHub Actions; fuller config is deferred.

**Packages:** Host-neutral hook types and behavior live in `internal/core` (`HookContext`, `TemplateBindings`, `Command`, `HookHandler`, `HookHandlerResult` in `hook_handler_result.go`, `HookHandlerFactory`). Subcommand constructors (e.g. `NewRunCommand`) live in `internal/commands`. Cursor-specific hook JSON, payloads, factories, default protocol response strings, and handler builders live in `internal/cursor`.

**I/O and execution:** User-facing stderr and hook stdout go through `cli.Console` (`cli.NewConsole` in production; `cli.NewMockConsole` in tests with `StderrBufferWriter` / `StdoutBufferWriter` and `StderrString`, `StdoutString`, etc.). Subprocesses use `watexec.NewRunner(stderr, console)` as `watexec.SubprocessRunner`, wired once in `app.Execute` and passed into `commands.NewRunCommand`. The runner runs `exec.LookPath` on `args[0]` and falls back to the system shell when the first token is not a PATH executable (e.g. Windows shell builtins such as `dir` or `copy`).

**Hook stdin:** `app.Execute` reads hook stdin with `cli.ReadHookStdinJSON` (shared JSON-object validation; whitespace-only input may yield a nil body). Cursor’s path requires a payload and treats empty or missing JSON as an error; other hosts may allow empty bodies. Cursor parses shared stdin fields inside `internal/cursor` and wires `TemplateBindings` before `Command.Execute`.

**Build and CI:** Prefer build output under a dedicated directory such as `bin/`, not the repo root. CI lives under `.github/workflows` (e.g. `ci.yml`).

**Runtime hook pipeline (current design):** Parse subcommand and shared host flags, build the `Command`, read hook stdin JSON, obtain a host `HookHandlerFactory` (e.g. `cursor.NewHookHandlerFactory()`), `HookHandlerFromJSON`, then `Handle(cmd)` → `core.HookHandlerResult` (`Code`, `Output`). `app.Execute` writes `Output` with `console.Write` (hook stdout, including when empty) and returns `Code`. The host handler sets `Output` (Cursor’s default sets `defaultHookResponseLine` after `Command.Execute`). Early failures (unsupported host, stdin JSON parse error, `HookHandlerFromJSON` error) go to stderr only; hook stdout stays empty.

## Rules

- **CLI and help:** Follow Go-style tooling: short description, `Usage:` block, command list. Use raw multiline strings when they keep output clearer.
- **Tests:** Assert real behavior; no empty or placeholder tests. When validating the real binary path, prefer building and running `cmd/wat` over only calling internal helpers.
- **Naming:** Use names that match roles (`entrypoint`, `console`, `subprocessRunner`, `parseCommandAndCommonParameters`, dedicated factory filenames). Avoid cryptic short names (`d`, `rest`) and single-letter parameters for important API values. Split growing factory maps/registries into dedicated files (e.g. hook handler builders beside the factory). For `cli.MockConsole`, name buffers and helpers after the streams (`stdoutBuf`, `StderrString` / `StdoutString`); do not relabel stdout/stderr.
- **Exit codes:** Use named exit-code constants; document exit codes in README in a dedicated list. Keep semantics aligned with the product (e.g. unsupported `--host`, empty command after templating as bad input where specified).
- **Layering:** Keep host-neutral abstractions in `internal/core`. Put host-specific JSON, payload types, handler factories, and default hook protocol strings in each host package (e.g. Cursor’s `defaultHookResponseLine` in `internal/cursor`), not in `cli`. Avoid one monolithic stdin envelope for every host. Enforce host-specific stdin rules in that host’s layer (e.g. Cursor non-empty JSON); shared readers alone are not enough.
- **Preserve the pipeline contract:** Do not break the orchestration above: early errors stderr-only and empty hook stdout; successful path writes `Output` via `console.Write`.
- **Shared host flags:** Parse `-H` / `--host` (with value) in `Execute` *after* the subcommand so usage is `wat <subcommand> -H <name> …`, not `wat -H <name> <subcommand>`. Pass only unconsumed argv into subcommand handlers.
- **Help on errors:** Root help for shared/common parsing failures; subcommands print their own help for subcommand-specific validation (e.g. empty run template when building the run command).
- **Wiring:** Prefer explicit factories in `cmd/wat` over blank imports and global registration for `HookHandlerFactory` implementations.
- **API surface:** Keep types, methods, and constructors unexported unless another package needs them. No dead code or exports only for tests if the same behavior can be asserted from package-internal tests. Put small shared utilities in `internal/helpers` instead of duplicating on domain types.
