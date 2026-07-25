# Using wat

This guide describes the CLI contract and the expected lifecycle of a `.wat/`
hook project.

## Requirements and installation

`wat` and generated hook projects require Go 1.26. Install the CLI from a
checkout with:

```bash
go install ./cmd/wat
```

The executable must contain module version or VCS build information because
`wat init` pins that version in the generated `.wat/go.mod`. Build normally
with Go's VCS stamping enabled; do not use `-buildvcs=false`.

Check the available commands with:

```bash
wat help
wat init -h
wat install -h
```

## Project discovery

A wat project is a `.wat/` directory containing regular files named `hooks.go`
and `go.mod`.

Commands that need a project walk upward from the current directory until they
find one. Set `WAT_PROJECT_DIR` to a workspace root to bypass that search:

```bash
WAT_PROJECT_DIR=/path/to/project wat doctor
```

The variable points to the project root, not to `.wat/` itself.

## Initialize a project

```bash
wat init [path]
```

Initialization:

- creates `.wat/.cache/`;
- creates `.wat/go.mod` if it is missing;
- writes `.wat/hooks.go`;
- runs `go mod tidy` inside `.wat/`.

Re-running `wat init` preserves an existing `go.mod` and refuses to overwrite
`hooks.go`. Use `wat init --force` only when replacing `hooks.go` is intended.

The generated file is an importable package, not a `main` package. Its required
contract is:

```go
var Hooks = []run.Hooks{
	agnostic.UseHooks().OnPreTool(preTool),
	cursor.UseHooks().BeforeShellExecution(cursorShell),
}
```

Registration expressions and package `init` functions run during install,
doctor, testing, and live invocation. Keep them deterministic and free of
external side effects. Perform runtime work inside handlers.

## Install hooks

```bash
wat install [--agent all|claude|copilot|cursor] [--wat-path path]
```

The default is `--agent all`. When `--wat-path` is omitted, `wat` resolves
itself from `PATH`.

Installation builds the hook project, inspects the exported `Hooks`, expands
portable handlers to native events, and reconciles only wat-managed command
entries in:

| Agent | Project configuration |
|---|---|
| Claude Code | `.claude/settings.json` |
| GitHub Copilot | `.github/hooks/wat.json` |
| Cursor | `.cursor/hooks.json` |

Other handlers and unrelated JSON fields are preserved. Wat-managed entries
for events no longer registered by the project are removed. Multiple authored
handlers for one native event still produce one installed command; they run in
registration order inside the generated binary.

Run `wat install` whenever the registration set changes. Handler implementation
changes do not require reinstalling because `wat run` rebuilds the cached
binary when sources change.

## Live hook execution

Installed commands have this shape:

```text
wat run --agent <dialect> --event <native-event>
```

`wat run` passes stdin, stdout, and stderr through to the generated hook
binary. `--agent` and `--event` identify the managed config entry for install
and doctor checks; the runtime decoder selects the dialect and event from the
payload itself.

On a cache miss, wat builds a bootstrap that imports `.wat/hooks.go` and calls
`run.Serve(Hooks...)`. The cache key includes all files under `.wat/` except
`.cache/`, the wat version, Go version, target OS/architecture, `GOFLAGS`,
`CGO_ENABLED`, and bootstrap format. The binary is stored below:

```text
.wat/.cache/<content-hash>/hooks[.exe]
```

Build failures are fail-open by default (exit 1). `wat run --fail-closed`
returns exit 2 for a build failure so permission-gating hosts can deny the
operation.

## Test with fixtures

```bash
wat test --agent <dialect> --fixture <file|-> [--verbose]
```

`--agent` and `--fixture` are required. A fixture must be non-empty JSON with a
registered native `hook_event_name`. The accepted `--event` flag does not
override the fixture; event identity is read from its payload. Use `-` to read
the fixture from stdin:

```bash
Get-Content fixture.json | wat test --agent cursor --fixture -
```

The report includes the fixture dialect/event, hook stdout, a recognized
decision field, and the hook exit code. `--verbose` also includes hook stderr.
The command returns the hook's exit code, so a deliberately denied fixture can
make the test command non-zero.

Repository examples live under [`testdata/fixtures/`](../testdata/fixtures/).
Use native event names and payload shapes for the selected dialect.

## Diagnose a project

```bash
wat doctor
```

Doctor prints `PASS`, `WARN`, and `FAIL` results for:

- Go availability and compatibility with `.wat/go.mod`;
- required project files and compilation;
- cache writability and warm/cold state;
- authored native registration manifest;
- installed config presence, disabled hooks, stale entries, duplicates, and
  command validity.

Warnings do not fail the command. One or more failed checks return exit 4 and
include a suggested fix when available.

## Handler execution semantics

For one native event:

1. handlers run in contribution and registration order;
2. `nil` or zero results add no output;
3. non-zero results merge into one native response;
4. conflicting replacement fields use the later value and emit a warning on
   stderr;
5. a terminal accumulated result (for example, deny/block or
   `continue: false`) prevents later handlers from running; an `Ask` result is
   not terminal;
6. the final result is encoded once.

Handler errors, decode errors, merge errors, and encode errors write a
diagnostic to stderr and return exit 1. Exact response fields and denial exit
codes remain agent-specific.

## Exit codes

| Code | Meaning |
|---|---|
| `0` | Successful CLI operation or hook execution |
| `1` | Invalid CLI usage; hook build failure in default mode; or hook runtime/protocol error |
| `2` | `wat run --fail-closed` build failure, or a native host denial that uses exit 2 |
| `3` | CLI runtime failure such as missing project, unreadable fixture, or failed install |
| `4` | At least one `wat doctor` check failed |

For `wat test`, a successfully executed fixture returns the hook binary's exit
code rather than translating it to a CLI-only code.

## Common problems

### No `.wat/` project found

Run the command below the intended project root, set `WAT_PROJECT_DIR` to that
root, or run `wat init`.

### The fixture has no registered handler

The fixture uses a native event name. Add the matching portable/native
registration, then confirm the expanded manifest with `wat doctor`.

### Agent config is stale

Run `wat install`, then `wat doctor`. Install reconciliation is scoped to
wat-managed command entries and will not remove unrelated hooks.

### A source change does not run

All files under `.wat/` participate in the cache key. Check that the command is
resolving the intended project and that `WAT_PROJECT_DIR` is not pointing
elsewhere.
