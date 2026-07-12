# wat

**wat** is a Go CLI and library stack for authoring hook scripts that run under Claude Code, GitHub Copilot, and Cursor. The redesign provides a unified `agenthooks` library, per-agent SDKs, and a `wat` command that compiles and installs user hook scripts.

## Build

```bash
go build ./cmd/wat
```

On Windows the toolchain writes `wat.exe`; on Unix, `wat`. Avoid `go build -o wat` on Windows — use `./cmd/wat` or `-o wat.exe`.

## Test and lint

```bash
go vet ./...
go test ./...
golangci-lint run ./...
```

Install golangci-lint once:

```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.8
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

Browse package docs: `go doc github.com/sviatsviatsviat/wat/agenthooks`

Architecture and task breakdown live in a local-only `plan/` directory (gitignored).

## License

Apache License 2.0 — see [LICENSE](LICENSE).
