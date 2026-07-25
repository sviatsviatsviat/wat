# wat

**wat** is a Go CLI and library stack for authoring hook scripts that run under Claude Code, GitHub Copilot, and Cursor. The redesign provides a unified `sdk/agnostic` library, per-agent SDKs under `sdk/`, and a `wat` command that compiles and installs user hook scripts.

## Build

```bash
go build ./cmd/wat
```

On Windows the toolchain writes `wat.exe`; on Unix, `wat`. Avoid `go build -o wat` on Windows — use `./cmd/wat` or `-o wat.exe`.

## Commands

| Command | Description |
|---------|-------------|
| `wat init` | Scaffold a `.wat/` hook project |
| `wat install` | Install only the native events registered by `.wat/hooks.go` |
| `wat run` | Execute `.wat/hooks.go` on hook invocation |
| `wat test` | Run hook script against fixture payloads |
| `wat doctor` | Verify toolchain, script, cache, and install state |

Run `wat help` for the command list or `wat <command> -h` for per-command flags.

## Hook projects

`wat init` creates an importable `.wat/hooks.go` package. Export one `[]run.Hooks`
slice containing portable and/or agent-native fluent registrations:

```go
package hooks

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic"
	"github.com/sviatsviatsviat/wat/sdk/cursor"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// Hooks contains this project's hook registrations.
var Hooks = []run.Hooks{
	agnostic.UseHooks().
		OnPreTool(preTool).
		OnStop(stop),
	cursor.UseHooks().
		BeforeShellExecution(beforeShell),
}
```

`wat` generates the executable bootstrap and caches the resulting binary.
`wat install` inspects `Hooks`, expands portable registrations to their native
agent events, installs one entry per registered event, and removes stale
wat-managed entries. Multiple handlers for the same native event share one
installed entry and run in registration order.

Package initialization runs when `wat install`, `wat doctor`, `wat test`, or a
live hook loads the binary. Keep registration expressions and `init` functions
free of external side effects; perform runtime work inside handlers.

## Exit codes

| Code | Meaning |
|------|---------|
| `0` | Success |
| `1` | Usage error (unknown command, invalid flags) |
| `2` | `wat run --fail-closed` build failure (block/deny) |
| `3` | Runtime failure (missing project, I/O error) |
| `4` | `wat doctor` check failure |

## Test and lint

```bash
go vet ./...
go test ./...
golangci-lint run ./...
```

Install golangci-lint once:

```bash
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2
```

Verify: `go vet ./... && go test ./... && golangci-lint run ./...`

## Dev Container

Linux development environment with Go 1.26 and pinned lint tools. See [.devcontainer/README.md](.devcontainer/README.md) for setup (named volume at `/workspaces/wat`, manual `git clone`).

## Documentation

| Doc | Purpose |
|-----|---------|
| [docs/README.md](docs/README.md) | Documentation index |
| [docs/agent-formats.md](docs/agent-formats.md) | Tool and MCP naming across agents |
| [CHANGELOG.md](CHANGELOG.md) | Release notes |
| [AGENTS.md](AGENTS.md) | Contributor and agent workflow conventions |

Browse package docs: `go doc github.com/sviatsviatsviat/wat/sdk/agnostic`

Architecture and task breakdown live in a local-only `plan/` directory (gitignored).

## License

Apache License 2.0 — see [LICENSE](LICENSE).
