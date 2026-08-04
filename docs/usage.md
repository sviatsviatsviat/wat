# Using wat

This guide describes the CLI contract and the expected lifecycle of a `.wat/`
hook project.

## Requirements and installation

`wat` and generated hook projects require Go 1.26. The CLI runs on Linux,
macOS, and Windows; CI exercises `go test ./...` and `go build ./cmd/wat` on
Ubuntu and Windows. Install the CLI with:

```bash
go install github.com/sviatsviatsviat/wat/cmd/wat@v0.4.0-alpha
```

From a repository checkout, `go install ./cmd/wat` also works. The executable
must contain module version or VCS build information because `wat init` pins
that version in the generated `.wat/go.mod`. Build normally with Go's VCS
stamping enabled; do not use `-buildvcs=false`.

Check the available commands with:

```bash
wat help
wat init -h
wat install -h
wat version
```

## Print the CLI version

```bash
wat version
```

`wat version` prints the same module version string that `wat init` pins in
`.wat/go.mod` and that the hook build cache key includes. Root flags
`-version` and `--version` are aliases.

Tagged installs print the release tag (for example `v0.4.0-alpha`). Local
builds with Go VCS stamping print a pseudo-version
`v0.0.0-<timestamp>-<revision>`. Builds without module or VCS version
information exit `3` with an error on stderr.

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
- creates `.wat/.gitignore` (ignores `.cache/`) when missing;
- creates `.wat/go.mod` if it is missing;
- writes `.wat/hooks.go`;
- writes starter fixtures and expect sidecars under `.wat/testdata/` when
  missing (one `session_start` case per agent);
- runs `go mod tidy` inside `.wat/`.

Re-running `wat init` preserves an existing `go.mod` and `.gitignore`. A missing
`.gitignore` is still created on a non-force re-run. Re-runs refuse to overwrite
`hooks.go` unless `--force` is set. Starter fixtures use write-if-missing, so
custom cases under `.wat/testdata/` are kept. Use `wat init --force` only when
replacing `hooks.go` is intended.

The generated file is an importable package, not a `main` package. Its required
contract is:

```go
var Hooks = []run.Hooks{
	agnostic.UseHooks().OnSessionStart(sessionStart),
}
```

`wat init` scaffolds a single `OnSessionStart` handler that injects observable
context (`"wat hooks are active"`) so you can confirm install and `wat test`
quickly. Add more `On*` registrations as needed.

Registration expressions and package `init` functions run during install,
doctor, testing, and live invocation. Keep them deterministic and free of
external side effects. Perform runtime work inside handlers.

## Install hooks

```bash
wat install [--agent all|claude|copilot|cursor] [--wat-path path]
```

The default is `--agent all`. When `--wat-path` is omitted, `wat` resolves
itself from `PATH` and writes an absolute path into each managed command.
Paths that contain whitespace or shell metacharacters are double-quoted with
escapes so host shells invoke the binary literally.

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
binary, and forwards `--agent` / `--event` as argv so `run.Serve` can force
dialect and event selection. When those flags are omitted, Serve detects the
dialect from the payload and peeks `hook_event_name`. Hint vs payload
disagreements warn on stderr and still use the hint. The same flags also
identify managed config entries for install and doctor. When the selected event
has no registered handlers, Serve warns on stderr and exits 0 without decoding.

On a cache miss, wat builds a bootstrap that imports `.wat/hooks.go` and calls
`run.Serve(Hooks...)`. The cache key includes all files under `.wat/` except
`.cache/` and `testdata/`, the wat version, Go version, target OS/architecture,
`GOFLAGS`, `CGO_ENABLED`, and bootstrap format. The binary is stored below:

```text
.wat/.cache/<content-hash>/hooks[.exe]
```

Build failures are fail-open by default (exit 1). `wat run --fail-closed`
returns exit 2 for a build failure so permission-gating hosts can deny the
operation.

## Test with fixtures

```bash
wat test --agent <dialect> --fixture <file|-> [--expect <file>] [--verbose]
```

`--agent` and `--fixture` are required. A fixture must be non-empty JSON. Event
identity comes from `--event` when set, otherwise from the fixture’s
`hook_event_name` (required when `--event` is omitted). Optional `--event` is
forwarded to the hooks binary as a dispatch hint (same as `wat run`). Use `-`
to read the fixture from stdin:

```powershell
Get-Content fixture.json | wat test --agent cursor --fixture -
```

```bash
cat fixture.json | wat test --agent cursor --fixture -
```

The report includes the fixture dialect/event, hook stdout, a recognized
decision field, and the hook exit code. `--verbose` also includes hook stderr.

### Optional expect documents

When `<fixture>.expect.json` exists, or when `--expect` points at a JSON file,
`wat test` asserts the hook outcome after printing the report:

| Field | Meaning |
|---|---|
| `exit` | Exact hook process exit code |
| `decision` | Recognized decision value from stdout (`deny`, `allow`, …) |
| `stdout_contains` | Substrings that must appear in hook stdout |
| `stdout_empty` | Whether trimmed stdout must be empty |

Unknown JSON fields are rejected. A matching expect run exits `0` even when the
hook itself returned a denial exit code. A mismatch exits `1` and lists failed
assertions on stderr. Without an expect document, the command still returns the
hook binary's exit code.

`wat init` seeds `.wat/testdata/<agent>/session_start.json` plus matching
`.expect.json` sidecars. E2E-only payloads live under
[`e2e/testdata/`](../e2e/testdata/). Use native event names and payload shapes
for the selected dialect.

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

Missing install wiring (absent agent configs or hook entries) and `wat` not
being on `PATH` are warnings: hooks simply will not be invoked until you run
`wat install`. Broken configs, duplicates, and toolchain/project failures
remain hard failures.

Status labels are colored on a TTY (`PASS` green, `WARN` yellow, `FAIL` red).
Set `NO_COLOR` to disable, or `FORCE_COLOR` to enable when stdout is not a
terminal.

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
| `0` | Successful CLI operation or hook execution; or a matching `wat test` expect run |
| `1` | Invalid CLI usage; hook build failure in default mode; hook runtime/protocol error; or `wat test` expect mismatch |
| `2` | `wat run --fail-closed` build failure, or a native host denial that uses exit 2 |
| `3` | CLI runtime failure such as missing project, unreadable fixture, failed install, or missing version build info |
| `4` | At least one `wat doctor` check failed |

For `wat test` without an expect document, a successfully executed fixture
returns the hook binary's exit code rather than translating it to a CLI-only
code. With an expect document, exit `0` means the assertions passed.

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

Source files under `.wat/` (except `.cache/` and `testdata/`) participate in the
cache key. Check that the command is resolving the intended project and that
`WAT_PROJECT_DIR` is not pointing elsewhere.
