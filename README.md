# wat

`wat` lets you write one Go hook project and run it under Claude Code, GitHub
Copilot, and Cursor. It provides:

- a CLI that scaffolds, builds, installs, tests, and diagnoses hook projects;
- a portable SDK for behavior shared by all three agents;
- native SDKs for agent-specific events and response fields.

The project currently requires Go 1.26.

## Quick start

Build and install the CLI from this checkout:

```bash
go install ./cmd/wat
```

From the repository where you want hooks:

```bash
wat init
wat install
wat doctor
```

`wat init` creates an independent Go module in `.wat/`. Its `hooks.go` exports
the registrations that `wat install` inspects and installs:

```go
package hooks

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/agnostic"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// Hooks is the hook project's public contract with wat.
var Hooks = []run.Hooks{
	agnostic.UseHooks().OnSessionStart(func(ctx context.Context, hook agnostic.SessionStartEvent, r agnostic.SessionStartResults) (agnostic.SessionStartResult, error) {
		return r.Context("wat hooks are active"), nil
	}),
}
```

Returning `nil` means that a handler has no opinion. Response-producing
handlers receive a hook-specific result builder; observe-only handlers return
only an error.

Test the hook without starting an agent:

```bash
wat test --agent claude --fixture testdata/fixtures/claude/session_start.json
```

Run the command from a directory inside the target project, or set
`WAT_PROJECT_DIR` to its root.

## Commands

| Command | Purpose |
|---|---|
| `wat init [path]` | Create `.wat/go.mod`, `.wat/hooks.go`, and the build cache directory |
| `wat install` | Reconcile registered native events into the selected agent configs |
| `wat run` | Build on a cache miss, then execute the hook process against stdin |
| `wat test` | Execute a registered hook against a JSON fixture and print a report |
| `wat doctor` | Check Go, the hook project, cache, manifest, and installed configs |

Use `wat help` and `wat <command> -h` for the exact flags.

## Choose an SDK

Use [`sdk/agnostic`](sdk/agnostic/) when the behavior must work across all
agents. It normalizes lifecycle, prompt, tool, subagent, stop, and compaction
events and maps portable results back to each native protocol.

Use [`sdk/claude`](sdk/claude/), [`sdk/copilot`](sdk/copilot/), or
[`sdk/cursor`](sdk/cursor/) when you need a native event or response field that
cannot be represented portably. Portable and native registrations can coexist
in the same `Hooks` slice.

Use [`sdk/run`](sdk/run/) only when building a standalone hook executable or
inspecting registrations programmatically. The `wat` CLI generates this
bootstrap for `.wat/` projects.

## Documentation

| Document | Read it for |
|---|---|
| [Using wat](docs/usage.md) | Installation, project discovery, commands, cache behavior, exit codes, and troubleshooting |
| [SDK API](docs/sdk.md) | Portable and native registration APIs, event/result patterns, and examples |
| [Architecture](docs/architecture.md) | Package boundaries, request lifecycle, and implementation patterns |
| [Agent protocols](docs/agent-formats.md) | Event coverage, tool normalization, and native protocol differences |
| [Contributing](CONTRIBUTING.md) | Development workflow, tests, docs, commits, and pull requests |
| [Changelog](CHANGELOG.md) | User-visible capabilities planned for the next release |

Browse package documentation locally:

```bash
go doc github.com/sviatsviatsviat/wat/sdk/agnostic
go doc github.com/sviatsviatsviat/wat/sdk/claude
go doc github.com/sviatsviatsviat/wat/sdk/run
```

## Development

```bash
go vet ./...
go test ./...
golangci-lint run ./...
go build ./cmd/wat
```

The repository also provides a pinned Linux development environment in
[`.devcontainer/`](.devcontainer/README.md).

## License

Apache License 2.0. See [LICENSE](LICENSE).
